package Cards

import "escape-engine/Models"

type Mutation struct {
	CardBase
}

func (m Mutation) GetName() string {
	return m.Name
}

func (m Mutation) GetDescription() string {
	return m.Description
}

func (m Mutation) Play(gameState *Models.GameState) {
	activePlayer := gameState.GetCurrentPlayer()

	activePlayer.Team = Models.PlayerTeam_Alien
}
