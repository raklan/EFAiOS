package Engine

import (
	"encoding/json"
	"escape-engine/Models"
	"log"
	"net/http"
	"slices"
	"sync"

	"github.com/gorilla/websocket"
)

// The different types of messages the server might send to a client connected via websocket.
const (
	WebsocketMessage_Changelog = "Changelog"
	WebsocketMessage_Close     = "Close"
	WebsocketMessage_Error     = "Error"
	WebsocketMessage_GameOver  = "GameOver"
	WebsocketMessage_GameState = "GameState"
	WebsocketMessage_LobbyInfo = "LobbyInfo"
)

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
	PlayerID  string       `json:"playerID"`
	LobbyInfo Models.Lobby `json:"lobbyInfo"`
}

// I know SocketError and SocketClose have the same exact structure, but I've separated them for both clarity in the code and in case we end up wanting to put additional (unique) data in one or both

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

// Outer key is the room code, inner key is the playerId. Value is the actual connection object
var gamesClients = make(map[string]map[string]*websocket.Conn)

// Mutex to access gamesClients in a thread-safe manner
var gamesClientsMutex = sync.Mutex{}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  32768, // Setting read buffer size to 32 KB
	WriteBufferSize: 32768, // Setting write buffer size to 32 KB
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// HostLobby creates the waiting lobby, joins on behalf of the given player, and upgrades the host into a websocket.
// The lobby (which contains the Room Code used for other people to join) is then passed back into the websocket.
func HostLobby(w http.ResponseWriter, r *http.Request) {
	log.Println("Starting hostLobby")
	mapId := r.URL.Query().Get("mapId")
	playerName := r.URL.Query().Get("playerName")

	if mapId == "" || playerName == "" {
		log.Println("Missing mapId or playerName in request")
		http.Error(w, "Missing mapId or playerName in request", http.StatusBadRequest)
		return
	}

	roomCode, err := CreateRoom(mapId) // Assuming Engine.CreateRoom initializes room in DB
	if err != nil {
		http.Error(w, "Unable to create room", http.StatusInternalServerError)
		return
	}
	lobbyInfo, playerID, err := JoinRoom(roomCode, playerName)
	if err != nil {
		log.Printf("Error joining room: %v\n", err)
		http.Error(w, "Unable to join room", http.StatusInternalServerError)
		return
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "WebSocket upgrade failed", http.StatusInternalServerError)
		return
	}
	msg := WebsocketMessage{
		Type: WebsocketMessage_LobbyInfo,
		Data: LobbyInfo{
			PlayerID:  playerID,
			LobbyInfo: lobbyInfo,
		},
	}
	conn.WriteJSON(msg)
	gamesClientsMutex.Lock()
	if _, exists := gamesClients[roomCode]; !exists {
		gamesClients[roomCode] = make(map[string]*websocket.Conn)
	}
	gamesClients[roomCode][playerID] = conn
	gamesClientsMutex.Unlock()

	go manageClient(roomCode, gamesClients[roomCode], playerID, conn)
}

// Given a playerName and roomCode, tries to join that room for the given player. If joining was successul,
// the client's connection is upgraded to a websocket. Once complete, the client receives the lobby info
func HandleJoinLobby(w http.ResponseWriter, r *http.Request) {
	roomCode := r.URL.Query().Get("roomCode")
	playerName := r.URL.Query().Get("playerName")

	if roomCode == "" || playerName == "" {
		http.Error(w, "Please provide roomCode and playerName", http.StatusBadRequest)
		return
	}

	_, playerID, err := JoinRoom(roomCode, playerName)
	if err != nil {
		log.Printf("Error joining room: %v\n", err)
		http.Error(w, "Unable to join room", http.StatusNotFound)
		return
	}

	gamesClientsMutex.Lock()
	if _, exists := gamesClients[roomCode]; !exists {
		http.Error(w, "Lobby is not being tracked by server.", http.StatusInternalServerError)
		gamesClientsMutex.Unlock()
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "WebSocket upgrade failed", http.StatusInternalServerError)
		return
	}

	// msg := WebsocketMessage{
	// 	Type: WebsocketMessage_LobbyInfo,
	// 	Data: LobbyInfo{
	// 		PlayerID:  playerID,
	// 		LobbyInfo: lobbyInfo,
	// 	},
	// }
	//conn.WriteJSON(msg)
	gamesClients[roomCode][playerID] = conn
	gamesClientsMutex.Unlock()
	go handShake(roomCode, playerID)
}

// Given a roomCode and a playerId (which should have been generated by the backend upon Hosting or Joining that lobby), will attempt to
// reinsert the player into the Lobby. If successful, connection is upgraded to a websocket, then, depending on whether the game has started yet
// or not, either a LobbyInfo or GameState is sent to the player, after which they're added into the regular flow of listening for messages
func HandleRejoinLobby(w http.ResponseWriter, r *http.Request) {
	log.Println("Received request to rejoin a lobby!")

	roomCode := r.URL.Query().Get("roomCode")
	playerId := r.URL.Query().Get("playerId")

	if roomCode == "" || playerId == "" {
		http.Error(w, "Please send both roomCode and playerId", http.StatusBadRequest)
		return
	}

	lobbyInfo, err := LoadLobbyFromRedis(roomCode)
	if err != nil {
		http.Error(w, "Could not find requested lobby", http.StatusNotFound)
		return
	}

	//Make sure this player has joined the game before
	log.Printf("Making sure player {%s} has joined this game before", playerId)
	if !slices.ContainsFunc(lobbyInfo.Players, func(p Models.Player) bool { return p.Id == playerId }) {
		http.Error(w, "No player with given ID found in lobby", http.StatusNotFound)
		return
	}

	//Make sure that player does not already have an open connection
	log.Printf("Making sure player {%s} does not already have an open connection", playerId)
	gamesClientsMutex.Lock()
	if _, exists := gamesClients[roomCode][playerId]; exists {
		http.Error(w, "Found already open connection for player", http.StatusBadRequest)
		gamesClientsMutex.Unlock()
		return
	}
	gamesClientsMutex.Unlock()

	//Now we should know the player is allowed to rejoin. Upgrade to websocket
	log.Printf("Player {%s} is allowed to rejoin. Upgrading connection to websocket.", playerId)
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "WebSocket upgrade failed", http.StatusInternalServerError)
		return
	}

	log.Println("Upgrade finished. Checking status of lobby to give accurate first message(s)...")

	//If the game has ended, indicate as such and close connection
	if lobbyInfo.Status == Models.LobbyStatus_Ended {
		log.Println("Game has already ended! Player cannot join!")
		msg := WebsocketMessage{
			Type: WebsocketMessage_Error,
			Data: SocketError{
				Message: "Game has already ended. Cannot rejoin",
			},
		}
		conn.WriteJSON(msg)
		conn.Close()
		return
	}

	//Send LobbyInfo
	msg := WebsocketMessage{
		Type: WebsocketMessage_LobbyInfo,
		Data: LobbyInfo{
			PlayerID:  playerId,
			LobbyInfo: lobbyInfo,
		},
	}
	conn.WriteJSON(msg)

	//If the game has started, send a GameState
	if lobbyInfo.Status == Models.LobbyStatus_InProgress {
		log.Println("Game has been marked as 'In Progress' - Sending GameState...")
		gameState, err := GetCachedGameStateFromRedis(lobbyInfo.GameStateId)
		if err != nil {
			conn.WriteJSON(WebsocketMessage{
				Type: WebsocketMessage_Error,
				Data: SocketError{Message: err.Error()},
			})
		}
		conn.WriteJSON(WebsocketMessage{Type: WebsocketMessage_GameState, Data: gameState})
	}

	log.Printf("Player {%s} has been given first message(s). Beginning to track connection for further communication...", playerId)

	//Store the connection, give the connection its own goroutine, and begin listening for more messages
	gamesClientsMutex.Lock()
	gamesClients[roomCode][playerId] = conn
	room := gamesClients[roomCode]
	gamesClientsMutex.Unlock()

	go manageClient(roomCode, room, playerId, conn)
}

// handShake sends out the lobby info to everyone currently in the room, along with the
// name of the freshly joined player
func handShake(roomCode string, newPlayerId string) {
	gamesClientsMutex.Lock()
	room := gamesClients[roomCode]
	gamesClientsMutex.Unlock()
	jsonLobby, err := LoadLobbyFromRedis(roomCode)

	msg := WebsocketMessage{
		Type: WebsocketMessage_LobbyInfo,
		Data: LobbyInfo{
			PlayerID:  newPlayerId,
			LobbyInfo: jsonLobby,
		},
	}
	for playerId, conn := range room {
		if err != nil {
			log.Printf("error connecting: {%s}", err)
		}
		err := conn.WriteJSON(msg)
		if err != nil {
			log.Println("Error sending handshake, aborting connection ", playerId)
			gamesClientsMutex.Lock()
			room[playerId].Close()
			delete(room, playerId)
			gamesClientsMutex.Unlock()
			continue // I don't handle disconnection rn
		}
	}

	//Last thing we need to do is start listening for messages from this player
	go manageClient(roomCode, room, newPlayerId, room[newPlayerId])
}

func manageClient(lobbyCode string, room map[string]*websocket.Conn, playerId string, conn *websocket.Conn) {
	defer socketRecovery(room, playerId)

	log.Printf("Managing Connection for playerId %s and waiting for message", playerId)
	_, msg, err := conn.ReadMessage()
	if err != nil {
		//Disconnect client by closing connection and removing player from lobby
		gamesClientsMutex.Lock()
		room[playerId].Close()
		delete(room, playerId)
		gamesClientsMutex.Unlock()
	}

	processMessage(lobbyCode, playerId, msg)
}

// This gets called on loop per lobby
func processMessage(roomCode string, playerId string, message []byte) {

	var msg struct {
		JsonType string          `json:"jsonType"`
		Data     json.RawMessage `json:"data"` //Raw message delays the parsing
	}
	json.Unmarshal(message, &msg)

	gamesClientsMutex.Lock()
	room := gamesClients[roomCode]
	gamesClientsMutex.Unlock()

	switch msg.JsonType {
	case "startGame":
		log.Println("Received request to start game")
		config := Models.GameConfig{}
		if err := json.Unmarshal(msg.Data, &config); err != nil {
			log.Printf("error decoding startGame config: %s", err)
			socketError := WebsocketMessage{
				Type: WebsocketMessage_Error,
				Data: SocketError{
					Message: err.Error(),
				},
			}
			room[playerId].WriteJSON(socketError)
			break
		}

		game, err := GetInitialGameState(roomCode, config)
		if err != nil {
			log.Printf("ERROR: GAME NOT STARTED, ABORTING...%s", err)
			break
		}

		sendMessageToAllPlayers(room, WebsocketMessage{Type: WebsocketMessage_GameState, Data: game})
	case "endGame":
		err := EndGame(roomCode, playerId)
		if err != nil {
			log.Printf("ERROR: Trying to end game...%s", err)
			return
		}
		cleanUpRoom(room, roomCode)
	case "submitAction":
		var action struct {
			GameId string                 `json:"gameId"`
			Action Models.SubmittedAction `json:"action"`
		}
		if err := json.Unmarshal(msg.Data, &action); err != nil {
			log.Printf("error decoding submitAction: {%s}", err)
			socketError := WebsocketMessage{
				Type: WebsocketMessage_Error,
				Data: SocketError{
					Message: err.Error(),
				},
			}
			room[playerId].WriteJSON(socketError)
			break
		}

		//Supply PlayerId with the Id of the player belonging to this connection
		action.Action.PlayerId = playerId

		gameState, err := SubmitAction(action.GameId, action.Action)
		if err != nil {
			log.Printf("error with submitAction: {%s}", err)
			socketError := WebsocketMessage{
				Type: WebsocketMessage_Error,
				Data: SocketError{
					Message: err.Error(),
				},
			}
			room[playerId].WriteJSON(socketError)
			break
		}

		sendMessageToAllPlayers(room, WebsocketMessage{Type: WebsocketMessage_GameState, Data: gameState})
	case "leaveLobby":
		updatedLobby, err := endPlayerConnection(roomCode, playerId, room)

		if err != nil {
			socketError := WebsocketMessage{
				Type: WebsocketMessage_Error,
				Data: SocketError{
					Message: err.Error(),
				},
			}
			room[playerId].WriteJSON(socketError)
			break
		}

		sendMessageToAllPlayers(room, WebsocketMessage{Type: WebsocketMessage_LobbyInfo, Data: LobbyInfo{PlayerID: "", LobbyInfo: updatedLobby}})
	case "kickPlayer":
		var action struct {
			PlayerToKick string `json:"playerToKick"`
		}
		if err := json.Unmarshal(msg.Data, &action); err != nil {
			log.Printf("Error trying to unmarshal kick request into struct with field 'playerToKick' ... Please ensure field exists in Data")
			socketError := WebsocketMessage{
				Type: WebsocketMessage_Error,
				Data: SocketError{
					Message: "Message is malformed. Please ensure field 'playerToKick' is found in message object's 'Data' field!",
				},
			}
			room[playerId].WriteJSON(socketError)
			break
		}

		dbLobby, err := LoadLobbyFromRedis(roomCode)

		if err != nil {
			log.Printf("Error trying to find lobby")
			socketError := WebsocketMessage{
				Type: WebsocketMessage_Error,
				Data: SocketError{
					Message: "Could not find lobby. Something has gone terribly wrong",
				},
			}
			room[playerId].WriteJSON(socketError)
			break
		}

		if dbLobby.Host.Id != playerId {
			socketError := WebsocketMessage{
				Type: WebsocketMessage_Error,
				Data: SocketError{
					Message: "Player submitting kick request is not the host of the lobby!",
				},
			}
			room[playerId].WriteJSON(socketError)
			break
		}

		updatedLobby, err := endPlayerConnection(roomCode, action.PlayerToKick, room)

		if err != nil {
			socketError := WebsocketMessage{
				Type: WebsocketMessage_Error,
				Data: SocketError{
					Message: err.Error(),
				},
			}
			room[playerId].WriteJSON(socketError)
			break
		}

		sendMessageToAllPlayers(room, WebsocketMessage{Type: WebsocketMessage_LobbyInfo, Data: LobbyInfo{PlayerID: "", LobbyInfo: updatedLobby}})
	case "disconnect":
		log.Printf("Player %s is requesting a disconnect!", playerId)
		conn := room[playerId]
		msg := WebsocketMessage{
			Type: WebsocketMessage_Close,
			Data: SocketClose{
				Message: "Request acknowledged, closing connection",
			},
		}

		conn.WriteJSON(msg)
		log.Println("Close message has been sent, closing connection")
		conn.Close()
		delete(room, playerId)
		log.Println("Connection closed and server has stopped tracking websocket connection")
		return //Return so we don't go back into manageClient
	default:
		log.Println("Unknown type sent, ignoring message recieved", msg)
	}

	//Listen for the next message from this client. Not using `go manageClient(...)` because this is already happening in a goroutine
	manageClient(roomCode, room, playerId, room[playerId])
}

func endPlayerConnection(roomCode string, playerId string, room map[string]*websocket.Conn) (Models.Lobby, error) {
	//Tell the engine to remove the player from the DB copy of the lobby
	updatedLobby, err := LeaveRoom(roomCode, playerId)
	if err != nil {
		return Models.Lobby{}, err
	}

	//If the removed client has a currently open connection, tell that the client that the connection is closing, then close connection
	if conn, exists := room[playerId]; exists {
		msg := WebsocketMessage{
			Type: WebsocketMessage_Close,
			Data: SocketClose{
				Message: "Player has been removed from Lobby. Closing connection",
			},
		}
		conn.WriteJSON(msg)
		conn.Close()
		//Remove connection from lobby map so we don't try to send them any more messages
		delete(room, playerId)
	}

	return updatedLobby, nil
}

func cleanUpRoom(room map[string]*websocket.Conn, roomCode string) {
	gameOverMessage := WebsocketMessage{
		Type: WebsocketMessage_GameOver,
		Data: GameOver{}, //Maybe put the final GameState here?
	}
	closeMessage := WebsocketMessage{
		Type: WebsocketMessage_Close,
		Data: SocketClose{
			Message: "Game has ended. Closing connection",
		},
	}

	//Send the messages to every player
	for playerId, conn := range room {
		err := conn.WriteJSON(gameOverMessage)
		if err != nil {
			log.Printf("Error sending GameOver Message to %s. Aborting message, but closing connection anyways", playerId)
		}
		err = conn.WriteJSON(closeMessage)
		if err != nil {
			log.Printf("Error sending Close Message to %s. Aborting message, but closing connection anyways", playerId)
		}
		conn.Close()
		delete(room, playerId)
	}

	//Clean up the gamesClients map
	gamesClientsMutex.Lock()
	delete(gamesClients, roomCode)
	gamesClientsMutex.Unlock()
}

func sendMessageToAllPlayers(room map[string]*websocket.Conn, message WebsocketMessage) {
	//log.Printf("Sending message to every player: %s", message)

	if message.Type == "" {
		log.Println("WARNING: Websocket message being sent has no Type set! Frontend will likely not know how to handle the message!")
	}

	for playerId, conn := range room {
		err := conn.WriteJSON(message)
		if err != nil {
			log.Println("Error sending message, skipping meesage to ", playerId)
			continue // I don't handle disconnection rn
		}
	}
}

// Defer this function whenever you try to read from a socket. If ReadMessage panics, this will kick in. Note: This must be set up (deferred) **BEFORE** calling ReadMessage
func socketRecovery(room map[string]*websocket.Conn, playerId string) {
	if r := recover(); r != nil {
		log.Printf("Something went wrong trying to read from the connection of Player, likely due to an unexpected closing of the Websocket connection: {%s} -- %s", playerId, r)
		//Disconnect client by closing connection and removing player from lobby
		gamesClientsMutex.Lock()
		//lobby[playerId].Close() Assume the socket is already closed if we're here
		delete(room, playerId)
		gamesClientsMutex.Unlock()
	}
}
