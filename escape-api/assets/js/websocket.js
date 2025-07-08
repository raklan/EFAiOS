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
                    row: -99,
                    col: "!"
                }
            }
        }
        sendWsMessage(ws, 'submitAction', actionToSend)
    } else if (cardEvent.type == "Green") {
        if (thisPlayer.statusEffects?.some(se => se.name === "Feline")) {
            clickMode = ClickModes.CatGreen
            showPlayerChoicePopup('cat-green')
        } else {
            clickMode = ClickModes.Noise
            showPlayerChoicePopup('greenCard')
        }
    } else if (cardEvent.type == "Red") {
        if (thisPlayer.statusEffects?.some(se => se.name === "Feline")) {
            clickMode = ClickModes.CatRed
            showPlayerChoicePopup('cat-red')
        } else {
            showPlayerChoicePopup('redCard')
        }
    }

}

async function handleCloseMessage(messageData) {
    if (!gameHasEnded) {
        showNotification(messageData.message, 'Connection Lost')
    }
    ws.close();
}

async function handleErrorMessage(socketError) {
    showNotification(socketError.message, 'Error')
}

async function handleGameEventMessage(gameEvent) {
    showNotification(gameEvent.description, 'Alert')
    let matches = [...gameEvent.description.matchAll(playerNameExtractor)]
    let playersMentionedInThisEvent = []

    if (matches?.length > 0) {
        for (let match of matches) {
            if (!playersMentionedInThisEvent.includes(match.groups.PlayerName)) { //Only add one entry if a player is mentioned multiple times
                addEvent(currentTurn, match.groups.PlayerName, gameEvent.description)
                saveEventToLocalStorage(match.groups.PlayerName, gameEvent.description)
                playersMentionedInThisEvent.push(match.groups.PlayerName)
            }
        }
    }
}

async function handleTurnEnd(turnEnd) {

    if (isThisPlayersTurn) {
        thisPlayer = turnEnd.playerCurrentState;
        renderPlayerHand();
        renderStatusEffects();
        clickMode = ClickModes.None;
        document.getElementById("endTurn-button").style.display = ''
    }
}

async function handleGameOverMessage(messageData) {
    gameHasEnded = true;
    showGameOver();
    let roomCode = JSON.parse(window.localStorage.getItem("efaios-connectionInfo")).roomCode
    document.getElementById('gameover-recap-link').setAttribute('href', `/recap?roomCode=${roomCode}`)
    window.localStorage.removeItem('efaios-connectionInfo')
    window.localStorage.removeItem('efaios-eventlog')
}

async function handleGameStateMessage(gameState) {
    thisPlayer = gameState.players?.find(p => p.id == thisPlayer.id)
    isThisPlayersTurn = gameState.currentPlayer == thisPlayer?.id
    gamePlayerList = gameState.players.filter(player => player.team != PlayerTeams.Spectator).map(player => {
        return {
            id: player.id,
            name: player.name,
            isThisPlayersTurn: gameState.currentPlayer == player.id
        }
    })
    
    if (!thisGameStateId) {
        document.getElementById("lobby").style.display = 'none';
        document.getElementById('gameplay').style.display = 'flex';
        thisGameStateId = gameState.id
        initializeEventLog(gameState.players)
        initializeCanvas();
        await fetch(`/api/role?name=${thisPlayer.role}`).then(resp => resp.json()).then(apiObj => {
            roleDescription = apiObj.roleDescription;
        })

        window.onresize = _ => {
            drawMap(gameState.gameMap)
            resizeCanvasToDisplaySize(document.getElementById("can"), document.getElementById("gridParent"))
            var playerSpace = document.getElementById(`hex-${thisPlayer.col}-${thisPlayer.row}`)
            playerSpace.classList.add('player')
        }
    }
    drawMap(gameState.gameMap)

    document.querySelectorAll('.player').forEach(x => x.classList.remove('player'))
    document.querySelectorAll('.player-human').forEach(x => x.classList.remove('player-human'))
    document.querySelectorAll('.player-alien').forEach(x => x.classList.remove('player-alien'))
    if (thisPlayer.team != PlayerTeams.Spectator) {
        var playerSpace = document.getElementById(`hex-${thisPlayer.col}-${thisPlayer.row}`)
        playerSpace.setAttribute('tooltip-text', 'You')
        playerSpace.setAttribute('tooltip-color', `var(--space-player)`)
        playerSpace.onmousemove = (event) => showSpaceTooltip(event)
        playerSpace.onmouseleave = (event) => hideSpaceTooltip(event)
        playerSpace.classList.add('player')
    } else {
        renderSpectatorView(gameState)
    }

    if (isThisPlayersTurn && !playerHasMoved) {
        if(showYourTurnNotification){
            showNotification('Your Turn', 'Your Turn');
        }else{
            showYourTurnNotification = true;
        }
        clickMode = ClickModes.Moving;
        sendWsMessage(ws, 'getAllowedMoves', {
            gameId: thisGameStateId
        })
    } else if (!isThisPlayersTurn) {
        clickMode = ClickModes.None
    }

    renderTeamCard();
    renderRoleCard();
    renderStatusEffects();
    renderPlayerHand();
    renderTurnOrder();
    renderTurnNumber(gameState.turn, gameState.gameMap.gameConfig.numTurns, gameState.gameMap.name);
}

async function handleLobbyInfoMessage(messageData) {
    if (!thisPlayer) {
        thisPlayer = messageData.lobbyInfo?.players?.find(p => p.id == messageData.playerID)
        const connectionInfo = {
            type: 'rejoin',
            playerId: thisPlayer.id,
            roomCode: messageData.lobbyInfo.roomCode
        }
        window.localStorage.setItem('efaios-connectionInfo', JSON.stringify(connectionInfo))
        setConfigFormFromObject(messageData.lobbyInfo.mapConfig)
    }
    document.getElementById("lobby-roomCode").innerHTML = `Room Code: ${messageData.lobbyInfo.roomCode}`
    document.getElementById('lobby-mapTitle').innerText = `Map: ${messageData.lobbyInfo.mapName}`

    //Player List Rendering
    var playerList = document.getElementById("lobby-playerList")
    playerList.replaceChildren() //Important: Clear the player list so new players joining don't cause duplicate rendering

    for (let player of messageData.lobbyInfo.players) {
        playerEntry = document.createElement("div")
        playerEntry.innerText = `${player.name}${player.name === messageData.lobbyInfo.host.name ? " (Host)" : ""}`
        playerEntry.style.border = "1px solid black"
        playerEntry.style.margin = '5px'

        playerList.appendChild(playerEntry)
    }

    //Host Controls
    if (thisPlayer?.id?.length > 0 && thisPlayer.id == messageData.lobbyInfo?.host?.id) {
        var startButton = document.getElementById("lobby-startButton")
        startButton.style.display = '';
        startButton.onclick = () => {
            let gameConfig = getGameConfig()
            console.info('starting game with config', gameConfig);
            if (gameConfig.numHumans + gameConfig.numAliens != messageData.lobbyInfo.players.length) {
                showNotification("# of Humans + # of Aliens must add up to # of Players in lobby!", "Error")
                return;
            }
            sendWsMessage(ws, 'startGame', getGameConfig())
        }

        var configButton = document.getElementById("lobby-gameConfigButton");
        configButton.style.display = '';
        configButton.onclick = () => {
            showConfig();
        }
    }
}

async function handleMovementResponse(movementEvent) {
    clickMode = ClickModes.None
    thisPlayer.row = movementEvent.newRow;
    thisPlayer.col = movementEvent.newCol;

    //If needed, this can be moved to before updating thisPlayer.row and just search for that row and col instead of querySelectorAll
    document.querySelectorAll('.player').forEach(x => x.classList.remove('player'))
    document.querySelectorAll('.hexfield.potential-move').forEach(x => x.classList.remove('potential-move'))

    var playerSpace = document.getElementById(`hex-${thisPlayer.col}-${thisPlayer.row}`)
    playerSpace.classList.add('player')

    //For now, just automatically don't let humans do anything after moving. In the future, we'll pause here to let them choose whether to play cards
    if (thisPlayer.team != PlayerTeams.Alien) {
        var actionToSend = {
            gameId: thisGameStateId,
            action: {
                type: 'Attack',
                turn: {
                    row: -99,
                    col: "!"
                }
            }
        }

        sendWsMessage(ws, 'submitAction', actionToSend);
    } else if (thisPlayer.team == PlayerTeams.Alien) {
        showPlayerChoicePopup('attack')
    }
}

async function handleAvailableMovementMessage(availableMovement) {
    availableMovement.spaces.forEach(space => {
        let spaceElement = document.getElementById(`hex-${space}`)
        if (spaceElement) {
            spaceElement.classList.add("potential-move")

            //The gold stroke won't show up properly on the right and bottom edges because those spaces are later in the HTML
            // and therefore rendered later than this one/on top of this one. To get around this, we'll remove the potential move spaces
            // and then re-add them to the svg, which will put them at the end of the markup i.e. on top
            let svgParent = document.getElementById("polycontainer")
            let spaceText = spaceElement.nextSibling

            svgParent.removeChild(spaceElement)
            svgParent.appendChild(spaceElement)
            svgParent.removeChild(spaceText)
            svgParent.appendChild(spaceText)
        }
    })
}