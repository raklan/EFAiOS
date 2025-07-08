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

// #region Hiding Spot
type HidingSpot struct {
	CardBase
}

func NewHidingSpot() *HidingSpot {
	return &HidingSpot{
		CardBase: CardBase{
			Name:        "Hiding Spot",
			Description: "Allows you to end your turn in the same spot you started it.",
			Type:        Card_White,
		},
	}
}

func (c HidingSpot) Play(gameState *GameState, details CardPlayDetails) string {
	activePlayer := gameState.GetCurrentPlayer()

	activePlayer.AddStatusEffect(StatusEffect_Lurking, NewLurking)

	go Recap.AddDataToRecap(gameState.Id, activePlayer.Id, gameState.Turn, fmt.Sprintf("Played %s, gaining one stack of %s", c.Name, StatusEffect_Lurking))

	return fmt.Sprintf("Player '%s' used a Hiding Spot! They can choose not to move at any time!", activePlayer.Name)
}

// #region Cloaking Device
type CloakingDevice struct {
	CardBase
}

func NewCloakingDevice() *CloakingDevice {
	return &CloakingDevice{
		CardBase: CardBase{
			Name:        "Cloaking Device",
			Description: "Hides you from the next Spotlight or Sensor that would reveal your location",
			Type:        Card_White,
		},
	}
}

func (c CloakingDevice) Play(gameState *GameState, details CardPlayDetails) string {
	activePlayer := gameState.GetCurrentPlayer()

	activePlayer.AddStatusEffect(StatusEffect_Invisible, NewInvisible)

	go Recap.AddDataToRecap(gameState.Id, activePlayer.Id, gameState.Turn, fmt.Sprintf("Played %s, gaining one stack of %s", c.Name, StatusEffect_Invisible))

	return fmt.Sprintf("Player '%s' used a Cloaking Device! They will be hidden from the next Spotlight or Sensor that would reveal their location!", activePlayer.Name)
}

// #region Engineering Manual
type EngineeringManual struct {
	CardBase
}

func NewEngineeringManual() *EngineeringManual {
	return &EngineeringManual{
		CardBase: CardBase{
			Name:        "Engineering Manual",
			Description: "Allows you to draw 2 Escape Pod cards from the next Escape Pod you try to use",
			Type:        Card_White,
		},
	}
}

func (c EngineeringManual) Play(gameState *GameState, details CardPlayDetails) string {
	activePlayer := gameState.GetCurrentPlayer()

	activePlayer.AddStatusEffect(StatusEffect_Knowhow, NewKnowhow)

	go Recap.AddDataToRecap(gameState.Id, activePlayer.Id, gameState.Turn, fmt.Sprintf("Played %s, gaining one stack of %s", c.Name, StatusEffect_Knowhow))

	return fmt.Sprintf("Player '%s' used an Engineering Manual! They can draw 2 Escape Pod cards from the next Escape Pod they enter!", activePlayer.Name)
}

// #region Noisemaker
type Noisemaker struct {
	CardBase
}

func NewNoisemaker() *Noisemaker {
	return &Noisemaker{
		CardBase: CardBase{
			Name:        "Noisemaker",
			Description: "Allows you to choose a sector to make noise in the next time you draw a White Card.",
			Type:        Card_White,
		},
	}
}

func (c Noisemaker) Play(gameState *GameState, details CardPlayDetails) string {
	activePlayer := gameState.GetCurrentPlayer()

	activePlayer.AddStatusEffect(StatusEffect_Deceptive, NewDeceptive)

	go Recap.AddDataToRecap(gameState.Id, activePlayer.Id, gameState.Turn, fmt.Sprintf("Played %s, gaining one stack of %s", c.Name, StatusEffect_Deceptive))

	return fmt.Sprintf("Player '%s' used a Noisemaker! They can choose a sector to make noise in the next time they draw a White Card!", activePlayer.Name)
}
