package main

import (
	"encoding/json"
	"escape-engine/Models"
	CardConstructors "escape-engine/Models/Cards/Constructors"
	"fmt"
)

func main() {
	myStr := []byte("{\"id\":\"1738684162119FSTAFSG3U5\",\"gameMap\":{\"id\":\"1738684152520Y46PTK44EM\",\"name\":\"Test\",\"rows\":5,\"cols\":10,\"spaces\":{\"A-0\":{\"row\":\"A\",\"col\":0,\"type\":2},\"A-1\":{\"row\":\"A\",\"col\":1,\"type\":2},\"A-2\":{\"row\":\"A\",\"col\":2,\"type\":2},\"A-3\":{\"row\":\"A\",\"col\":3,\"type\":2},\"A-4\":{\"row\":\"A\",\"col\":4,\"type\":2},\"A-5\":{\"row\":\"A\",\"col\":5,\"type\":2},\"A-6\":{\"row\":\"A\",\"col\":6,\"type\":2},\"A-7\":{\"row\":\"A\",\"col\":7,\"type\":2},\"A-8\":{\"row\":\"A\",\"col\":8,\"type\":2},\"A-9\":{\"row\":\"A\",\"col\":9,\"type\":2},\"B-0\":{\"row\":\"B\",\"col\":0,\"type\":2},\"B-1\":{\"row\":\"B\",\"col\":1,\"type\":2},\"B-2\":{\"row\":\"B\",\"col\":2,\"type\":2},\"B-3\":{\"row\":\"B\",\"col\":3,\"type\":2},\"B-4\":{\"row\":\"B\",\"col\":4,\"type\":2},\"B-5\":{\"row\":\"B\",\"col\":5,\"type\":2},\"B-6\":{\"row\":\"B\",\"col\":6,\"type\":2},\"B-7\":{\"row\":\"B\",\"col\":7,\"type\":2},\"B-8\":{\"row\":\"B\",\"col\":8,\"type\":2},\"B-9\":{\"row\":\"B\",\"col\":9,\"type\":2},\"C-0\":{\"row\":\"C\",\"col\":0,\"type\":2},\"C-1\":{\"row\":\"C\",\"col\":1,\"type\":2},\"C-2\":{\"row\":\"C\",\"col\":2,\"type\":2},\"C-3\":{\"row\":\"C\",\"col\":3,\"type\":2},\"C-4\":{\"row\":\"C\",\"col\":4,\"type\":2},\"C-5\":{\"row\":\"C\",\"col\":5,\"type\":5},\"C-6\":{\"row\":\"C\",\"col\":6,\"type\":2},\"C-7\":{\"row\":\"C\",\"col\":7,\"type\":2},\"C-8\":{\"row\":\"C\",\"col\":8,\"type\":2},\"C-9\":{\"row\":\"C\",\"col\":9,\"type\":2},\"D-0\":{\"row\":\"D\",\"col\":0,\"type\":2},\"D-1\":{\"row\":\"D\",\"col\":1,\"type\":2},\"D-2\":{\"row\":\"D\",\"col\":2,\"type\":2},\"D-3\":{\"row\":\"D\",\"col\":3,\"type\":2},\"D-4\":{\"row\":\"D\",\"col\":4,\"type\":2},\"D-5\":{\"row\":\"D\",\"col\":5,\"type\":2},\"D-6\":{\"row\":\"D\",\"col\":6,\"type\":2},\"D-7\":{\"row\":\"D\",\"col\":7,\"type\":2},\"D-8\":{\"row\":\"D\",\"col\":8,\"type\":2},\"D-9\":{\"row\":\"D\",\"col\":9,\"type\":2},\"E-0\":{\"row\":\"E\",\"col\":0,\"type\":2},\"E-1\":{\"row\":\"E\",\"col\":1,\"type\":2},\"E-2\":{\"row\":\"E\",\"col\":2,\"type\":2},\"E-3\":{\"row\":\"E\",\"col\":3,\"type\":2},\"E-4\":{\"row\":\"E\",\"col\":4,\"type\":2},\"E-5\":{\"row\":\"E\",\"col\":5,\"type\":2},\"E-6\":{\"row\":\"E\",\"col\":6,\"type\":2},\"E-7\":{\"row\":\"E\",\"col\":7,\"type\":2},\"E-8\":{\"row\":\"E\",\"col\":8,\"type\":2},\"E-9\":{\"row\":\"E\",\"col\":9,\"type\":2},\"null-null\":{\"row\":\"\",\"col\":0,\"type\":0}}},\"gameConfig\":{\"numHumans\":1,\"numAliens\":0,\"numWorkingPods\":0,\"numBrokenPods\":0},\"deck\":[{\"name\":\"Red Card\",\"description\":\"Make a noise in the sector you just moved into\",\"type\":\"Red\"},{\"name\":\"Green Card\",\"description\":\"Make a noise in any sector of your choosing\",\"type\":\"Green\"},{\"name\":\"Teleport\",\"description\":\"Teleports the Player to a random Human Start Sector\",\"type\":\"White\"}],\"discardPile\":null,\"players\":[{\"id\":\"1738684159369CR2XL2PCOP\",\"name\":\"R\",\"team\":\"Human\",\"role\":\"\",\"statusEffects\":null,\"hand\":null,\"row\":\"C\",\"col\":5}],\"currentPlayer\":\"1738684159369CR2XL2PCOP\"}")

	test := GameStateTest{}

	json.Unmarshal(myStr, &test)

	fmt.Println(test)
}

type GameStateTest struct {
	//This is solely for book-keeping. The front end should submit this Id along with SubmittedActions to update the GameState
	Id string `json:"id"`
	//The map used by this Game
	GameMap Models.GameMap `json:"gameMap"`
	//GameState-specific config as defined by the Host
	GameConfig Models.GameConfig `json:"gameConfig"`
	//All cards used by this Game
	Deck []Models.Card `json:"deck"`
	//Used cards. Will be automatically reshuffled into the deck when empty
	DiscardPile []Models.Card `json:"discardPile"`
	//A list of the states of each Player in the game.
	Players []Models.Player `json:"players"`
	//Id of the Player whose turn it currently is
	CurrentPlayer string `json:"currentPlayer"`
}

func (gameState *GameStateTest) UnmarshalJSON(data []byte) error {
	intermediate := struct {
		//This is solely for book-keeping. The front end should submit this Id along with SubmittedActions to update the GameState
		Id string `json:"id"`
		//The map used by this Game
		GameMap Models.GameMap `json:"gameMap"`
		//GameState-specific config as defined by the Host
		GameConfig Models.GameConfig `json:"gameConfig"`
		//All cards used by this Game
		Deck []map[string]string `json:"deck"`
		//Used cards. Will be automatically reshuffled into the deck when empty
		DiscardPile []map[string]string `json:"discardPile"`
		//A list of the states of each Player in the game.
		Players []Models.Player `json:"players"`
		//Id of the Player whose turn it currently is
		CurrentPlayer string `json:"currentPlayer"`
	}{}

	err := json.Unmarshal(data, &intermediate)

	*gameState = GameStateTest{
		Id:            intermediate.Id,
		GameMap:       intermediate.GameMap,
		GameConfig:    intermediate.GameConfig,
		Players:       intermediate.Players,
		CurrentPlayer: intermediate.CurrentPlayer,
		Deck:          make([]Models.Card, len(intermediate.Deck)),
		DiscardPile:   make([]Models.Card, len(intermediate.DiscardPile)),
	}

	for i, card := range intermediate.Deck {
		switch card["name"] {
		case "Red Card":
			gameState.Deck[i] = CardConstructors.NewRedCard()
		case "Green Card":
			gameState.Deck[i] = CardConstructors.NewGreenCard()
		case "Mutation":
			gameState.Deck[i] = CardConstructors.NewMutation()
		case "Adrenaline":
			gameState.Deck[i] = CardConstructors.NewAdrenaline()
		case "Teleport":
			gameState.Deck[i] = CardConstructors.NewTeleport()
		}
	}

	for i, card := range intermediate.DiscardPile {
		switch card["name"] {
		case "Red Card":
			gameState.DiscardPile[i] = CardConstructors.NewRedCard()
		case "Green Card":
			gameState.DiscardPile[i] = CardConstructors.NewGreenCard()
		case "Mutation":
			gameState.DiscardPile[i] = CardConstructors.NewMutation()
		case "Adrenaline":
			gameState.DiscardPile[i] = CardConstructors.NewAdrenaline()
		case "Teleport":
			gameState.DiscardPile[i] = CardConstructors.NewTeleport()
		}
	}

	return err
}
