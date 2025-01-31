package CardConstructors

import "escape-engine/Models/Cards"

func NewMutation() *Cards.Mutation {
	return &Cards.Mutation{
		CardBase: Cards.CardBase{
			Name:        "Mutation",
			Description: "Turns the Player into an Alien!",
		},
	}
}

func NewAdrenaline() *Cards.Adrenaline {
	return &Cards.Adrenaline{
		CardBase: Cards.CardBase{
			Name:        "Adrenaline",
			Description: "Gives a rush of adrenaline, allowing one extra space of movement",
		},
	}
}

func NewTeleport() *Cards.Teleport {
	return &Cards.Teleport{
		CardBase: Cards.CardBase{
			Name:        "Teleport",
			Description: "Teleports the Player to a random Human Start Sector",
		},
	}
}
