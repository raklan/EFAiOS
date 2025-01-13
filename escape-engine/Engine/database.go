package Engine

import (
	"encoding/json"
	"escape-api/LogUtil"
	"escape-engine/Models"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
)

// Yes, I know. I just REALLY didn't want to bring in an entire database JUST for this and Redis shouldn't be used for it
func SaveMapToDB(m Models.GameMap) (Models.GameMap, error) {
	if m.Id == "" {
		m.Id = GenerateId()
	}

	asJson, err := json.Marshal(m)
	if err != nil {
		return m, err
	}

	filename := "map_" + m.Id + ".json"
	f, err := os.Create(fmt.Sprintf("./maps/%s", filename))
	if err != nil {
		fmt.Println("FML")
		f.Close()
		return m, err
	}
	_, err = f.Write(asJson)
	f.Close()
	if err != nil {
		return m, err
	}

	return m, nil
}

func GetMapFromDB(mapId string) (Models.GameMap, error) {
	funcLogPrefix := "==GetMapFromDB=="
	defer LogUtil.EnsureLogPrefixIsReset()
	LogUtil.SetLogPrefix(ModuleLogPrefix, PackageLogPrefix)

	log.Printf("%s Getting map from DB with id == {%s}", funcLogPrefix, mapId)
	data, err := os.ReadFile(fmt.Sprintf("./maps/map_%s.json", mapId))
	if err != nil {
		return Models.GameMap{}, err
	}

	parsed := Models.GameMap{}

	err = json.Unmarshal(data, &parsed)
	if err != nil {
		return parsed, err
	}

	return parsed, nil
}

func GetAllMaps() ([]string, error) {
	funcLogPrefix := "==GetAllMaps=="
	defer LogUtil.EnsureLogPrefixIsReset()
	LogUtil.SetLogPrefix(ModuleLogPrefix, PackageLogPrefix)

	log.Printf("%s Getting all maps from DB", funcLogPrefix)

	files, err := os.ReadDir("./maps/")
	toReturn := []string{}
	if err != nil {
		return []string{}, err
	}

	for _, file := range files {
		if !file.IsDir() {
			toReturn = append(toReturn, file.Name())
		}
	}

	log.Printf("%s Found %d maps, returning list...", funcLogPrefix, len(toReturn))

	return toReturn, nil
}

func SaveLobbyInRedis(lobby Models.Lobby) (Models.Lobby, error) {
	funcLogPrefix := "==SaveLobbyInRedis=="
	defer LogUtil.EnsureLogPrefixIsReset()
	LogUtil.SetLogPrefix(ModuleLogPrefix, PackageLogPrefix)

	log.Printf("%s Recieved request to save lobby in Redis", funcLogPrefix)

	asJson, err := json.Marshal(lobby)
	if err != nil {
		LogError(funcLogPrefix, err)
		return Models.Lobby{}, err
	}

	key := "lobby:" + lobby.RoomCode
	expiry, _ := time.ParseDuration("168h")
	err = RDB.Set(ctx, key, asJson, expiry).Err()
	if err != nil {
		LogError(funcLogPrefix, err)
		return Models.Lobby{}, err
	}

	log.Printf("%s Lobby saved in Redis with key == {%s}", funcLogPrefix, key)
	return lobby, nil
}

func LoadLobbyFromRedis(roomCode string) (Models.Lobby, error) {
	defer LogUtil.EnsureLogPrefixIsReset()
	LogUtil.SetLogPrefix(ModuleLogPrefix, PackageLogPrefix)
	funcLogPrefix := "==LoadLobbyFromRedis==:"

	log.Printf("%s Retrieving Lobby with RoomCode=={%s} from Redis", funcLogPrefix, roomCode)

	lobby := Models.Lobby{}

	//Catch empty ID
	if roomCode == "" {
		log.Printf("%s ERROR! RoomCode cannot be empty. Returning empty Lobby", funcLogPrefix)
		return lobby, fmt.Errorf("%s Id cannot be empty", funcLogPrefix)
	}

	//Try to get the Game from Redis. If it doesn't exist, give a specific error for that
	def, err := RDB.Get(ctx, "lobby:"+roomCode).Result()
	if err == redis.Nil {
		log.Printf("%s Could not find cached lobby for roomCode \"%s\"...Returning Empty Lobby", funcLogPrefix, roomCode)
		return lobby, fmt.Errorf("%s No game for Id=={%s} found", funcLogPrefix, roomCode)
	} else if err != nil {
		LogError(funcLogPrefix, err)
		return lobby, err
	}

	//Result is just a JSON string, so we still need to deserialize/unmarshal it
	err = json.Unmarshal([]byte(def), &lobby)
	if err != nil {
		LogError(funcLogPrefix, err)
		return lobby, err
	}

	log.Printf("%s Found a lobby, returning result", funcLogPrefix)
	return lobby, nil
}

// Caches the given [gameState] in redis. Returns nil for [error] if everything goes well
func CacheGameStateInRedis(gameState Models.GameState) (Models.GameState, error) {
	funcLogPrefix := "==CacheGameStateInRedis==:"
	defer LogUtil.EnsureLogPrefixIsReset()
	LogUtil.SetLogPrefix(ModuleLogPrefix, PackageLogPrefix)

	log.Printf("%s Received GameState to cache", funcLogPrefix)

	//If the gameState doesn't have an ID yet,
	//Generate one for it by simply using the Current UNIX time in milliseconds
	id := gameState.Id
	if id == "" {
		log.Printf("%s GameState does not yet have an ID. Generating new one.", funcLogPrefix)
		id = GenerateId()
		log.Printf("%s ID successfully generated. Assigning ID {%s} to GameState", funcLogPrefix, id)
		gameState.Id = id
	}

	//Convert to string and save to Redis
	asJson, err := json.Marshal(gameState)
	if err != nil {
		LogError(funcLogPrefix, err)
		return gameState, err
	}

	key := "gameState:" + id
	expiry, _ := time.ParseDuration("168h")
	err = RDB.Set(ctx, key, asJson, expiry).Err()
	if err != nil {
		LogError(funcLogPrefix, err)
		return gameState, err
	}

	log.Printf("%s GameState cached with key=={%s}", funcLogPrefix, key)
	return gameState, nil
}

// Retrieves a gameState with an id == [id] from Redis. If everything goes well, then [error] is nil
func GetCachedGameStateFromRedis(id string) (Models.GameState, error) {
	funcLogPrefix := "==GetCachedGameStateFromRedis==:"
	defer LogUtil.EnsureLogPrefixIsReset()
	LogUtil.SetLogPrefix(ModuleLogPrefix, PackageLogPrefix)

	log.Printf("%s Received request to get cached GameState from Redis", funcLogPrefix)

	gameState := Models.GameState{}

	//Catch empty id string early
	if id == "" {
		log.Printf("%s ERROR! Id cannot be empty. Returning empty GameState", funcLogPrefix)
		return gameState, fmt.Errorf("%s Id cannot be empty", funcLogPrefix)
	}

	//Try to get the game from Redis. If it doesn't exist, fail gracefully
	game, err := RDB.Get(ctx, "gameState:"+id).Result()
	if err == redis.Nil {
		log.Printf("%s Could not find cached GameState for key \"%s\"...Returning Empty GameState", funcLogPrefix, id)
		return gameState, fmt.Errorf("%s No game for Id=={%s} found", funcLogPrefix, id)
	} else if err != nil {
		LogError(funcLogPrefix, err)
		return gameState, err
	}

	//game is a JSON string of a GameState, so unmarshal it
	err = json.Unmarshal([]byte(game), &gameState)
	if err != nil {
		LogError(funcLogPrefix, err)
		return gameState, err
	}

	log.Printf("%s Found a GameState, returning result", funcLogPrefix)
	return gameState, nil
}

// Given an id to a Game defition, constructs and returns an initial GameState for it. This is essentially
// how to start the game
func GetInitialGameState(roomCode string) (Models.GameState, error) {
	funcLogPrefix := "==GetInitialGameState=="
	defer LogUtil.EnsureLogPrefixIsReset()
	LogUtil.SetLogPrefix(ModuleLogPrefix, PackageLogPrefix)

	gameState := Models.GameState{}

	lobby, err := LoadLobbyFromRedis(roomCode)
	if err != nil {
		LogError(funcLogPrefix, err)
		return gameState, err
	}

	//Check if the lobby is already started.
	if lobby.Status == Models.LobbyStatus_InProgress {
		err := fmt.Errorf("tried to start game, but Lobby {%s} has been marked as In Progress and has a GameStateId == {%s}", lobby.RoomCode, lobby.GameStateId)
		LogError(funcLogPrefix, err)
		return gameState, err
	}

	mapDef, err := GetMapFromDB(lobby.MapId)
	if err != nil {
		LogError(funcLogPrefix, err)
		return gameState, err
	}

	gameState.MapId = mapDef.Id

	gameState.Players = []Models.Player{}

	for _, element := range lobby.Players { //TODO: Assign team, role, and starting position
		gameState.Players = append(gameState.Players, Models.Player{
			Id:   element.Id,
			Name: element.Name,
		})
	}

	gameState.CurrentPlayer = gameState.Players[0].Id //TODO: Make a better way to determine a starting player maybe?

	gameState, err = CacheGameStateInRedis(gameState)
	if err != nil {
		LogError(funcLogPrefix, err)
		return gameState, err
	}

	//Mark the lobby as started and fill in GameStateId
	lobby.GameStateId = gameState.Id
	lobby.Status = Models.LobbyStatus_InProgress
	_, err = SaveLobbyInRedis(lobby)
	if err != nil {
		LogError(funcLogPrefix, err)
		return gameState, err
	}

	return gameState, nil
}

func EndGame(roomCode string, playerId string) error {
	funcLogPrefix := "==EndGame=="
	defer LogUtil.EnsureLogPrefixIsReset()
	LogUtil.SetLogPrefix(ModuleLogPrefix, PackageLogPrefix)

	lobby, err := LoadLobbyFromRedis(roomCode)
	if err != nil {
		LogError(funcLogPrefix, err)
		return err
	}

	//Make sure that A) this player is the host and therefore allowed to end the game, and B) this game isn't already ended

	if lobby.Host.Id != playerId {
		return fmt.Errorf("player trying to end game is not host of lobby")
	}

	if lobby.Status == Models.LobbyStatus_Ended {
		return fmt.Errorf("game has already been marked as ended")
	}

	//Mark Game as ended and resave
	lobby.Status = Models.LobbyStatus_Ended

	_, err = SaveLobbyInRedis(lobby)

	//Return any error that occurred during saving, if any
	return err
}
