package CardConstructors

import (
	"escape-engine/Models"
	"escape-engine/Models/Cards"
)

func NewRedCard() *Cards.RedCard {
	return &Cards.RedCard{
		CardBase: Cards.CardBase{
			Name:        "Red Card",
			Description: "Make a noise in the sector you just moved into",
			Type:        Models.Card_Red,
		},
	}
}

func NewGreenCard() *Cards.GreenCard {
	return &Cards.GreenCard{
		CardBase: Cards.CardBase{
			Name:        "Green Card",
			Description: "Make a noise in any sector of your choosing",
			Type:        Models.Card_Green,
		},
	}
}

func NewMutation() *Cards.Mutation {
	return &Cards.Mutation{
		CardBase: Cards.CardBase{
			Name:        "Mutation",
			Description: "Turns the Player into an Alien!",
			Type:        Models.Card_White,
		},
	}
}

func NewAdrenaline() *Cards.Adrenaline {
	return &Cards.Adrenaline{
		CardBase: Cards.CardBase{
			Name:        "Adrenaline",
			Description: "Gives a rush of adrenaline, allowing one extra space of movement",
			Type:        Models.Card_White,
		},
	}
}

func NewTeleport() *Cards.Teleport {
	return &Cards.Teleport{
		CardBase: Cards.CardBase{
			Name:        "Teleport",
			Description: "Teleports the Player to a random Human Start Sector",
			Type:        Models.Card_White,
		},
	}
}
