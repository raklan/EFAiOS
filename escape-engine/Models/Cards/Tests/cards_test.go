package CardTests

import (
	"escape-engine/Models"
	CardConstructors "escape-engine/Models/Cards/Constructors"
	"testing"
)

func TestAdrenaline(t *testing.T) {
	player := Models.Player{
		StatusEffects: []Models.StatusEffect{},
		Id:            "testPlayer",
	}
	gameState := Models.GameState{
		Players:       []Models.Player{player},
		CurrentPlayer: player.Id,
	}

	card := CardConstructors.NewAdrenaline()

	card.Play(&gameState)

	if len(gameState.GetCurrentPlayer().StatusEffects) != 1 {
		t.Fatal("Player does not have 1 status effect")
	}
}

func TestMutation(t *testing.T) {
	player := Models.Player{
		StatusEffects: []Models.StatusEffect{},
		Id:            "testPlayer",
		Team:          Models.PlayerTeam_Human,
	}
	gameState := Models.GameState{
		Players:       []Models.Player{player},
		CurrentPlayer: player.Id,
	}

	card := CardConstructors.NewMutation()

	card.Play(&gameState)

	if len(gameState.GetCurrentPlayer().StatusEffects) != 0 {
		t.Fatal("Player does not have 0 status effects")
	}

	if gameState.GetCurrentPlayer().Team != Models.PlayerTeam_Alien {
		t.Fatal("Player is not an alien!")
	}
}

func TestTeleport(t *testing.T) {
	player := Models.Player{
		StatusEffects: []Models.StatusEffect{},
		Id:            "testPlayer",
		Row:           "A",
		Col:           0,
	}
	gameState := Models.GameState{
		GameMap: Models.GameMap{
			Spaces: map[string]Models.Space{
				"A-0": {
					Row:  "A",
					Col:  0,
					Type: Models.Space_Dangerous,
				},
				"A-1": {
					Row:  "A",
					Col:  1,
					Type: Models.Space_HumanStart,
				},
			},
		},
		Players:       []Models.Player{player},
		CurrentPlayer: player.Id,
	}

	card := CardConstructors.NewTeleport()

	card.Play(&gameState)

	if len(gameState.GetCurrentPlayer().StatusEffects) != 0 {
		t.Fatal("Player does not have 0 status effects")
	}

	if gameState.GetCurrentPlayer().Row != "A" || gameState.GetCurrentPlayer().Col != 1 {
		t.Fatal("Player is not at human start")
	}
}
