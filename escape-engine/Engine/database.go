package Engine

import (
	"encoding/json"
	"escape-api/LogUtil"
	"escape-engine/Models"
	"fmt"
	"log"
	"os"
	"time"
)

func PrepareFilesystem() {
	os.Mkdir("./maps", 0666)
	os.Mkdir("./lobbies", 0666)
	os.Mkdir("./gameStates", 0666)
	os.Mkdir("./recaps", 0666)
}

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
		fmt.Println("Ran into unrecoverable error trying to save map")
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

func SaveLobbyToFs(lobby Models.Lobby) (Models.Lobby, error) {
	funcLogPrefix := "==SaveLobbyToFs=="
	defer LogUtil.EnsureLogPrefixIsReset()
	LogUtil.SetLogPrefix(ModuleLogPrefix, PackageLogPrefix)
	asJson, err := json.Marshal(lobby)
	if err != nil {
		return lobby, err
	}

	filename := "lobby_" + lobby.RoomCode + ".json"

	f, err := os.Create(fmt.Sprintf("./lobbies/%s", filename))
	if err != nil {
		log.Printf("%s Ran into unrecoverable error trying to save lobby to filesystem", funcLogPrefix)
		f.Close()
		return lobby, err
	}
	_, err = f.Write(asJson)
	f.Close()
	if err != nil {
		return lobby, err
	}

	//Kick off goroutine clearing out unused lobbies
	go clearOutOldFiles("./lobbies/")

	return lobby, nil
}

func GetLobbyFromFs(roomCode string) (Models.Lobby, error) {
	funcLogPrefix := "==GetLobbyFromFs=="
	defer LogUtil.EnsureLogPrefixIsReset()
	LogUtil.SetLogPrefix(ModuleLogPrefix, PackageLogPrefix)

	log.Printf("%s Getting lobby from FS with RoomCode == {%s}", funcLogPrefix, roomCode)
	data, err := os.ReadFile(fmt.Sprintf("./lobbies/lobby_%s.json", roomCode))
	if err != nil {
		LogError(funcLogPrefix, err)
		return Models.Lobby{}, err
	}

	parsed := Models.Lobby{}

	err = json.Unmarshal(data, &parsed)
	if err != nil {
		return parsed, err
	}

	return parsed, nil
}

func SaveGameStateToFs(gameState Models.GameState) (Models.GameState, error) {
	funcLogPrefix := "==SaveGameStateToFs==:"
	defer LogUtil.EnsureLogPrefixIsReset()
	LogUtil.SetLogPrefix(ModuleLogPrefix, PackageLogPrefix)

	log.Printf("%s Received GameState to save", funcLogPrefix)

	//If the gameState doesn't have an ID yet,
	//Generate one for it by simply using the Current UNIX time in milliseconds
	id := gameState.Id
	if id == "" {
		log.Printf("%s GameState does not yet have an ID. Generating new one.", funcLogPrefix)
		id = GenerateId()
		log.Printf("%s ID successfully generated. Assigning ID {%s} to GameState", funcLogPrefix, id)
		gameState.Id = id
	}

	asJson, err := json.Marshal(gameState)
	if err != nil {
		return gameState, err
	}

	filename := "gameState_" + gameState.Id + ".json"
	f, err := os.Create(fmt.Sprintf("./gameStates/%s", filename))
	if err != nil {
		log.Printf("%s Ran into unrecoverable error trying to save GameState to filesystem", funcLogPrefix)
		f.Close()
		return gameState, err
	}
	_, err = f.Write(asJson)
	f.Close()
	if err != nil {
		return gameState, err
	}

	//Kick off goroutine clearing out unused lobbies
	go clearOutOldFiles("./gameStates/")
	return gameState, nil
}

func GetGameStateFromFs(id string) (Models.GameState, error) {
	funcLogPrefix := "==GetGameStateFromFs=="
	defer LogUtil.EnsureLogPrefixIsReset()
	LogUtil.SetLogPrefix(ModuleLogPrefix, PackageLogPrefix)

	log.Printf("%s Getting GameState from FS with ID == {%s}", funcLogPrefix, id)
	data, err := os.ReadFile(fmt.Sprintf("./gameStates/gameState_%s.json", id))
	if err != nil {
		return Models.GameState{}, err
	}

	parsed := Models.GameState{}

	err = json.Unmarshal(data, &parsed)
	if err != nil {
		return parsed, err
	}

	return parsed, nil
}

func clearOutOldFiles(directory string) {
	files, _ := os.ReadDir(directory)
	for _, file := range files {
		fullFileName := directory + file.Name()
		stats, _ := os.Stat(fullFileName)
		expirationTime := stats.ModTime().AddDate(0, 0, 7)
		if time.Now().After(expirationTime) {
			os.Remove(fullFileName)
		}
	}
}
