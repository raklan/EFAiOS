package Cards

import (
	"escape-engine/Models"
	"escape-engine/Models/StatusEffects"
	"slices"
)

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

func (a *Adrenaline) Play(gameState *Models.GameState) {
	activePlayer := gameState.GetCurrentPlayer()

	if indexOfEffect := slices.IndexFunc(activePlayer.StatusEffects, func(s Models.StatusEffect) bool { return s.GetName() == "Adrenaline Surge" }); indexOfEffect != -1 {
		activePlayer.StatusEffects[indexOfEffect].AddUse()
	} else {
		activePlayer.StatusEffects = append(activePlayer.StatusEffects, StatusEffects.NewAdrenalineSurge())
	}
}
