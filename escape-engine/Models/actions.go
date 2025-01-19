package Models

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"math/rand/v2"
	"slices"
)

const (
	Action_Movement = "Movement"
	Action_Attack   = "Attack"
)

// This is the way the frontend will send data to the backend during gameplay. They will
// send one of these objects, then the Rule Engine will take it, perform any updates to the
// internal model of the Game, then respond with a Changelog
type SubmittedAction struct {
	//The type of Turn you want to take. Should match exactly with the name of one of the below structs (i.e. "Movement", "Insertion", etc)
	Type string `json:"type"`
	//The actual turn object. Should have all the fields within the struct that you're wanting
	Turn json.RawMessage `json:"turn"`
	//ID of the player who is trying to submit this action. This is now supplied by the backend
	PlayerId string `json:"playerId"`
}

type Movement struct {
	ToRow int `json:"toRow"`
	ToCol int `json:"toCol"`
}

type Attack struct {
	Row int `json:"row"`
	Col int `json:"col"`
}

func (move Movement) Execute(gameState *GameState, playerId string) (*GameEvent, error) {
	var gameEvent *GameEvent = nil
	// if gameState.CurrentPlayer != playerId {
	// 	return nil, fmt.Errorf("player trying to execute turn is not current player")
	// }

	actingPlayerIndex := slices.IndexFunc(gameState.Players, func(p Player) bool { return playerId == p.Id })
	if actingPlayerIndex == -1 {
		return nil, fmt.Errorf("could not find acting player with ID == {%s}", playerId)
	}

	actingPlayer := &(gameState.Players[actingPlayerIndex])

	//Bounds check
	spaceKey := fmt.Sprintf("%d,%d", move.ToRow, move.ToCol)

	if space, exists := gameState.GameMap.Spaces[spaceKey]; exists {
		//Make sure it's an open space
		cantMoveInto := []int{
			Space_AlienStart,
			Space_HumanStart,
			Space_Wall,
			Space_UsedPod,
		}
		if slices.ContainsFunc(cantMoveInto, func(t int) bool { return space.Type == t }) {
			return nil, fmt.Errorf("space [%d,%d] is not allowed to be moved into", move.ToRow, move.ToCol)
		}

		//Make sure it's close enough
		allowedSpaces := 1 //TODO: Figure out how to deal with aliens getting 3 spaces of movement after a kill
		if actingPlayer.Team == PlayerTeam_Alien {
			allowedSpaces = 2
		}
		if !checkMovement(move.ToRow, actingPlayer.Row, move.ToCol, actingPlayer.Col, allowedSpaces) {
			//return fmt.Errorf("movement not allowed") TODO: Turned off for now because hex grids make counting spaces HARD
		}

		if space.Type == Space_Pod {
			if actingPlayer.Team == PlayerTeam_Alien {
				return nil, fmt.Errorf("aliens are not allowed to enter pods")
			}

			totalPodCards := gameState.GameConfig.NumWorkingPods + gameState.GameConfig.NumBrokenPods
			if totalPodCards == 0 {
				return nil, fmt.Errorf("no pod cards left")
			}
			podCard := (rand.IntN(totalPodCards) + 1)
			podIsWorking := podCard > gameState.GameConfig.NumWorkingPods
			if gameState.GameConfig.NumBrokenPods == 0 { //0 broken pods is an edge case that will effectively make 1 "working" card act as a broken card
				podIsWorking = true
			}
			if gameState.GameConfig.NumWorkingPods == 0 { //0 working pods is an edge case in the opposite direction
				podIsWorking = false
			}
			if podIsWorking {
				gameEvent = &GameEvent{
					Row:         move.ToRow,
					Col:         move.ToCol,
					Description: fmt.Sprintf("Player %s escaped using the Pod at (%d,%d)!", actingPlayer.Name, move.ToRow, move.ToCol),
				}
				actingPlayer.Team = PlayerTeam_Spectator
				actingPlayer.Row, actingPlayer.Col = -99, -99

				gameState.GameConfig.NumWorkingPods -= 1
				gameState.GameMap.Spaces[space.GetMapKey()] = Space{
					Row:  space.Row,
					Col:  space.Col,
					Type: Space_UsedPod,
				}

				return gameEvent, nil
			} else {
				gameEvent = &GameEvent{
					Row:         move.ToRow,
					Col:         move.ToCol,
					Description: fmt.Sprintf("Player %s tried the Pod at (%d,%d) but it didn't work!", actingPlayer.Name, move.ToRow, move.ToCol),
				}

				gameState.GameConfig.NumBrokenPods -= 1
				gameState.GameMap.Spaces[space.GetMapKey()] = Space{
					Row:  space.Row,
					Col:  space.Col,
					Type: Space_UsedPod,
				}
			}
		}

		//At this point, player is allowed to execute the move
		actingPlayer.Row = move.ToRow
		actingPlayer.Col = move.ToCol
	} else {
		return nil, fmt.Errorf("requested space [%d,%d] not found in map", move.ToRow, move.ToCol)
	}

	return gameEvent, nil
}

func checkMovement(toRow int, fromRow int, toCol int, fromCol int, allowedMovement int) bool {
	log.Printf("Row (%d->%d) Col(%d->%d), %d", fromRow, toRow, fromCol, toCol, allowedMovement)

	if int(math.Abs(float64(toRow-fromRow))) > allowedMovement {
		return false
	}

	allowedMovementOffset := allowedMovement - 1
	if fromCol%2 == 0 { //Even column
		log.Println("Even column")
		if toRow-fromRow > 0 {
			log.Println("Moving down")
			return int(math.Abs(float64(toCol-fromCol))) <= allowedMovementOffset
		} else {
			log.Println("moving up")
			return int(math.Abs(float64(toCol-fromCol))) <= allowedMovement
		}
	} else { //Odd column
		log.Println("Odd column")
		if toRow-fromRow < 0 { //Moving up
			log.Println("Moving up")
			return int(math.Abs(float64(toCol-fromCol))) <= allowedMovementOffset
		} else {
			log.Println("moving down")
			return int(math.Abs(float64(toCol-fromCol))) <= allowedMovement
		}
	}
}

func (attack Attack) Execute(gameState *GameState, playerId string) (*GameEvent, error) {
	actingPlayerIndex := slices.IndexFunc(gameState.Players, func(p Player) bool { return playerId == p.Id })
	if actingPlayerIndex == -1 {
		return nil, fmt.Errorf("could not find acting player with ID == {%s}", playerId)
	}

	actingPlayer := &(gameState.Players[actingPlayerIndex])

	var gameEvent *GameEvent = &GameEvent{
		Row:         attack.Row,
		Col:         attack.Col,
		Description: fmt.Sprintf("%s attacked [%d,%d]!\n", actingPlayer.Name, attack.Row, attack.Col),
	}
	alienStarts := gameState.GameMap.GetSpacesOfType(Space_AlienStart)

	for index, player := range gameState.Players {
		if player.Id == actingPlayer.Id { //Don't kill the player doing the attacking
			continue
		}
		newSpaceForPlayer := alienStarts[rand.IntN(len(alienStarts))]

		gameState.Players[index].Team = PlayerTeam_Alien
		gameState.Players[index].Row, gameState.Players[index].Col = newSpaceForPlayer.Row, newSpaceForPlayer.Col
		gameEvent.Description = string(fmt.Appendf([]byte(gameEvent.Description), "%s died!\n", gameState.Players[index].Name))
	}

	return gameEvent, nil
}
