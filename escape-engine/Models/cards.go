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

type CardPlayDetails struct {
	TargetRow    string
	TargetCol    int
	TargetPlayer string
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

func (r *RedCard) Play(gameState *GameState, details CardPlayDetails) {

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

func (g *GreenCard) Play(gameState *GameState, details CardPlayDetails) {
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

func (a *Adrenaline) Play(gameState *GameState, details CardPlayDetails) {
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

func (m *Mutation) Play(gameState *GameState, details CardPlayDetails) {
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

func (t Teleport) Play(gameState *GameState, details CardPlayDetails) {
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

// #region Clone

type Clone struct {
	CardBase
}

func (c Clone) GetName() string {
	return c.Name
}

func (c Clone) GetType() string {
	return c.Type
}

func (c Clone) GetDescription() string {
	return c.Description
}

func (c Clone) Play(gameState *GameState, details CardPlayDetails) {
	activePlayer := gameState.GetCurrentPlayer()

	if indexOfEffect := slices.IndexFunc(activePlayer.StatusEffects, func(s StatusEffect) bool { return s.GetName() == "Cloned" }); indexOfEffect != -1 {
		activePlayer.StatusEffects[indexOfEffect].AddUse()
	} else {
		activePlayer.StatusEffects = append(activePlayer.StatusEffects, NewCloned())
	}
}

func NewClone() *Clone {
	return &Clone{
		CardBase: CardBase{
			Name:        "Clone",
			Description: "Creates a Clone of this player that automatically activates upon death",
			Type:        Card_White,
		},
	}
}

// #region Defence

type Defense struct {
	CardBase
}

func (c Defense) GetName() string {
	return c.Name
}

func (c Defense) GetType() string {
	return c.Type
}

func (c Defense) GetDescription() string {
	return c.Description
}

func (c Defense) Play(gameState *GameState, details CardPlayDetails) {
	activePlayer := gameState.GetCurrentPlayer()

	if indexOfEffect := slices.IndexFunc(activePlayer.StatusEffects, func(s StatusEffect) bool { return s.GetName() == "Armored" }); indexOfEffect != -1 {
		activePlayer.StatusEffects[indexOfEffect].AddUse()
	} else {
		activePlayer.StatusEffects = append(activePlayer.StatusEffects, NewArmored())
	}
}

func NewDefense() *Defense {
	return &Defense{
		CardBase: CardBase{
			Name:        "Defense",
			Description: "Makes this player invulnerable to the next attack that hits them",
			Type:        Card_White,
		},
	}
}

// #region Spotlight

type Spotlight struct {
	CardBase
}

func (c Spotlight) GetName() string {
	return c.Name
}

func (c Spotlight) GetType() string {
	return c.Type
}

func (c Spotlight) GetDescription() string {
	return c.Description
}

func (c Spotlight) Play(gameState *GameState, details CardPlayDetails) {
	if details.TargetRow == "" || details.TargetCol == -99 {
		panic("No details provided for spotlight")
	}

	//TODO: Implement

}

func NewSpotlight() *Spotlight {
	return &Spotlight{
		CardBase: CardBase{
			Name:        "Spotlight",
			Description: "Reveals any player in the chosen space or any adjacent space",
			Type:        Card_White,
		},
	}
}
