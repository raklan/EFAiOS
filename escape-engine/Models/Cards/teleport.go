package Cards

import (
	"escape-engine/Models"
	"math/rand"
)

type Teleport struct {
	CardBase
}

func (t Teleport) GetName() string {
	return t.Name
}

func (t Teleport) GetDescription() string {
	return t.Description
}

func (t Teleport) Play(gameState *Models.GameState) {
	activePlayer := gameState.GetCurrentPlayer()

	humanStarts := gameState.GameMap.GetSpacesOfType(Models.Space_HumanStart)

	toMoveTo := humanStarts[rand.Intn(len(humanStarts))]

	activePlayer.Row, activePlayer.Col = toMoveTo.Row, toMoveTo.Col
}
