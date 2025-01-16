package Engine

import (
	"escape-api/LogUtil"
	"escape-engine/Models"
	"fmt"
	"log"
	"math/rand"
)

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

	for _, element := range lobby.Players { //TODO: Assign role
		gameState.Players = append(gameState.Players, Models.Player{
			Id:   element.Id,
			Name: element.Name,
			//Using -99 to avoid any weird cases where that player might be close enough to get onto the Map
			Row: -99,
			Col: -99,
		})
	}

	assignTeams(&gameState)
	assignStartingPositions(&gameState, &mapDef)

	gameState.CurrentPlayer = gameState.Players[rand.Intn(len(gameState.Players))].Id

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

// #region Helper Functions

// Assigns teams randomly to all players in the GameState. If a player cannot be assigned for any reason, they are assigned as a spectator
func assignTeams(gameState *Models.GameState) {
	humansToAssign, aliensToAssign := gameState.GameConfig.NumHumans, gameState.GameConfig.NumAliens
	for index := range gameState.Players {
		if humansToAssign == 0 && aliensToAssign != 0 { //No alien slots left, must be human
			gameState.Players[index].Team = Models.PlayerTeam_Human
		} else if aliensToAssign == 0 && humansToAssign != 0 { //No human slots left, must be alien
			gameState.Players[index].Team = Models.PlayerTeam_Alien
		} else {
			if humansToAssign == 0 && aliensToAssign == 0 {
				log.Printf("assignTeams ERROR! No team slot left to assign to player %s! Assigning player as spectator!", gameState.Players[index].Name)
				gameState.Players[index].Team = Models.PlayerTeam_Spectator
			}
			if rand.Intn(2) == 0 {
				gameState.Players[index].Team = Models.PlayerTeam_Human
				humansToAssign--
			} else {
				gameState.Players[index].Team = Models.PlayerTeam_Alien
				aliensToAssign--
			}
		}
	}
}

// Assigns a random valid starting position to each Player from the pool of start spaces assigned to each team
func assignStartingPositions(gameState *Models.GameState, gameMap *Models.GameMap) {
	humanStarts, alienStarts := []Models.Space{}, []Models.Space{}
	for _, space := range gameMap.Spaces {
		if space.Type == Models.Space_AlienStart {
			alienStarts = append(alienStarts, space)
		} else if space.Type == Models.Space_HumanStart {
			humanStarts = append(humanStarts, space)
		}
	}

	for index, player := range gameState.Players {
		if player.Team == Models.PlayerTeam_Human {
			startingSpace := humanStarts[rand.Intn(len(humanStarts))]

			gameState.Players[index].Row, gameState.Players[index].Col = startingSpace.Row, startingSpace.Col
		} else if player.Team == Models.PlayerTeam_Alien {
			startingSpace := alienStarts[rand.Intn(len(humanStarts))]

			gameState.Players[index].Row, gameState.Players[index].Col = startingSpace.Row, startingSpace.Col
		} else if player.Team == Models.PlayerTeam_Spectator {
			gameState.Players[index].Row, gameState.Players[index].Col = -99, -99
		}
	}
}

//#endregion
