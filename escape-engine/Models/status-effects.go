package Models

import (
	"math/rand"
	"slices"
)

type StatusEffectBase struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	UsesLeft    int    `json:"usesLeft"`
}

// #region Adrenaline Surge

type AdrenalineSurge struct {
	StatusEffectBase
}

func NewAdrenalineSurge() *AdrenalineSurge {
	return &AdrenalineSurge{
		StatusEffectBase: StatusEffectBase{
			Name:        "Adrenaline Surge",
			Description: "The affected player may move 1 extra space",
			UsesLeft:    1,
		},
	}
}

func (a *AdrenalineSurge) GetName() string {
	return a.Name
}

func (a *AdrenalineSurge) GetDescription() string {
	return a.Description
}

func (a *AdrenalineSurge) GetUsesLeft() int {
	return a.UsesLeft
}

func (a *AdrenalineSurge) AddUse() int {
	a.UsesLeft++
	return a.GetUsesLeft()
}

func (a *AdrenalineSurge) Activate(gameState *GameState) {
	a.UsesLeft--
	activePlayer := gameState.GetCurrentPlayer()

	if a.UsesLeft <= 0 {
		activePlayer.StatusEffects = slices.DeleteFunc(activePlayer.StatusEffects, func(s StatusEffect) bool { return s == a })
	}
}

// #region Cat

// #region Cloned

type Cloned struct {
	StatusEffectBase
}

func NewCloned() *Cloned {
	return &Cloned{
		StatusEffectBase: StatusEffectBase{
			Name:        "Cloned",
			Description: "This player has a clone that will automatically activate upon death",
			UsesLeft:    1,
		},
	}
}

func (s *Cloned) GetName() string {
	return s.Name
}

func (s *Cloned) GetDescription() string {
	return s.Description
}

func (s *Cloned) GetUsesLeft() int {
	return s.UsesLeft
}

func (s *Cloned) AddUse() int {
	s.UsesLeft++
	return s.GetUsesLeft()
}

func (s *Cloned) Activate(gameState *GameState) {
	s.UsesLeft--
	activePlayer := gameState.GetCurrentPlayer()
	humanStarts := gameState.GameMap.GetSpacesOfType(Space_HumanStart)

	toMoveTo := humanStarts[rand.Intn(len(humanStarts))]

	activePlayer.Row, activePlayer.Col = toMoveTo.Row, toMoveTo.Col

	if s.UsesLeft <= 0 {
		activePlayer.StatusEffects = slices.DeleteFunc(activePlayer.StatusEffects, func(s2 StatusEffect) bool { return s2 == s })
	}
}

// #region Armored

type Armored struct {
	StatusEffectBase
}

func NewArmored() *Armored {
	return &Armored{
		StatusEffectBase: StatusEffectBase{
			Name:        "Armored",
			Description: "This player is defended from the next attack that hits them",
			UsesLeft:    1,
		},
	}
}

func (s *Armored) GetName() string {
	return s.Name
}

func (s *Armored) GetDescription() string {
	return s.Description
}

func (s *Armored) GetUsesLeft() int {
	return s.UsesLeft
}

func (s *Armored) AddUse() int {
	s.UsesLeft++
	return s.GetUsesLeft()
}

func (s *Armored) Activate(gameState *GameState) {
	s.UsesLeft--
	activePlayer := gameState.GetCurrentPlayer()

	if s.UsesLeft <= 0 {
		activePlayer.StatusEffects = slices.DeleteFunc(activePlayer.StatusEffects, func(s2 StatusEffect) bool { return s2 == s })
	}
}

// #region Hyperphagic

type Hyperphagic struct {
	StatusEffectBase
}

func NewHyperphagic() *Armored {
	return &Armored{
		StatusEffectBase: StatusEffectBase{
			Name:        "Hyperphagic",
			Description: "This Alien has fed on a human, gaining strength. But they want more...",
			UsesLeft:    1,
		},
	}
}

func (s *Hyperphagic) GetName() string {
	return s.Name
}

func (s *Hyperphagic) GetDescription() string {
	return s.Description
}

func (s *Hyperphagic) GetUsesLeft() int {
	return s.UsesLeft
}

func (s *Hyperphagic) AddUse() int {
	s.UsesLeft++
	return s.GetUsesLeft()
}

func (s *Hyperphagic) Activate(gameState *GameState) {
	//For now, don't subtract uses. It's a permanent bonus
	// s.UsesLeft--
	// activePlayer := gameState.GetCurrentPlayer()

	// if s.UsesLeft <= 0 {
	// 	activePlayer.StatusEffects = slices.DeleteFunc(activePlayer.StatusEffects, func(s2 StatusEffect) bool { return s2 == s })
	// }
}
