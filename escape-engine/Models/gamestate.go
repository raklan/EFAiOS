package Models

import (
	"encoding/json"
	"slices"
)

const (
	LobbyStatus_AwaitingStart = "Awaiting Start"
	LobbyStatus_InProgress    = "In Progress"
	LobbyStatus_Ended         = "Game Ended"

	PlayerTeam_Human     = "Human"
	PlayerTeam_Alien     = "Alien"
	PlayerTeam_Spectator = "Spectator"
)

// A lobby is a collection of players waiting for a game to start. This is created by the /createRoom endpoint
// which returns the room code. From there, you can pass the room code to /joinRoom which will put you in the room
// and return the state of the lobby
type Lobby struct {
	//The Room Code used for players to join this lobby. Passed to /joinRoom
	RoomCode string `json:"roomCode"`
	//The Map this Lobby is going to play
	MapId string `json:"mapId"`
	//The GameState created from this lobby. This field is empty until the game is started,
	//at which point the API server will fill it in.
	GameStateId string `json:"gameStateId"`
	//The status of the game, really only used for the backend to determine whether a game has started/ended. Will be one of the above constants
	Status string `json:"status"`
	//Current number of players in the lobby
	NumPlayers int `json:"numPlayers"`
	//Maximum allowed players in the lobby.
	MaxPlayers int `json:"maxPlayers"`
	//Current list of joined players
	Players []Player `json:"players"`
	//The player that created the Lobby
	Host Player `json:"host"`
}

// This is intended to be the actual data the backend sends to the frontend to have it render things for the players. This is separate from the Game
// definition found throughout the other packages
type GameState struct {
	//This is solely for book-keeping. The front end should submit this Id along with SubmittedActions to update the GameState
	Id string `json:"id"`
	//The map used by this Game
	GameMap GameMap `json:"gameMap"`
	//GameState-specific config as defined by the Host
	GameConfig GameConfig `json:"gameConfig"`
	//All cards used by this Game
	Deck CardCollection `json:"deck"`
	//Used cards. Will be automatically reshuffled into the deck when empty
	DiscardPile CardCollection `json:"discardPile"`
	//A list of the states of each Player in the game.
	Players []Player `json:"players"`
	//Id of the Player whose turn it currently is
	CurrentPlayer string `json:"currentPlayer"`
	//Priority list of StatusEffects
	StatusEffectPriorities map[string]int
}

func (gameState *GameState) GetCurrentPlayer() *Player {
	currentPlayerIndex := slices.IndexFunc(gameState.Players, func(p Player) bool { return gameState.CurrentPlayer == p.Id })
	if currentPlayerIndex == -1 {
		panic("Could not find current player in GameState's players!")
	}
	return &gameState.Players[currentPlayerIndex]
}

// GameState-specific config as defined by the Host
type GameConfig struct {
	//Number of Humans currently in the Game. The Game automatically ends when this number hits 0.
	NumHumans int `json:"numHumans"`
	//Number of Aliens currently in the Game
	NumAliens int `json:"numAliens"`
	//Number of Working Escape Pods left. The Game automatically ends when this number hits 0.
	NumWorkingPods int `json:"numWorkingPods"`
	//Number of Broken Escape Pods left
	NumBrokenPods int `json:"numBrokenPods"`
	//Which cards should be active, as well as how many of each
	ActiveCards map[string]int `json:"activeCards"`
	//Which StatusEffects should be active, as well as their priority
	ActiveStatusEffects map[string]int `json:"activeStatusEffects"`
}

type Player struct {
	Id            string         `json:"id"`
	Name          string         `json:"name"`
	Team          string         `json:"team"`
	Role          string         `json:"role"`
	StatusEffects []StatusEffect `json:"statusEffects"`
	Hand          CardCollection `json:"hand"`
	Row           string         `json:"row"`
	Col           int            `json:"col"`
}

func (p Player) HasStatusEffect(name string) bool {
	return slices.ContainsFunc(p.StatusEffects, func(s StatusEffect) bool { return s.GetName() == name })
}

func (player *Player) SubtractStatusEffect(name string) bool {
	if indexOfEffect := slices.IndexFunc(player.StatusEffects, func(s StatusEffect) bool { return s.GetName() == NewAdrenaline().GetName() }); indexOfEffect != -1 {
		player.StatusEffects[indexOfEffect].SubtractUse(player)
		return true
	}
	return false
}

type CardCollection struct {
	Cards []Card `json:"cards"`
}

func (c *CardCollection) UnmarshalJSON(data []byte) error {
	intermediate := struct {
		Cards []map[string]string
	}{}

	if err := json.Unmarshal(data, &intermediate); err != nil {
		return err
	}

	c.Cards = make([]Card, len(intermediate.Cards))

	for i, card := range intermediate.Cards {
		switch card["name"] {
		case "Red Card":
			c.Cards[i] = NewRedCard()
		case "Green Card":
			c.Cards[i] = NewGreenCard()
		case "Mutation":
			c.Cards[i] = NewMutation()
		case "Adrenaline":
			c.Cards[i] = NewAdrenaline()
		case "Teleport":
			c.Cards[i] = NewTeleport()
		case "Clone":
			c.Cards[i] = NewClone()
		case "Defence":
			c.Cards[i] = NewDefense()
		case "Spotlight":
			c.Cards[i] = NewSpotlight()
		case "Attack":
			c.Cards[i] = NewAttackCard()
		}
	}

	return nil
}

type Card interface {
	GetName() string
	GetType() string
	GetDescription() string
	Play(*GameState, CardPlayDetails) string
}

type StatusEffect interface {
	//Returns the Name of this StatusEffect
	GetName() string
	//Returns the Description of this StatusEffect
	GetDescription() string
	//Returns the number of "Uses" left that this effect can activate. This StatusEffect should be removed from the Player's StatusEffects array when it hits 0
	GetUsesLeft() int
	//Adds one use to this StatusEffect
	AddUse() int
	//Uses this StatusEffect, reducing the UsesLeft by 1. If there are zero uses left, it will remove itself from the Player's StatusEffects. Returns a boolean describing whether there are uses left in this StatusEffect
	SubtractUse(*Player) bool
}
