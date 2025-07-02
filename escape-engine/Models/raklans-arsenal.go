package Models

import (
	"escape-engine/Models/Recap"
	"fmt"
	"math/rand"
)

//#region Unstable Teleporter

type UnstableTeleporter struct {
	CardBase
}

func NewUnstableTeleporter() *UnstableTeleporter {
	return &UnstableTeleporter{
		CardBase: CardBase{
			Name:        "Unstable Teleporter",
			Description: "Swaps the position of 2 randomly selected players",
			Type:        Card_White,
		},
	}
}

func (c UnstableTeleporter) Play(gameState *GameState, details CardPlayDetails) string {
	activePlayer := gameState.GetCurrentPlayer()

	player1 := gameState.Players[rand.Intn(len(gameState.Players))]
	player2 := gameState.Players[rand.Intn(len(gameState.Players))]

	for player1.Id == player2.Id && len(gameState.Players) > 1 {
		player2 = gameState.Players[rand.Intn(len(gameState.Players))]
	}

	newCol, newRow := player1.Col, player1.Row

	player1.Col, player1.Row = player2.Col, player2.Row
	player2.Col, player2.Row = newCol, newRow

	go Recap.AddDataToRecap(gameState.Id, activePlayer.Id, gameState.Turn, fmt.Sprintf("Played %s, switching Players '%s' and '%s'", c.Name, player1.Name, player2.Name))
	go Recap.AddDataToRecap(gameState.Id, player1.Id, gameState.Turn, fmt.Sprintf("Targeted by %s, swapped positions with Player '%s' ending up in Sector [%s-%d]", c.Name, player2.Name, player1.Col, player1.Row))
	go Recap.AddDataToRecap(gameState.Id, player2.Id, gameState.Turn, fmt.Sprintf("Targeted by %s, swapped positions with Player '%s' ending up in Sector [%s-%d]", c.Name, player1.Name, player2.Col, player2.Row))

	return fmt.Sprintf("Player '%s' used an Unstable Teleporter! Player '%s' and Player '%s' switched places!", activePlayer.Name, player1.Name, player2.Name)
}
