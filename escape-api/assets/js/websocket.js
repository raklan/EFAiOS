function sendWsMessage(ws, msgType, data) {
    console.info("sending message", { jsonType: msgType, data: data })
    if (!ws || !ws?.OPEN) {
        console.assert(ws, 'WebSocket has not been initialized')
        console.assert(ws?.OPEN, 'WebSocket connection is not open')
        return;
    }
    ws.send(JSON.stringify({
        jsonType: msgType,
        data: data
    }))
}

const WS_CLOSE = "Close"
const WS_ERROR = "Error"
const WS_GAMESTATE = "GameState"
const WS_GAMEOVER = "GameOver"
const WS_LOBBYINFO = "LobbyInfo"

function handleWsMessage(message) {
    console.info("Message inbound", message)
    switch (message.type) {
        case WS_CLOSE:
            handleCloseMessage(message.data);
            break;
        case WS_ERROR:
            handleErrorMessage(message.data);
            break;
        case WS_GAMESTATE:
            handleGameStateMessage(message.data);
            break;
        case WS_GAMEOVER:
            handleGameOverMessage(message.data);
            break;
        case WS_LOBBYINFO:
            handleLobbyInfoMessage(message.data);
            break;
    }
}

async function handleCloseMessage(messageData) {

}

async function handleErrorMessage(messageData) {

}

async function handleGameOverMessage(messageData) {

}

async function handleGameStateMessage(gameState) {
    thisPlayer = gameState.players?.find(p => p.id == thisPlayer.id)
    drawMap(gameState.gameMap)
    if (!thisGameStateId) {
        thisGameStateId = gameState.id
    }
    document.getElementById("lobby").style.display = 'none';
    document.getElementById('gameplay').style.display = 'flex';

    document.querySelectorAll('.player').forEach(x => x.classList.remove('player'))
    for (let player of gameState.players.filter(p => p.team != PlayerTeams.Spectator)) {
        var playerSpace = document.getElementById(`hex-${player.row}-${player.col}`)
        playerSpace.classList = 'hexfield player'
    }
}

async function handleLobbyInfoMessage(messageData) {
    if (!thisPlayer) {
        thisPlayer = messageData.lobbyInfo?.players?.find(p => p.id == messageData.playerID)
        console.log(thisPlayer)
    }
    document.getElementById("lobby-roomCode").innerText = `Room Code: ${messageData.lobbyInfo.roomCode}`

    //#region Player List Rendering
    var playerList = document.getElementById("lobby-playerList")
    playerList.replaceChildren() //Important: Clear the player list so new players joining don't cause duplicate rendering

    for (let player of messageData.lobbyInfo.players) {
        playerEntry = document.createElement("div")
        playerEntry.innerText = player.name
        playerEntry.style.border = "1px solid black"

        playerList.appendChild(playerEntry)
    }
    //#endregion

    //#region Start Game Button
    if (thisPlayer?.id?.length > 0 && thisPlayer.id == messageData.lobbyInfo?.host?.id) {
        var startButton = document.getElementById("lobby-startButton")
        startButton.style.display = '';
        startButton.onclick = () => {
            config = {
                numHumans: 0,
                numAliens: 0,
                numWorkingPods: 0,
                numBrokenPods: 0
            }
            config.numHumans = parseInt(document.getElementById("config-numHumans")?.value ?? 0)
            config.numAliens = parseInt(document.getElementById("config-numAliens")?.value ?? 0)
            config.numWorkingPods = parseInt(document.getElementById("config-numWorkingPods")?.value ?? 0)
            config.numBrokenPods = parseInt(document.getElementById("config-numBrokenPods")?.value ?? 0)
            sendWsMessage(ws, 'startGame', config)
        }


    }
    //#endregion
} 