package Models

// The different types of messages the server might send to a client connected via websocket.
const (
	WebsocketMessage_Card             = "Card"
	WebsocketMessage_Close            = "Close"
	WebsocketMessage_Error            = "Error"
	WebsocketMessage_GameEvent        = "GameEvent"
	WebsocketMessage_GameOver         = "GameOver"
	WebsocketMessage_GameState        = "GameState" //In this case, the [Data] field will be a GameState struct
	WebsocketMessage_LobbyInfo        = "LobbyInfo"
	WebsocketMessage_MovementResponse = "MovementResponse"
	WebsocketMessage_TurnEnd          = "TurnEnd"
)

type WebsocketMessageListItem struct {
	Message         WebsocketMessage
	ShouldBroadcast bool
}

// A message sent from the server to a client. The frontend can check [Type] to determine how to parse the object in [Data]
type WebsocketMessage struct {
	//One of the above constants. That constant will tell you which of the below structs is found in the [Data] field
	Type string `json:"type"`
	//One of the below structs, a Changelog, or a GameState. Its exact type is recorded in [Type]
	Data any `json:"data"`
}

// A message containing a Player's assigned ID and the details of the lobby after they've joined it, whether by hosting it or joining a pre-existing lobby.
// The frontend should store this PlayerID.
type LobbyInfo struct {
	PlayerID  string `json:"playerID"`
	LobbyInfo Lobby  `json:"lobbyInfo"`
}

// If some message from a client causes any error, one of these is sent back to the client
type SocketError struct {
	Message string `json:"message"`
}

// If a connection is about to be closed by the server, it will send a SocketClose, followed by immediately closing the connection
type SocketClose struct {
	Message string `json:"message"`
}

type GameOver struct {
}

type TurnEnd struct {
}

type GameEvent struct {
	Description string `json:"description"`
	Row         string `json:"row"`
	Col         int    `json:"col"`
}

type CardEvent struct {
	Type string `json:"type"`
	Card any    `json:"card"` //This will be a Card object, but I can't call it that because it would cause an import cycle
}

type MovementEvent struct {
	NewRow string `json:"newRow"`
	NewCol int    `json:"newCol"`
}
