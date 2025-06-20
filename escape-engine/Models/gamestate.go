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
	Deck []Card `json:"deck"`
	//Used cards. Will be automatically reshuffled into the deck when empty
	DiscardPile []Card `json:"discardPile"`
	//A list of the states of each Player in the game.
	Players []Player `json:"players"`
	//Id of the Player whose turn it currently is
	CurrentPlayer string `json:"currentPlayer"`
	//Number of turn the game is on. The config sets the point at which this number will end the game
	Turn int `json:"turn"`
	//Priority list of StatusEffects
	StatusEffectPriorities map[string]int `json:"statusEffectPriorities"`
}

func (g *GameState) UnmarshalJSON(data []byte) error {
	intermediate := struct {
		Id                     string         `json:"id"`
		GameMap                GameMap        `json:"gameMap"`
		GameConfig             GameConfig     `json:"gameConfig"`
		Deck                   []CardBase     `json:"deck"`
		DiscardPile            []CardBase     `json:"discardPile"`
		Players                []Player       `json:"players"`
		CurrentPlayer          string         `json:"currentPlayer"`
		Turn                   int            `json:"turn"`
		StatusEffectPriorities map[string]int `json:"statusEffectPriorities"`
	}{}

	if err := json.Unmarshal(data, &intermediate); err != nil {
		return err
	}

	g.Id = intermediate.Id
	g.GameMap = intermediate.GameMap
	g.GameConfig = intermediate.GameConfig
	g.Players = intermediate.Players
	g.CurrentPlayer = intermediate.CurrentPlayer
	g.Turn = intermediate.Turn
	g.StatusEffectPriorities = intermediate.StatusEffectPriorities

	//Copy Deck
	g.Deck = GetUnmarshalledCardArray(intermediate.Deck)

	//Copy Discard Pile
	g.DiscardPile = GetUnmarshalledCardArray(intermediate.DiscardPile)

	return nil
}

func (gameState *GameState) GetCurrentPlayer() *Player {
	currentPlayerIndex := slices.IndexFunc(gameState.Players, func(p Player) bool { return gameState.CurrentPlayer == p.Id })
	if currentPlayerIndex == -1 {
		panic("Could not find current player in GameState's players!")
	}
	return &gameState.Players[currentPlayerIndex]
}

func (gameState *GameState) GetHumanPlayers() []Player {
	players := make([]Player, gameState.GameConfig.NumHumans)
	indexToAssign := 0
	for _, player := range gameState.Players {
		if player.Team == PlayerTeam_Human {
			players[indexToAssign] = player
			indexToAssign++
		}
	}
	return players
}

func (gameState *GameState) GetAlienPlayers() []Player {
	players := make([]Player, gameState.GameConfig.NumAliens)
	indexToAssign := 0
	for _, player := range gameState.Players {
		if player.Team == PlayerTeam_Alien {
			players[indexToAssign] = player
			indexToAssign++
		}
	}
	return players
}

func (gameState *GameState) GetSpectatorPlayers() []Player {
	players := []Player{}
	for _, player := range gameState.Players {
		if player.Team == PlayerTeam_Spectator {
			players = append(players, player)
		}
	}
	return players
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
	//Number of turns before the game should end
	NumTurns int `json:"numTurns"`
	//Which cards should be active, as well as how many of each
	ActiveCards map[string]int `json:"activeCards"`
	//Which roles should be active, as well as the maximum number allowed to be present. Should be >= that role's presence in RequiredRoles, if it's required
	ActiveRoles map[string]int `json:"activeRoles"`
	//Which roles should be guaranteed to be in the game, as well as the number of players that should have that role
	RequiredRoles map[string]int `json:"requiredRoles"`
	//Which StatusEffects should be active, as well as their priority
	ActiveStatusEffects map[string]int `json:"activeStatusEffects"`
}

type Player struct {
	Id            string         `json:"id"`
	Name          string         `json:"name"`
	Team          string         `json:"team"`
	Role          string         `json:"role"`
	StatusEffects []StatusEffect `json:"statusEffects"`
	Hand          []Card         `json:"hand"`
	Row           int            `json:"row"`
	Col           string         `json:"col"`
}

func (p *Player) UnmarshalJSON(data []byte) error {
	intermediate := struct {
		Id            string         `json:"id"`
		Name          string         `json:"name"`
		Team          string         `json:"team"`
		Role          string         `json:"role"`
		Hand          []CardBase     `json:"hand"`
		StatusEffects []StatusEffect `json:"statusEffects"`
		Row           int            `json:"row"`
		Col           string         `json:"col"`
	}{}

	if err := json.Unmarshal(data, &intermediate); err != nil {
		return err
	}

	//Copy the easy fields first
	p.Id = intermediate.Id
	p.Name = intermediate.Name
	p.Team = intermediate.Team
	p.Role = intermediate.Role
	p.Row = intermediate.Row
	p.Col = intermediate.Col
	p.StatusEffects = intermediate.StatusEffects

	//Copy hand
	p.Hand = GetUnmarshalledCardArray(intermediate.Hand)

	return nil
}

// Returns true if the player has a status effect with the given name
func (p Player) HasStatusEffect(name string) bool {
	return slices.ContainsFunc(p.StatusEffects, func(s StatusEffect) bool { return s.Name == name })
}

// Attempts to find a status effect on the player with the given name and subtracts a use if one is found. Returns a boolean indicating whether any status effect was found
func (player *Player) SubtractStatusEffect(name string) bool {
	if indexOfEffect := slices.IndexFunc(player.StatusEffects, func(s StatusEffect) bool { return s.Name == name }); indexOfEffect != -1 {
		player.StatusEffects[indexOfEffect].SubtractUse(player)
		return true
	}
	return false
}

func (player Player) GetSpaceMapKey() string {
	return GetMapKey(player.Row, player.Col)
}

type Card interface {
	GetName() string
	GetType() string
	GetDescription() string
	GetDestroyOnUse() bool
	SetDestroyOnUse(bool)
	Play(*GameState, CardPlayDetails) string
}
