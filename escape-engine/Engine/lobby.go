package Engine

import (
	"escape-api/LogUtil"
	"escape-engine/Models"
	"fmt"
	"log"
	"slices"
	"strings"
)

// Given a GameDefinition's ID and a player name, creates and saves a new lobby for that player's game, returning the Lobby's room code.
func CreateRoom(mapId string) (string, error) {
	funcLogPrefix := "==CreateRoom=="
	defer LogUtil.EnsureLogPrefixIsReset()
	LogUtil.SetLogPrefix(ModuleLogPrefix, PackageLogPrefix)

	log.Printf("%s Getting Map for Id == {%s}", funcLogPrefix, mapId)
	requestedMap, err := GetMapFromDB(mapId)
	if err != nil {
		LogError(funcLogPrefix, err)
		return "", err
	}

	log.Printf("%s Creating lobby object", funcLogPrefix)
	lobby := Models.Lobby{
		MapId:      requestedMap.Id,
		Status:     Models.LobbyStatus_AwaitingStart,
		MaxPlayers: 12, //Setting just on my own for now
		NumPlayers: 0,
		Players:    []Models.Player{},
		Host:       Models.Player{},
	}

	log.Printf("%s Generating Room Code", funcLogPrefix)
	roomCode := generateRoomCode()

	log.Printf("%s Room Code successfully generated. Assigning RoomCode {%s} to Lobby", funcLogPrefix, roomCode)
	lobby.RoomCode = roomCode

	log.Printf("%s Saving Lobby to Redis", funcLogPrefix)
	lobby, err = SaveLobbyToFs(lobby)
	if err != nil {
		LogError(funcLogPrefix, err)
		return "", err
	}

	log.Printf("%s Lobby Created & Saved. Returning RoomCode", funcLogPrefix)
	return roomCode, nil
}

// Creates a Player object for the given PlayerName and attempts to add them to the lobby with the given RoomCode. On Success, returns the new state
// of the lobby, the Player's assigned Id, and any error that occurred
func JoinRoom(roomCode string, playerName string) (Models.Lobby, string, error) {
	funcLogPrefix := "==JoinRoom=="
	defer LogUtil.EnsureLogPrefixIsReset()
	LogUtil.SetLogPrefix(ModuleLogPrefix, PackageLogPrefix)

	log.Printf("%s Recieved request from Player {%s} to join lobby with RoomCode == {%s}", funcLogPrefix, playerName, roomCode)

	lobby, err := GetLobbyFromFs(strings.ToUpper(roomCode))
	if err != nil {
		LogError(funcLogPrefix, err)
		return Models.Lobby{}, "", err
	}

	//Only allow player to join if there's room & the game hasn't started yet (i.e. Status == LobbyStatus_AwaitingStart)
	if lobby.NumPlayers >= lobby.MaxPlayers {
		log.Printf("%s ERROR: Lobby's max player count {%d} already reached. Player cannot join!", funcLogPrefix, lobby.MaxPlayers)
		return Models.Lobby{}, "", fmt.Errorf("Lobby's max player count {%d} already reached", lobby.MaxPlayers)
	}
	if lobby.Status != Models.LobbyStatus_AwaitingStart {
		log.Printf("%s Error: Game has already started. Player cannot join!", funcLogPrefix)
		return Models.Lobby{}, "", fmt.Errorf("Game has already started!")
	}
	if slices.ContainsFunc(lobby.Players, func(p Models.Player) bool { return p.Name == playerName }) {
		log.Printf("%s Error: Player name {%s} already taken. Player cannot join!", funcLogPrefix, playerName)
		return Models.Lobby{}, "", fmt.Errorf("Name already taken!")
	}

	thisPlayer := createPlayerObject(playerName)

	//Create a copy, in case anything goes wrong
	updatedLobby := lobby
	updatedLobby.Players = slices.Clone(lobby.Players)

	log.Printf("%s Adding player {%s} to lobby's Player List", funcLogPrefix, playerName)

	//If this player is the first to join, set them as the host
	if updatedLobby.NumPlayers == 0 {
		updatedLobby.Host = thisPlayer
	}

	updatedLobby.Players = append(lobby.Players, thisPlayer)
	updatedLobby.NumPlayers = len(updatedLobby.Players)

	log.Printf("%s Player added. Caching new Lobby", funcLogPrefix)
	saved, err := SaveLobbyToFs(updatedLobby)
	if err != nil { //If something goes wrong, re-save and return the version without any changes
		SaveLobbyToFs(lobby)
		return Models.Lobby{}, "", err
	}

	log.Printf("%s Lobby joined and saved. Returning Lobby", funcLogPrefix)
	return saved, thisPlayer.Id, nil
}

func LeaveRoom(roomCode string, playerId string) (Models.Lobby, error) {
	funcLogPrefix := "==LeaveRoom=="
	defer LogUtil.EnsureLogPrefixIsReset()
	LogUtil.SetLogPrefix(ModuleLogPrefix, PackageLogPrefix)

	log.Printf("%s Recieved request to remove Player {%s} from lobby with RoomCode == {%s}", funcLogPrefix, playerId, roomCode)

	lobby, err := GetLobbyFromFs(roomCode)
	if err != nil {
		LogError(funcLogPrefix, err)
		return Models.Lobby{}, err
	}

	//Create a copy, in case anything goes wrong
	updatedLobby := lobby
	updatedLobby.Players = slices.Clone(lobby.Players)

	log.Printf("%s Removing player {%s} from lobby's Player List", funcLogPrefix, playerId)

	newPlayers := []Models.Player{}

	for _, player := range updatedLobby.Players {
		if player.Id != playerId {
			newPlayers = append(newPlayers, player)
		}
	}

	updatedLobby.Players = newPlayers
	updatedLobby.NumPlayers = len(newPlayers)

	log.Printf("%s Player Removed. Caching new Lobby", funcLogPrefix)
	saved, err := SaveLobbyToFs(updatedLobby)
	if err != nil { //If something goes wrong, re-save and return the version without any changes
		SaveLobbyToFs(lobby)
		return Models.Lobby{}, err
	}

	//If the game has started, we need to remove them from the GameState too
	if saved.Status == Models.LobbyStatus_InProgress {
		log.Println("Player is being removed from an in-progress game. Removing player from GameState...")
		gameState, err := GetGameStateFromFs(saved.GameStateId)
		if err != nil {
			LogError(funcLogPrefix, err)
		}

		//If it's this player's turn, end their turn before removing them
		if gameState.CurrentPlayer == playerId {
			log.Println("GameState is listing Player as CurrentPlayer. Ending their turn before removal...")
			//Models.EndTurn{}.Execute(&gameState, playerId) TODO: Reimplement
		}

		//Remove from Player list
		currentPlayers := slices.Clone(gameState.Players)
		newPlayers = []Models.Player{}

		for _, player := range currentPlayers {
			if player.Id != playerId {
				newPlayers = append(newPlayers, player)
			}
		}

		gameState.Players = newPlayers

		log.Println("Player has been removed from GameState. Caching new GameState now...")
		_, err = SaveGameStateToFs(gameState)
		if err != nil {
			LogError(funcLogPrefix, err)
		}

	}

	log.Printf("%s Left Lobby. Returning Lobby", funcLogPrefix)
	return saved, nil
}

// Creates a player object for a given name. Does NOT assign team, role, or starting position
func createPlayerObject(name string) Models.Player {
	funcLogPrefix := "==CreatePlayerObject=="
	LogUtil.SetLogPrefix(ModuleLogPrefix, PackageLogPrefix)
	defer LogUtil.EnsureLogPrefixIsReset()

	log.Printf("%s Creating Player object for Player name {%s}", funcLogPrefix, name)

	return Models.Player{
		Id:   GenerateId(),
		Name: name,
	}
}
