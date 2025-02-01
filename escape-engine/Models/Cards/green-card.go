package Cards

import "escape-engine/Models"

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

func (g *GreenCard) Play(gameState *Models.GameState) {

}
