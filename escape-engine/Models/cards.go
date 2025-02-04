package Models

import (
	"math/rand"
	"slices"
)

const (
	Card_Red    = "Red"
	Card_Green  = "Green"
	Card_White  = "White"
	Card_NoCard = "SilentSector"
)

type CardBase struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Type        string `json:"type"`
}

// #region Red Card

type RedCard struct {
	CardBase
}

func (r *RedCard) GetName() string {
	return r.Name
}

func (r *RedCard) GetType() string {
	return r.Type
}

func (r *RedCard) GetDescription() string {
	return r.Description
}

func (r *RedCard) Play(gameState *GameState) {

}

func NewRedCard() *RedCard {
	return &RedCard{
		CardBase: CardBase{
			Name:        "Red Card",
			Description: "Make a noise in the sector you just moved into",
			Type:        Card_Red,
		},
	}
}

// #region Green Card

type GreenCard struct {
	CardBase
}

func (g *GreenCard) GetName() string {
	return g.Name
}

func (g *GreenCard) GetType() string {
	return g.Type
}

func (g *GreenCard) GetDescription() string {
	return g.Description
}

func (g *GreenCard) Play(gameState *GameState) {
}

func NewGreenCard() *GreenCard {
	return &GreenCard{
		CardBase: CardBase{
			Name:        "Green Card",
			Description: "Make a noise in any sector of your choosing",
			Type:        Card_Green,
		},
	}
}

// #region Adrenaline

type Adrenaline struct {
	CardBase
}

func (a *Adrenaline) GetName() string {
	return a.Name
}

func (a *Adrenaline) GetType() string {
	return a.Type
}

func (a *Adrenaline) GetDescription() string {
	return a.Description
}

func (a *Adrenaline) Play(gameState *GameState) {
	activePlayer := gameState.GetCurrentPlayer()

	if indexOfEffect := slices.IndexFunc(activePlayer.StatusEffects, func(s StatusEffect) bool { return s.GetName() == "Adrenaline Surge" }); indexOfEffect != -1 {
		activePlayer.StatusEffects[indexOfEffect].AddUse()
	} else {
		activePlayer.StatusEffects = append(activePlayer.StatusEffects, NewAdrenalineSurge())
	}
}

func NewAdrenaline() *Adrenaline {
	return &Adrenaline{
		CardBase: CardBase{
			Name:        "Adrenaline",
			Description: "Gives a rush of adrenaline, allowing one extra space of movement",
			Type:        Card_White,
		},
	}
}

// #region Mutation

type Mutation struct {
	CardBase
}

func (m *Mutation) GetName() string {
	return m.Name
}

func (m *Mutation) GetType() string {
	return m.Type
}

func (m *Mutation) GetDescription() string {
	return m.Description
}

func (m *Mutation) Play(gameState *GameState) {
	activePlayer := gameState.GetCurrentPlayer()

	activePlayer.Team = PlayerTeam_Alien
}

func NewMutation() *Mutation {
	return &Mutation{
		CardBase: CardBase{
			Name:        "Mutation",
			Description: "Turns the Player into an Alien!",
			Type:        Card_White,
		},
	}
}

// #region Teleport

type Teleport struct {
	CardBase
}

func (t Teleport) GetName() string {
	return t.Name
}

func (t Teleport) GetType() string {
	return t.Type
}

func (t Teleport) GetDescription() string {
	return t.Description
}

func (t Teleport) Play(gameState *GameState) {
	activePlayer := gameState.GetCurrentPlayer()

	humanStarts := gameState.GameMap.GetSpacesOfType(Space_HumanStart)

	toMoveTo := humanStarts[rand.Intn(len(humanStarts))]

	activePlayer.Row, activePlayer.Col = toMoveTo.Row, toMoveTo.Col
}

func NewTeleport() *Teleport {
	return &Teleport{
		CardBase: CardBase{
			Name:        "Teleport",
			Description: "Teleports the Player to a random Human Start Sector",
			Type:        Card_White,
		},
	}
}
