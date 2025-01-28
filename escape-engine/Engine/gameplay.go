package Engine

import (
	"encoding/json"
	"escape-api/LogUtil"
	"escape-engine/Models"
	"fmt"
	"log"
	"math/rand"
	"slices"
)

// Given an id to a Game defition, constructs and returns an initial GameState for it. This is essentially
// how to start the game
func GetInitialGameState(roomCode string, gameConfig Models.GameConfig) (Models.GameState, error) {
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

	gameState.GameMap = mapDef
	gameState.GameConfig = gameConfig

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

func SubmitAction(gameId string, action Models.SubmittedAction) ([]Models.WebsocketMessageListItem, error) {
	funcLogPrefix := "==SubmitAction=="
	defer LogUtil.EnsureLogPrefixIsReset()
	LogUtil.SetLogPrefix(ModuleLogPrefix, PackageLogPrefix)

	gameState, err := GetCachedGameStateFromRedis(gameId)
	if err != nil {
		LogError(funcLogPrefix, err)
		return []Models.WebsocketMessageListItem{}, err
	}

	//TODO: Add turn enforcement

	messages := []Models.WebsocketMessageListItem{}

	switch action.Type {
	case Models.Action_Attack:
		turn := Models.Attack{}
		err := json.Unmarshal(action.Turn, &turn)
		if err != nil {
			LogError(funcLogPrefix, err)
			return messages, err
		}

		data := Models.WebsocketMessage{}
		shouldBroadcast := false

		var result any = nil
		if turn.IsAttacking() {
			result, err = turn.Execute(&gameState, action.PlayerId)
			data.Type = Models.WebsocketMessage_GameEvent
			shouldBroadcast = true
		} else {
			// Cards are weird to deal with. You should only draw a card if you're in a dangerous sector,
			// so DrawCard will check where you are and set the type to Card_NoCard if so. In that case,
			// everyone should be told the player has moved into a safe sector, effectively skipping over
			// the Noise phase of their turn
			cardEvent, er := Models.DrawCard(&gameState, action.PlayerId)
			err = er
			if cardEvent.Type == Models.Card_NoCard {
				actingPlayerIndex := slices.IndexFunc(gameState.Players, func(p Models.Player) bool { return action.PlayerId == p.Id })
				if actingPlayerIndex == -1 {
					err = fmt.Errorf("could not find acting player with ID == {%s}", action.PlayerId)
				}

				actingPlayer := &(gameState.Players[actingPlayerIndex])

				data.Type = Models.WebsocketMessage_GameEvent
				shouldBroadcast = true
				result = Models.GameEvent{
					Row:         -99,
					Col:         -99,
					Description: fmt.Sprintf("Player '%s' is in a safe sector", actingPlayer.Name),
				}
			} else {
				data.Type = Models.WebsocketMessage_Card
				shouldBroadcast = false
				result = cardEvent
			}
		}

		if err != nil {
			LogError(funcLogPrefix, err)
			return messages, err
		}

		data.Data = result
		messages = append(messages, Models.WebsocketMessageListItem{
			Message:         data,
			ShouldBroadcast: shouldBroadcast,
		})
	case Models.Action_Movement:
		turn := Models.Movement{}
		err := json.Unmarshal(action.Turn, &turn)
		if err != nil {
			LogError(funcLogPrefix, err)
			return messages, err
		}

		result, err := turn.Execute(&gameState, action.PlayerId)
		if err != nil {
			LogError(funcLogPrefix, err)
			return messages, err
		}

		messages = append(messages, Models.WebsocketMessageListItem{
			Message: Models.WebsocketMessage{
				Type: Models.WebsocketMessage_MovementResponse,
				Data: result,
			},
			ShouldBroadcast: false,
		})
	case Models.Action_Noise:
		turn := Models.Noise{}
		err := json.Unmarshal(action.Turn, &turn)
		if err != nil {
			LogError(funcLogPrefix, err)
			return messages, err
		}

		result, err := turn.Execute(&gameState, action.PlayerId)
		if err != nil {
			LogError(funcLogPrefix, err)
			return messages, err
		}

		messages = append(messages, Models.WebsocketMessageListItem{
			ShouldBroadcast: true,
			Message: Models.WebsocketMessage{
				Type: Models.WebsocketMessage_GameEvent,
				Data: result,
			},
		})

	case Models.Action_EndTurn:
		turn := Models.EndTurn{}
		err := json.Unmarshal(action.Turn, &turn)
		if err != nil {
			LogError(funcLogPrefix, err)
			return messages, err
		}

		result, event, err := turn.Execute(&gameState, action.PlayerId) //TODO: Figure out how to tell clients someone used a pod
		if err != nil {
			LogError(funcLogPrefix, err)
			return messages, err
		}

		if event != nil {
			messages = append(messages, Models.WebsocketMessageListItem{
				Message: Models.WebsocketMessage{
					Type: Models.WebsocketMessage_GameEvent,
					Data: event,
				},
				ShouldBroadcast: true,
			})
		}

		if result != nil {
			messages = append(messages, Models.WebsocketMessageListItem{
				Message: Models.WebsocketMessage{
					Type: Models.WebsocketMessage_GameState,
					Data: result,
				},
				ShouldBroadcast: true,
			})
		}
	}

	//Re-save gamestate
	_, err = CacheGameStateInRedis(gameState)
	if err != nil {
		LogError(funcLogPrefix, err)
		return messages, err
	}

	return messages, nil
}

// #region Helper Functions

// Assigns teams randomly to all players in the GameState. If a player cannot be assigned for any reason, they are assigned as a spectator
func assignTeams(gameState *Models.GameState) {
	log.Println("Assigning teams")
	humansToAssign, aliensToAssign := gameState.GameConfig.NumHumans, gameState.GameConfig.NumAliens
	for index := range gameState.Players {
		if humansToAssign == 0 && aliensToAssign != 0 { //No human slots left, must be human
			gameState.Players[index].Team = Models.PlayerTeam_Alien
		} else if aliensToAssign == 0 && humansToAssign != 0 { //No alien slots left, must be alien
			gameState.Players[index].Team = Models.PlayerTeam_Human
		} else {
			if humansToAssign == 0 && aliensToAssign == 0 {
				gameState.Players[index].Team = Models.PlayerTeam_Spectator
			}
			if rand.Intn(2) == 0 {
				gameState.Players[index].Team = Models.PlayerTeam_Human
				humansToAssign = humansToAssign - 1
			} else {
				gameState.Players[index].Team = Models.PlayerTeam_Alien
				aliensToAssign = aliensToAssign - 1
			}
		}
	}
}

// Assigns a random valid starting position to each Player from the pool of start spaces assigned to each team
func assignStartingPositions(gameState *Models.GameState, gameMap *Models.GameMap) {
	log.Println("Assigning starting postitions")
	humanStarts, alienStarts := gameMap.GetSpacesOfType(Models.Space_HumanStart), gameMap.GetSpacesOfType(Models.Space_AlienStart)
	log.Println(humanStarts)
	for index, player := range gameState.Players {
		if player.Team == Models.PlayerTeam_Human {
			startingSpace := humanStarts[rand.Intn(len(humanStarts))]

			gameState.Players[index].Row, gameState.Players[index].Col = startingSpace.Row, startingSpace.Col
		} else if player.Team == Models.PlayerTeam_Alien {
			startingSpace := alienStarts[rand.Intn(len(alienStarts))]

			gameState.Players[index].Row, gameState.Players[index].Col = startingSpace.Row, startingSpace.Col
		} else if player.Team == Models.PlayerTeam_Spectator {
			gameState.Players[index].Row, gameState.Players[index].Col = -99, -99
		}
	}
}

//#endregion
