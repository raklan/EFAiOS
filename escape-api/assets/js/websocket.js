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

const WS_CARD = "Card"
const WS_CLOSE = "Close"
const WS_ERROR = "Error"
const WS_GAMEEVENT = "GameEvent"
const WS_GAMEOVER = "GameOver"
const WS_GAMESTATE = "GameState"
const WS_LOBBYINFO = "LobbyInfo"
const WS_MOVEMENTRESPONSE = "MovementResponse"
const WS_TURNEND = "TurnEnd"
const WS_AVAILABLEMOVEMENT = "AvailableMovement"

function handleWsMessage(message) {
    console.info("Message inbound", message)
    let handler = null;
    switch (message.type) {
        case WS_CARD:
            handler = handleCardMessage;
            break;
        case WS_CLOSE:
            handler = handleCloseMessage;
            break;
        case WS_ERROR:
            handler = handleErrorMessage;
            break;
        case WS_GAMEEVENT:
            handler = handleGameEventMessage;
            break;
        case WS_GAMESTATE:
            handler = handleGameStateMessage;
            break;
        case WS_GAMEOVER:
            handler = handleGameOverMessage;
            break;
        case WS_LOBBYINFO:
            handler = handleLobbyInfoMessage;
            break;
        case WS_MOVEMENTRESPONSE:
            handler = handleMovementResponse;
            break;
        case WS_TURNEND:
            handler = handleTurnEnd;
            break;
        case WS_AVAILABLEMOVEMENT:
            handler = handleAvailableMovementMessage;
            break;
    }

    if (handler) {
        handler(message.data)
    }
}

async function handleCardMessage(cardEvent) {
    if (cardEvent.type == "White") {
        const actionToSend = {
            gameId: thisGameStateId,
            action: {
                type: 'Noise',
                turn: {
                    row: "!",
                    col: -99
                }
            }
        }
        sendWsMessage(ws, 'submitAction', actionToSend)
    } else if (cardEvent.type == "Green") {
        clickMode = ClickModes.Noise
        showPlayerChoicePopup('greenCard')
    } else if (cardEvent.type == "Red") {
        showPlayerChoicePopup('redCard')
    }

}

async function handleCloseMessage(messageData) {
    showNotification(messageData.message, 'Connection Lost')
    ws.close();
}

async function handleErrorMessage(socketError) {
    showNotification(socketError.message, 'Error')
}

async function handleGameEventMessage(gameEvent) {
    showNotification(gameEvent.description, 'Alert')
    let regRes = gameEvent.description.match(playerNameExtractor)

    if(regRes?.groups?.PlayerName){
        addEvent(regRes.groups.PlayerName, gameEvent.description)
    }
}

async function handleTurnEnd(turnEnd) {
    //For now, a TurnEnd should immediately end the turn
    if (isThisPlayersTurn) {
        const actionToSend = {
            gameId: thisGameStateId,
            action: {
                type: 'EndTurn',
                turn: {}
            }
        }
        sendWsMessage(ws, 'submitAction', actionToSend)
    }
}

async function handleGameOverMessage(messageData) {
    showNotification("The Game has ended", "Game Over!")
}

async function handleGameStateMessage(gameState) {
    thisPlayer = gameState.players?.find(p => p.id == thisPlayer.id)
    isThisPlayersTurn = gameState.currentPlayer == thisPlayer?.id
    drawMap(gameState.gameMap)
    if (!thisGameStateId) {
        thisGameStateId = gameState.id
        initializeEventLog(gameState.players)
    }
    document.getElementById("lobby").style.display = 'none';
    document.getElementById('gameplay').style.display = 'flex';

    document.querySelectorAll('.player').forEach(x => x.classList.remove('player'))
    if (thisPlayer.team != PlayerTeams.Spectator) {
        var playerSpace = document.getElementById(`hex-${thisPlayer.row}-${thisPlayer.col}`)
        playerSpace.classList.add('player')
    } else {
        for (let player of gameState.players.filter(p => p.team != PlayerTeams.Spectator)) {
            var playerSpace = document.getElementById(`hex-${player.row}-${player.col}`)
            playerSpace.classList.add('player')
        }
    }

    if (gameState.currentPlayer == thisPlayer.id) {
        clickMode = ClickModes.Moving
        sendWsMessage(ws, 'getAllowedMoves', {
            gameId: thisGameStateId
        })
    } else {
        clickMode = ClickModes.None
    }

    renderRoleCard();

    renderPlayerHand();
}

async function handleLobbyInfoMessage(messageData) {
    if (!thisPlayer) {
        thisPlayer = messageData.lobbyInfo?.players?.find(p => p.id == messageData.playerID)
        console.log(thisPlayer)
    }
    document.getElementById("lobby-roomCode").innerHTML = `Room Code: ${messageData.lobbyInfo.roomCode}`

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

    //#region Host Controls
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

        var gameConfig = document.getElementById("lobby-gameConfig")
        gameConfig.style.display = '';
    }
    //#endregion
}

async function handleMovementResponse(movementEvent) {
    thisPlayer.row = movementEvent.newRow;
    thisPlayer.col = movementEvent.newCol;

    //If needed, this can be moved to before updating thisPlayer.row and just search for that row and col instead of querySelectorAll
    document.querySelectorAll('.player').forEach(x => x.classList.remove('player'))

    var playerSpace = document.getElementById(`hex-${thisPlayer.row}-${thisPlayer.col}`)
    playerSpace.classList.add('player')

    //For now, just automatically don't let humans do anything after moving. In the future, we'll pause here to let them choose whether to play cards
    if (thisPlayer.team != PlayerTeams.Alien) {
        var actionToSend = {
            gameId: thisGameStateId,
            action: {
                type: 'Attack',
                turn: {
                    row: "!",
                    col: -99
                }
            }
        }

        sendWsMessage(ws, 'submitAction', actionToSend);
    } else if (thisPlayer.team == PlayerTeams.Alien) {
        showPlayerChoicePopup('attack')
    }
}

async function handleAvailableMovementMessage(availableMovement){
    availableMovement.spaces.forEach(space => {
        let spaceElement = document.getElementById(`hex-${space}`)
        if(spaceElement){
            spaceElement.classList = [cssClass, 'potential-move'].join(' ');
        }
    })
}