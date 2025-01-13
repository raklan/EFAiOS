package Models

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
	//The ID of the map that this game state uses
	MapId string `json:"mapId"`
	//A list of the states of each Player in the game.
	Players []Player `json:"players"`
	//Id of the Player whose turn it currently is
	CurrentPlayer string `json:"currentPlayer"`
}

type Player struct {
	Id   string `json:"id"`
	Name string `json:"name"`
	Team string `json:"team"`
	Role string `json:"role"`
	Row  int    `json:"row"`
	Col  int    `json:"col"`
}
