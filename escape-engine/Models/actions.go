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
	Action_Attack   = "Attack"
	Action_Movement = "Movement"
	Action_Noise    = "Noise"
	Action_EndTurn  = "EndTurn"
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

// Set Row & Col to -99 to indicate no attack
type Attack struct {
	Row int `json:"row"`
	Col int `json:"col"`
}

func (attack Attack) IsAttacking() bool {
	return attack.Row != -99 && attack.Col != -99
}

type Noise struct {
	Row int `json:"row"`
	Col int `json:"col"`
}

func (noise Noise) IsNoisy() bool {
	return noise.Row != -99 && noise.Col != -99
}

type EndTurn struct {
}

func (move Movement) Execute(gameState *GameState, playerId string) (MovementEvent, error) {
	movement := MovementEvent{}
	// if gameState.CurrentPlayer != playerId {
	// 	return nil, fmt.Errorf("player trying to execute turn is not current player")
	// }

	actingPlayerIndex := slices.IndexFunc(gameState.Players, func(p Player) bool { return playerId == p.Id })
	if actingPlayerIndex == -1 {
		return movement, fmt.Errorf("could not find acting player with ID == {%s}", playerId)
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
			return movement, fmt.Errorf("space [%d,%d] is not allowed to be moved into", move.ToRow, move.ToCol)
		}

		//Make sure it's close enough
		allowedSpaces := 1 //TODO: Figure out how to deal with aliens getting 3 spaces of movement after a kill
		if actingPlayer.Team == PlayerTeam_Alien {
			allowedSpaces = 2
		}
		if !checkMovement(move.ToRow, actingPlayer.Row, move.ToCol, actingPlayer.Col, allowedSpaces) {
			//return fmt.Errorf("movement not allowed") TODO: Turned off for now because hex grids make counting spaces HARD
		}

		//At this point, player is allowed to execute the move
		actingPlayer.Row, actingPlayer.Col = move.ToRow, move.ToCol
		movement.NewRow, movement.NewCol = actingPlayer.Row, actingPlayer.Col
	} else {
		return movement, fmt.Errorf("requested space [%d,%d] not found in map", move.ToRow, move.ToCol)
	}

	return movement, nil
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
		if player.Row != attack.Row || player.Col != attack.Col {
			continue
		}
		newSpaceForPlayer := alienStarts[rand.IntN(len(alienStarts))]

		gameState.Players[index].Team = PlayerTeam_Alien
		gameState.Players[index].Row, gameState.Players[index].Col = newSpaceForPlayer.Row, newSpaceForPlayer.Col
		gameEvent.Description = string(fmt.Appendf([]byte(gameEvent.Description), "%s died!\n", gameState.Players[index].Name))
	}

	return gameEvent, nil
}

func DrawCard(gameState *GameState, playerId string) (CardEvent, error) { //TODO: Full implementation
	event := CardEvent{}
	actingPlayerIndex := slices.IndexFunc(gameState.Players, func(p Player) bool { return playerId == p.Id })
	if actingPlayerIndex == -1 {
		return event, fmt.Errorf("could not find acting player with ID == {%s}", playerId)
	}

	actingPlayer := &(gameState.Players[actingPlayerIndex])

	currentSpace := Space{
		Row: actingPlayer.Row,
		Col: actingPlayer.Col,
	}

	log.Println("Player's Space:", currentSpace)
	log.Println("Space in Map:", gameState.GameMap.Spaces[currentSpace.GetMapKey()])

	if gameState.GameMap.Spaces[currentSpace.GetMapKey()].Type == Space_Safe {
		event.Type = Card_NoCard
	} else {
		switch rand.IntN(3) {
		case 0:
			event.Type = Card_Green
		case 1:
			event.Type = Card_Red
		case 2:
			event.Type = Card_White
		}
	}

	return event, nil
}

func (noise Noise) Execute(gameState *GameState, playerId string) (*GameEvent, error) {
	actingPlayerIndex := slices.IndexFunc(gameState.Players, func(p Player) bool { return playerId == p.Id })
	if actingPlayerIndex == -1 {
		return nil, fmt.Errorf("could not find acting player with ID == {%s}", playerId)
	}

	actingPlayer := &(gameState.Players[actingPlayerIndex])

	event := GameEvent{
		Row: noise.Row,
		Col: noise.Col,
	}
	if noise.IsNoisy() {
		event.Description = fmt.Sprintf("Player '%s' made noise at [%d,%d]!", actingPlayer.Name, noise.Row, noise.Col)
	} else {
		event.Description = fmt.Sprintf("Player '%s' avoided making noise", actingPlayer.Name)
	}
	return &event, nil
}

func (endTurn EndTurn) Execute(gameState *GameState, playerId string) (*GameState, *GameEvent, error) {
	var gameEvent *GameEvent = nil

	actingPlayerIndex := slices.IndexFunc(gameState.Players, func(p Player) bool { return playerId == p.Id })
	if actingPlayerIndex == -1 {
		return gameState, nil, fmt.Errorf("could not find acting player with ID == {%s}", playerId)
	}

	actingPlayer := &(gameState.Players[actingPlayerIndex])

	space := gameState.GameMap.Spaces[Space{Row: actingPlayer.Row, Col: actingPlayer.Col}.GetMapKey()]
	if space.Type == Space_Pod {
		if actingPlayer.Team == PlayerTeam_Alien {
			return gameState, nil, fmt.Errorf("aliens are not allowed to enter pods")
		}

		totalPodCards := gameState.GameConfig.NumWorkingPods + gameState.GameConfig.NumBrokenPods
		if totalPodCards == 0 {
			return gameState, nil, fmt.Errorf("no pod cards left")
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
				Row:         actingPlayer.Row,
				Col:         actingPlayer.Col,
				Description: fmt.Sprintf("Player %s escaped using the Pod at (%d,%d)!", actingPlayer.Name, actingPlayer.Row, actingPlayer.Col),
			}
			actingPlayer.Team = PlayerTeam_Spectator
			actingPlayer.Row, actingPlayer.Col = -99, -99

			gameState.GameConfig.NumWorkingPods -= 1
			gameState.GameMap.Spaces[space.GetMapKey()] = Space{
				Row:  space.Row,
				Col:  space.Col,
				Type: Space_UsedPod,
			}

			return gameState, gameEvent, nil
		} else {
			gameEvent = &GameEvent{
				Row:         actingPlayer.Row,
				Col:         actingPlayer.Col,
				Description: fmt.Sprintf("Player %s tried the Pod at (%d,%d), but it didn't work!", actingPlayer.Name, actingPlayer.Row, actingPlayer.Col),
			}

			gameState.GameConfig.NumBrokenPods -= 1
			gameState.GameMap.Spaces[space.GetMapKey()] = Space{
				Row:  space.Row,
				Col:  space.Col,
				Type: Space_UsedPod,
			}
		}
	}

	gameState.CurrentPlayer = getNextPlayerId(gameState.Players, actingPlayerIndex)
	return gameState, gameEvent, nil //TODO: End the game when no human players remain
}

func getNextPlayerId(players []Player, currentIndex int) string {
	nextIndex := currentIndex + 1
	if currentIndex >= len(players)-1 {
		nextIndex = 0
	}
	if players[nextIndex].Team == PlayerTeam_Spectator {
		return getNextPlayerId(players, nextIndex)
	} else {
		return players[nextIndex].Id
	}
}
