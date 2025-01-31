package StatusEffects

import (
	"escape-engine/Models"
	"slices"
)

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

func (a *AdrenalineSurge) Activate(gameState *Models.GameState) {
	a.UsesLeft--

	if a.UsesLeft <= 0 {
		activePlayer := gameState.GetCurrentPlayer()
		activePlayer.StatusEffects = slices.DeleteFunc(activePlayer.StatusEffects, func(s Models.StatusEffect) bool { return s == a })
	}
}
