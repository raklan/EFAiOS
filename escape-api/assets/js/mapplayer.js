//Heavily inspired by https://github.com/gojko/hexgridwidget, but altered to not require JQuery

const GAME_CONFIG_DEFAULT = {
    workingPods: 4,
    brokenPods: 1,
    red: 24,
    green: 26,
    silent: 4,
    adrenaline: 3,
    attack: 1,
    cat: 2,
    clone: 1,
    defense: 1,
    mutation: 1,
    sedatives: 1,
    sensor: 1,
    spotlight: 2,
    teleport: 1
}

const SpaceTypes = {
    Wall: 0,
    Safe: 1,
    Dangerous: 2,
    Pod: 3,
    UsedPod: 4,
    HumanStart: 5,
    AlienStart: 6
}

var selectedSpace = {
    row: '!',
    col: -99
}

var selectedSpace2 = {
    row: '!',
    col: -99
}

const PlayerTeams = {
    Human: 'Human',
    Alien: 'Alien',
    Spectator: 'Spectator'
}

const ClickModes = {
    Moving: 'Moving',
    Noise: 'Noise',
    Spotlight: 'Spotlight',
    CatGreen: 'CatGreen',
    CatRed: 'CatRed',
    None: 'None'
}

const playerNameExtractor = new RegExp(/Player \'(?<PlayerName>[^\']+)\'.+/);

var clickMode = ClickModes.None;

var cssClass = 'hexfield';//If you change this, change it in hexClick() too

var MAP = null

function createGrid(rows, columns) {
    let radius = Math.min(rows * columns / (rows + columns) * (window.screen.width / 150), 50)
    var grid = document.getElementById("gameplay-gridParent");

    var createSVG = function (tag) {
        var newElement = document.createElementNS('http://www.w3.org/2000/svg', tag || 'svg');
        if (tag !== 'svg') //Only add to the polygons
            newElement.addEventListener('click', hexClick);
        return newElement;
    };
    var toPoint = function (dx, dy) {
        return Math.round(dx + center.x) + ',' + Math.round(dy + center.y);
    };

    var height = Math.sqrt(3) / 2 * radius;
    svgParent = createSVG('svg');
    svgParent.setAttribute('tabindex', 1);
    svgParent.setAttribute('id', 'polycontainer')
    grid.appendChild(svgParent);
    svgParent.style.width = `${(1.5 * columns + 0.5) * radius}px`;
    svgParent.style.height = `${(2 * rows + 1) * height}px`;

    for (row = 0; row < rows; row++) {
        for (column = 0; column < columns; column++) {
            center = { x: Math.round((1 + 1.5 * column) * radius), y: Math.round(height * (1 + row * 2 + (column % 2))) };
            let poly = createSVG('polygon');
            poly.setAttribute('points', [
                toPoint(-1 * radius / 2, -1 * height),
                toPoint(radius / 2, -1 * height),
                toPoint(radius, 0),
                toPoint(radius / 2, height),
                toPoint(-1 * radius / 2, height),
                toPoint(-1 * radius, 0)
            ].join(' '));
            poly.setAttribute('class', [cssClass, 'safe'].join(' '));
            poly.setAttribute('tabindex', 1);
            poly.setAttribute('hex-row', numberToLetter(row));
            poly.setAttribute('hex-column', column);
            poly.setAttribute('hex-type', SpaceTypes.Safe);
            poly.setAttribute('id', `hex-${numberToLetter(row)}-${column}`)
            poly.innerHTML = `<title>[${numberToLetter(row)}-${column}]</title>`; //TODO: This won't work for mobile
            svgParent.appendChild(poly);

            var polyText = document.createElementNS("http://www.w3.org/2000/svg", "text")
            polyText.setAttribute('x', `${center.x}`)
            polyText.setAttribute('y', `${center.y}`)
            polyText.setAttribute('fill', 'black')
            polyText.setAttribute('text-anchor', 'middle')
            polyText.setAttribute('font-size', 'small')
            polyText.innerHTML = `[${numberToLetter(row)}-${column}]`
            polyText.style.pointerEvents = 'none'
            svgParent.appendChild(polyText)
        }
    }
}

function drawMapOnPage() {
    if (!MAP) {
        return;
    }

    Object.values(MAP.spaces).forEach(space => {
        var el = document.getElementById(`hex-${space.row}-${space.col}`)
        if (el) {
            var spaceClass = 'safe'
            switch (space.type) {
                case SpaceTypes.Wall: spaceClass = 'wall'; break;
                case SpaceTypes.Safe: spaceClass = 'safe'; break;
                case SpaceTypes.Pod: spaceClass = 'pod'; break;
                case SpaceTypes.UsedPod: spaceClass = 'pod-used'; break;
                case SpaceTypes.Dangerous: spaceClass = 'dangerous'; break;
                case SpaceTypes.HumanStart: spaceClass = 'humanstart'; break;
                case SpaceTypes.AlienStart: spaceClass = 'alienstart'; break;
            }

            el.classList = [cssClass, spaceClass].join(' ');
            el.setAttribute('hex-type', space.type);
        }
    });
}

function drawMap(map) {
    if (map) {
        MAP = map;
        clearGrid();
        createGrid(MAP.rows, MAP.cols);
        drawMapOnPage();
    } else {
        console.error("No map given")
    }
}

function hexClick(event) {
    if (thisPlayer.team == PlayerTeams.Spectator || !isThisPlayersTurn) {
        showNotification('It\'s not your turn!', 'Error')
        return
    }

    var row = event.target.getAttribute('hex-row') ?? "!";
    var col = parseInt(event.target.getAttribute('hex-column') ?? -99);

    var actionToSend = {}
    if (clickMode == ClickModes.Moving) {
        actionToSend = {
            gameId: thisGameStateId,
            action: {
                type: 'Movement',
                turn: {
                    toRow: row,
                    toCol: col
                }
            }
        }
        sendWsMessage(ws, 'submitAction', actionToSend);
    } else if (clickMode == ClickModes.Noise) {
        selectedSpace = {
            row: row,
            col: col
        }
        document.querySelectorAll('.hexfield.selected').forEach(x => x.classList.remove('selected'))
        event.target.classList.add('selected')

        document.getElementById("greenCard-confirm").style.display = ''
    } else if (clickMode == ClickModes.Spotlight) {
        selectedSpace = {
            row: row,
            col: col
        }
        document.querySelectorAll('.hexfield.selected').forEach(x => x.classList.remove('selected'))
        event.target.classList.add('selected')

        document.getElementById("spotlight-confirm").style.display = ''
    } else if (clickMode == ClickModes.CatRed) {
        selectedSpace = {
            row: thisPlayer.row,
            col: thisPlayer.col
        }
        selectedSpace2 = {
            row: row,
            col: col
        }
        document.querySelectorAll('.hexfield.selected').forEach(x => x.classList.remove('selected'))
        event.target.classList.add('selected')

        document.getElementById("cat-confirm").style.display = ''
    } else if (clickMode == ClickModes.CatGreen) {
        var deselectedSpace = document.getElementById(`hex-${selectedSpace2.row}-${selectedSpace2.col}`)
        if (deselectedSpace) {
            deselectedSpace.classList.remove("selected")
        }
        selectedSpace2 = {
            row: selectedSpace.row,
            col: selectedSpace.col
        }
        selectedSpace = {
            row: row,
            col: col
        }
        event.target.classList.add('selected')

        if (selectedSpace2.row != '!' && selectedSpace2.col != -99) {
            document.getElementById("cat-confirm").style.display = ''
        } 
    }



}

function clearGrid() {
    var polycontainer = document.getElementById("polycontainer")
    polycontainer?.remove();
}

function showPlayerChoicePopup(mode) {
    let popup = document.getElementById("playerChoice-popup");
    let title = document.getElementById("playerChoice-title");
    let content = document.getElementById("playerChoice-content");

    title.innerHTML = '';

    for (let child of content.children) {
        child.style.display = 'none'
    }

    if (mode == 'greenCard') {
        document.getElementById("greenCard-confirm").style.display = 'none'
        document.getElementById("playerChoice-greenCard").style.display = '';

        popup.style.color = 'lime'
        popup.style.border = '2px solid lime'

        let content_info = document.getElementById('playerChoice-greenCard-content')
        content_info.innerHTML = ''

        typeWord(title, 'Green Card Drawn')
        typeWord(content_info, 'Choose a space to make noise in')
    } else if (mode == 'redCard') {
        document.getElementById("playerChoice-redCard").style.display = '';
        typeWord(title, 'Red Card Drawn')

        popup.style.color = 'red'
        popup.style.border = '2px solid red'

        let content_info = document.getElementById("playerChoice-redCard-content")
        content_info.innerHTML = ''

        typeWord(content_info, "You're about to make noise in your space")
    } else if (mode == 'attack') {
        document.getElementById("playerChoice-attack").style.display = '';
        typeWord(title, 'Attack Space?')

        popup.style.color = 'white'
        popup.style.border = '2px solid white'

        let content_info = document.getElementById("playerChoice-attack-content")
        content_info.innerHTML = ''

        typeWord(content_info, 'Would you like to attack this space?')
    } else if (mode == 'Spotlight') {
        document.getElementById("playerChoice-spotlight").style.display = '';
        typeWord(title, 'Using Spotlight')

        popup.style.color = 'white'
        popup.style.border = '2px solid white'

        let content_info = document.getElementById("playerChoice-spotlight-content")
        content_info.innerHTML = ''

        typeWord(content_info, 'Choose a space to reveal with the Spotlight')
    } else if (mode == 'Sensor') {
        document.getElementById("playerChoice-sensor").style.display = '';
        typeWord(title, 'Using Sensor')

        popup.style.color = 'white'
        popup.style.border = '2px solid white'

        let content_info = document.getElementById("playerChoice-sensor-content")
        content_info.innerHTML = ''

        typeWord(content_info, 'Choose a player to reveal with the Sensor')
    } else if (mode === 'cat-green') {
        document.getElementById("cat-confirm").style.display = 'none'
        document.getElementById("playerChoice-cat").style.display = '';

        popup.style.color = 'lime'
        popup.style.border = '2px solid lime'

        let content_info = document.getElementById('playerChoice-cat-content')
        content_info.innerHTML = ''

        typeWord(title, 'Green Card Drawn + Cat Activated')
        typeWord(content_info, 'Choose 2 spaces to make noise in')
    } else if (mode === 'cat-red') {
        document.getElementById("cat-confirm").style.display = 'none'
        document.getElementById("playerChoice-cat").style.display = '';

        popup.style.color = 'red'
        popup.style.border = '2px solid red'

        let content_info = document.getElementById('playerChoice-cat-content')
        content_info.innerHTML = ''

        typeWord(title, 'Red Card Drawn + Cat Activated')
        typeWord(content_info, 'Choose an extra space to make noise in')
    }

    popup.classList.add('notification-displayed')
}

function hidePlayerChoicePopup() {
    var popup = document.getElementById("playerChoice-popup");
    popup.classList.remove('notification-displayed')
}

function redCardConfirm() {
    const actionToSend = {
        gameId: thisGameStateId,
        action: {
            type: 'Noise',
            turn: {
                row: thisPlayer.row,
                col: thisPlayer.col,
                row2: "!",
                col2: -99,
            }
        }
    }
    sendWsMessage(ws, 'submitAction', actionToSend)
    hidePlayerChoicePopup();
}

function greenCardConfirm() {
    document.querySelectorAll('.hexfield.selected').forEach(x => x.classList.remove('selected'))
    clickMode = ClickModes.Moving
    const actionToSend = {
        gameId: thisGameStateId,
        action: {
            type: 'Noise',
            turn: {
                row: selectedSpace.row,
                col: selectedSpace.col,
                row2: '!',
                col2: -99
            }
        }
    }
    sendWsMessage(ws, 'submitAction', actionToSend)
    hidePlayerChoicePopup();
    selectedSpace = {
        row: thisPlayer.row,
        col: thisPlayer.col
    }
}

function catConfirm() {
    document.querySelectorAll('.hexfield.selected').forEach(x => x.classList.remove('selected'))
    clickMode = ClickModes.Moving
    const actionToSend = {
        gameId: thisGameStateId,
        action: {
            type: 'Noise',
            turn: {
                row: selectedSpace.row,
                col: selectedSpace.col,
                row2: selectedSpace2.row,
                col2: selectedSpace2.col
            }
        }
    }
    sendWsMessage(ws, 'submitAction', actionToSend)
    hidePlayerChoicePopup();
    selectedSpace = {
        row: thisPlayer.row,
        col: thisPlayer.col
    }
}

function spotlightConfirm() {
    document.querySelectorAll('.hexfield.selected').forEach(x => x.classList.remove('selected'))
    clickMode = ClickModes.Moving
    const actionToSend = {
        gameId: thisGameStateId,
        action: {
            type: 'PlayCard',
            turn: {
                name: 'Spotlight',
                row: selectedSpace.row,
                col: selectedSpace.col
            }
        }
    }
    sendWsMessage(ws, 'submitAction', actionToSend)
    hidePlayerChoicePopup();
    selectedSpace = {
        row: thisPlayer.row,
        col: thisPlayer.col
    }
}

function sensorConfirm(playerId) {
    hidePlayerChoicePopup();
    clickMode = ClickModes.Moving;

    const actionToSend = {
        gameId: thisGameStateId,
        action: {
            type: 'PlayCard',
            turn: {
                name: 'Sensor',
                targetPlayer: playerId
            }
        }
    }

    sendWsMessage(ws, 'submitAction', actionToSend)
}

function attack(isAttacking) {
    var actionToSend = {
        gameId: thisGameStateId,
        action: {
            type: 'Attack',
            turn: {
                row: isAttacking ? thisPlayer.row : "!",
                col: isAttacking ? thisPlayer.col : -99
            }
        }
    }

    sendWsMessage(ws, 'submitAction', actionToSend);
    hidePlayerChoicePopup();
}

function renderPlayerHand() {
    let hand = document.getElementById("cards")
    hand.replaceChildren()

    if (thisPlayer?.hand?.length > 0) {
        for (let card of thisPlayer?.hand) {
            let node = document.createElement("div")
            node.classList = 'card'
            node.innerHTML = `${card.name}`
            node.onclick = () => cardClick(card)

            hand.appendChild(node)
        }
    }
}

function cardClick(card) {
    if (thisPlayer.team == PlayerTeams.Spectator || !isThisPlayersTurn) {
        showNotification('It\'s not your turn!', 'Error')
        return
    }

    if (card.name === "Spotlight") {
        clickMode = ClickModes.Spotlight;
        showPlayerChoicePopup(card.name)
        return;
    }
    else if (card.name === 'Sensor') {
        var playerList = document.getElementById("playerChoice-sensor-playerList")
        playerList.replaceChildren() //Important: Clear the player list so new players joining don't cause duplicate rendering

        for (let player of gamePlayerList) {
            playerEntry = document.createElement("button")
            playerEntry.innerText = player.name
            playerEntry.classList = ['redCard-confirm']
            playerEntry.style.color = 'red'
            playerEntry.onclick = () => {
                sensorConfirm(player.id)
            }

            playerList.appendChild(playerEntry)
        }

        showPlayerChoicePopup(card.name)

        return;
    }
    let toSend = {
        gameId: thisGameStateId,
        action: {
            type: 'PlayCard',
            turn: {
                name: card.name
            }
        }
    }
    sendWsMessage(ws, 'submitAction', toSend)
}

function renderRoleCard() {
    var roleCard = document.getElementById("role")
    roleCard.innerHTML = `${thisPlayer.team}`
    roleCard.style.setProperty('--team-color', getTeamColor())
}

function getTeamColor() {
    switch (thisPlayer.team) {
        case PlayerTeams.Human:
            return "deepskyblue"
        case PlayerTeams.Alien:
            return "red"
        case PlayerTeams.Spectator:
            return "white"
    }
}

function initializeEventLog(players) {
    const tablist = document.getElementById('tab-list')
    const eventLog = document.getElementById("event-log")

    eventLog.onmouseleave = () => {
        // Get all elements with class="tabcontent" and hide them
        tabcontent = document.getElementsByClassName("tabcontent");
        for (i = 0; i < tabcontent.length; i++) {
            tabcontent[i].style.display = "none";
        }
    }

    for (let player of players) {
        let button = document.createElement("button")
        button.classList.add("tablinks")
        button.onclick = () => viewPlayerEvents(player.name)
        button.innerHTML = `${player.name}`
        tablist.appendChild(button)

        let log = document.createElement("div")
        log.id = `event-log-${player.name}`
        log.classList.add("tabcontent")
        eventLog.appendChild(log)
    }
}

function viewPlayerEvents(playerName) {
    // Declare all variables
    var i, tabcontent, tablinks;

    // Get all elements with class="tabcontent" and hide them
    tabcontent = document.getElementsByClassName("tabcontent");
    for (i = 0; i < tabcontent.length; i++) {
        tabcontent[i].style.display = "none";
    }

    // Get all elements with class="tablinks" and remove the class "active"
    tablinks = document.getElementsByClassName("tablinks");
    for (i = 0; i < tablinks.length; i++) {
        tablinks[i].className = tablinks[i].className.replace(" active", "");
    }

    // Show the current tab, and add an "active" class to the button that opened the tab
    document.getElementById(`event-log-${playerName}`).style.display = "block";
}

function addEvent(playerName, event) {
    const eventLogContainer = document.getElementById(`event-log-${playerName}`)
    let eventDesc = document.createElement("p")
    eventDesc.innerHTML = event
    eventLogContainer.appendChild(eventDesc)
}

function setConfigForm(configObject) {
    let configForm = document.getElementById("lobby-gameConfig")

    configForm['config-numHumans'].value = 0;
    configForm['config-numAliens'].value = 0;

    configForm['config-numWorkingPods'].value = configObject.workingPods;
    configForm['config-numBrokenPods'].value = configObject.brokenPods;

    configForm['config-numRedCards'].value = configObject.red;
    configForm['config-numGreenCards'].value = configObject.green;
    configForm['config-numWhiteCards'].value = configObject.silent;

    configForm['config-numTeleport'].value = configObject.teleport;
    configForm['config-numClone'].value = configObject.clone;
    configForm['config-numDefense'].value = configObject.defense;

    configForm['config-numSpotlight'].value = configObject.spotlight;
    configForm['config-numAttack'].value = configObject.attack;
    configForm['config-numSensor'].value = configObject.sensor;

    configForm['config-numAdrenaline'].value = configObject.adrenaline;
    configForm['config-numSedatives'].value = configObject.sedatives;
    configForm['config-numCat'].value = configObject.cat;
    configForm['config-numMutation'].value = configObject.mutation;
}

function getGameConfig() {
    let configForm = document.getElementById("lobby-gameConfig")
    let config = {};

    config.numHumans = getConfigValue("config-numHumans")
    config.numAliens = getConfigValue("config-numAliens")

    config.numWorkingPods = getConfigValue('config-numWorkingPods')
    config.numBrokenPods = getConfigValue('config-numBrokenPods')

    config.activeCards = {
        "Red Card": getConfigValue('config-numRedCards'),
        "Green Card": getConfigValue('config-numGreenCards'),
        "White Card": getConfigValue('config-numWhiteCards'),
        "Teleport": getConfigValue('config-numTeleport'),
        "Clone": getConfigValue('config-numClone'),
        "Defense": getConfigValue('config-numDefense'),
        "Spotlight": getConfigValue('config-numSpotlight'),
        "Attack": getConfigValue('config-numAttack'),
        "Sensor": getConfigValue('config-numSensor'),
        "Adrenaline": getConfigValue('config-numAdrenaline'),
        "Sedatives": getConfigValue('config-numSedatives'),
        "Cat": getConfigValue('config-numCat'),
        "Mutation": getConfigValue('config-numMutation'),
    }

    config.activeStatusEffects = {
        "Armored": 2,
        "Cloned": 1
    }

    function getConfigValue(inputKey) {
        return configForm[inputKey]?.value?.length > 0 ? parseInt(configForm[inputKey].value) : 0;
    }

    return config;
}