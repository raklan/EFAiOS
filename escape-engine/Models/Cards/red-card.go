package Cards

import "escape-engine/Models"

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

func (r *RedCard) Play(gameState *Models.GameState) {

}
