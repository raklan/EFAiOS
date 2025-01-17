package Models

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
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

func (move Movement) Execute(gameState *GameState, gameMap GameMap, playerId string) error {
	if gameState.CurrentPlayer != playerId {
		return fmt.Errorf("player trying to execute turn is not current player")
	}

	actingPlayerIndex := slices.IndexFunc(gameState.Players, func(p Player) bool { return playerId == p.Id })
	if actingPlayerIndex == -1 {
		return fmt.Errorf("could not find acting player with ID == {%s}", playerId)
	}

	actingPlayer := &(gameState.Players[actingPlayerIndex])

	//Bounds check
	spaceKey := fmt.Sprintf("%d,%d", move.ToRow, move.ToCol)

	if space, exists := gameMap.Spaces[spaceKey]; exists {
		//Make sure it's an open space
		cantMoveInto := []int{
			Space_AlienStart,
			Space_HumanStart,
			Space_Wall,
		}
		if slices.ContainsFunc(cantMoveInto, func(t int) bool { return space.Type == t }) {
			return fmt.Errorf("space [%d,%d] is not allowed to be moved into", move.ToRow, move.ToCol)
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
		actingPlayer.Row = move.ToRow
		actingPlayer.Col = move.ToCol
	} else {
		return fmt.Errorf("requested space [%d,%d] not found in map", move.ToRow, move.ToCol)
	}

	return nil
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

func (attack Attack) Execute(gameState *GameState, gameMap GameMap, playerId string) error {
	return nil
}
