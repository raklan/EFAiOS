//Heavily inspired by https://github.com/gojko/hexgridwidget, but altered to not require JQuery

const GAME_CONFIG_DEFAULT = {
    workingPods: 4,
    brokenPods: 1,

    numTurns: 40,
    aliensRespawn: false,

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
    teleport: 1,

    numVanillaRolePossible: 1,
    numVanillaRoleRequired: 0,
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
    row: -99,
    col: '!'
}

var selectedSpace2 = {
    row: -99,
    col: '!'
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

let gameHasEnded = false;
let playerHasMoved = false;

const playerNameExtractor = new RegExp(/Player \'(?<PlayerName>[^\']+)\'/g);

var clickMode = ClickModes.None;

var cssClass = 'hexfield';//If you change this, change it in hexClick() too

function drawMapOnPage() {
    if (!MAP) {
        return;
    }

    Object.values(MAP.spaces).forEach(space => {
        var el = document.getElementById(`hex-${space.col}-${space.row}`)
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
        createGrid(MAP.rows, MAP.cols, 35);
        drawMapOnPage();
    } else {
        console.error("No map given")
    }
}

function hexClick(event) {
    if (gameHasEnded) {
        showNotification('The Game has ended', 'Error')
        return
    }

    if (thisPlayer.team == PlayerTeams.Spectator || !isThisPlayersTurn) {
        showNotification('It\'s not your turn!', 'Error')
        return
    }

    var row = parseInt(event.target.getAttribute('hex-row') ?? -99);
    var col = event.target.getAttribute('hex-column') ?? "!";

    var actionToSend = {}
    if (clickMode == ClickModes.Moving) {
        playerHasMoved = true;
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
        var deselectedSpace = document.getElementById(`hex-${selectedSpace2.col}-${selectedSpace2.row}`)
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

        if (selectedSpace2.row != -99 && selectedSpace2.col != '!') {
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
        document.getElementById("playerChoice-targeted").style.display = '';
        typeWord(title, 'Using Sensor')

        popup.style.color = 'white'
        popup.style.border = '2px solid white'

        let content_info = document.getElementById("playerChoice-targeted-content")
        content_info.innerHTML = ''

        typeWord(content_info, 'Choose a player to reveal with the Sensor')
    } else if (mode == 'Scanner') {
        document.getElementById("playerChoice-targeted").style.display = '';
        typeWord(title, 'Using Scanner')

        popup.style.color = 'white'
        popup.style.border = '2px solid white'

        let content_info = document.getElementById("playerChoice-targeted-content")
        content_info.innerHTML = ''

        typeWord(content_info, 'Choose a player to reveal their Role/Team with the Scanner')
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
                row2: -99,
                col2: "!",
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
                row2: -99,
                col2: "!"
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

function scannerConfirm(playerId) {
    hidePlayerChoicePopup();
    clickMode = ClickModes.Moving;

    const actionToSend = {
        gameId: thisGameStateId,
        action: {
            type: 'PlayCard',
            turn: {
                name: 'Scanner',
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
                row: isAttacking ? thisPlayer.row : -99,
                col: isAttacking ? thisPlayer.col : "!"
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

            let tooltip = document.createElement("div")
            tooltip.classList.add("tooltip")
            tooltip.innerText = card.description;
            node.appendChild(tooltip)

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
        var playerList = document.getElementById("playerChoice-targeted-playerList")
        playerList.replaceChildren() //Important: Clear the player list so new players joining don't cause duplicate rendering

        for (let player of gamePlayerList) {
            playerEntry = document.createElement("button")
            playerEntry.innerText = player.name
            playerEntry.style.setProperty('--button-color', 'red')
            playerEntry.onclick = () => {
                sensorConfirm(player.id)
            }

            playerList.appendChild(playerEntry)
        }

        showPlayerChoicePopup(card.name)

        return;
    }
    else if (card.name === 'Scanner') {
        var playerList = document.getElementById("playerChoice-targeted-playerList")
        playerList.replaceChildren() //Important: Clear the player list so new players joining don't cause duplicate rendering

        for (let player of gamePlayerList) {
            playerEntry = document.createElement("button")
            playerEntry.innerText = player.name
            playerEntry.style.setProperty('--button-color', 'red')
            playerEntry.onclick = () => {
                scannerConfirm(player.id)
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

function renderTeamCard() {
    var teamCard = document.getElementById("team")
    teamCard.innerHTML = `<span>${thisPlayer.team}</span>`;
    teamCard.style.setProperty('--team-color', getTeamColor())
}

function renderRoleCard() {
    var roleCard = document.getElementById("role")
    roleCard.innerHTML = `<span>${thisPlayer.role}</span>`;
    roleCard.style.setProperty('--team-color', getTeamColor())
}

function renderStatusEffects() {
    var statusEffectList = document.getElementById("status-effects")
    statusEffectList.replaceChildren()
    statusEffectList.style.setProperty('--team-color', "white")

    let title = document.createElement("h5")
    title.innerText = thisPlayer?.statusEffects?.length > 0 ? "Current Status Effects" : "No Status Effects"
    statusEffectList.appendChild(title)

    for (let statusEffect of thisPlayer.statusEffects) {
        let entry = document.createElement("span")
        entry.innerText = `${statusEffect.name} (${statusEffect.usesLeft})`
        entry.classList.add('status-effect-entry')

        let tooltip = document.createElement("div")
        tooltip.classList.add("tooltip")
        tooltip.innerText = statusEffect.description;
        entry.appendChild(tooltip);

        statusEffectList.appendChild(entry)
    }
}

function renderTurnOrder() {
    var turnOrderList = document.getElementById('turn-order')
    turnOrderList.replaceChildren();
    turnOrderList.style.setProperty('--team-color', "white")

    let title = document.createElement("h5")
    title.innerText = "Turn Order"
    turnOrderList.appendChild(title)

    for (let player of gamePlayerList) {
        let entry = document.createElement("span")
        entry.innerText = `${player.name}`
        if (player.name === thisPlayer.name) {
            entry.innerText += ' (You)'
        }
        if (player.isThisPlayersTurn) {
            entry.classList.add("current-player-turn")
        }

        turnOrderList.appendChild(entry);
    }
}

function renderTurnNumber(turnNum, maxTurns, mapName){
    var turnNumContainer = document.getElementById("turn-number")
    turnNumContainer.innerHTML = `<h4 style="margin-top: 5px">${mapName}</h4><div>Turn ${turnNum} / ${maxTurns}<div>`
}

function getTeamColor() {
    switch (thisPlayer.team) {
        case PlayerTeams.Human:
            return "deepskyblue"
        case PlayerTeams.Alien:
            return "red"
        case PlayerTeams.Spectator:
            return "white"
        default:
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

function setAllConfigAsDefault() {
    setGeneralConfigAsDefault();
    setCardConfigAsDefault();
    setRoleConfigAsDefault();
}

function setGeneralConfigAsDefault() {
    setGeneralConfig(GAME_CONFIG_DEFAULT)
}

function setCardConfigAsDefault() {
    setCardConfig(GAME_CONFIG_DEFAULT)
}

function setRoleConfigAsDefault() {
    setRoleConfig(GAME_CONFIG_DEFAULT)
}

function setConfigForm(configObject) {
    setGeneralConfig(configObject);
    setCardConfig(configObject);
    setRoleConfig(configObject);
}

function setGeneralConfig(configObject) {
    let configForm = document.getElementById("lobby-gameConfig")

    configForm['config-numHumans'].value = 0;
    configForm['config-numAliens'].value = 0;

    configForm['config-numWorkingPods'].value = configObject.workingPods;
    configForm['config-numBrokenPods'].value = configObject.brokenPods;

    configForm['config-numTurns'].value = configObject.numTurns;
    configForm['config-aliensRespawn'].checked = configObject.aliensRespawn
}

function setCardConfig(configObject) {
    let configForm = document.getElementById("lobby-gameConfig")

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

function setRoleConfig(configObject) {
    let configForm = document.getElementById("lobby-gameConfig")
    configForm['config-numCaptain'].value = configObject.numVanillaRolePossible;
    configForm['config-numPilot'].value = configObject.numVanillaRolePossible;
    configForm['config-numCopilot'].value = configObject.numVanillaRolePossible;
    configForm['config-numSoldier'].value = configObject.numVanillaRolePossible;
    configForm['config-numEngineer'].value = configObject.numVanillaRolePossible;
    configForm['config-numPsychologist'].value = configObject.numVanillaRolePossible;
    configForm['config-numEO'].value = configObject.numVanillaRolePossible;
    configForm['config-numMedic'].value = configObject.numVanillaRolePossible;

    configForm['config-numFast'].value = configObject.numVanillaRolePossible;
    configForm['config-numSurge'].value = configObject.numVanillaRolePossible;
    configForm['config-numBlink'].value = configObject.numVanillaRolePossible;
    configForm['config-numSilent'].value = configObject.numVanillaRolePossible;
    configForm['config-numBrute'].value = configObject.numVanillaRolePossible;
    configForm['config-numInvisible'].value = configObject.numVanillaRolePossible;
    configForm['config-numLurking'].value = configObject.numVanillaRolePossible;
    configForm['config-numPsychic'].value = configObject.numVanillaRolePossible;
}

function getGameConfig() {
    let configForm = document.getElementById("lobby-gameConfig")
    let config = {};

    config.numHumans = getConfigValue("config-numHumans")
    config.numAliens = getConfigValue("config-numAliens")

    config.numWorkingPods = getConfigValue('config-numWorkingPods')
    config.numBrokenPods = getConfigValue('config-numBrokenPods')

    config.numTurns = getConfigValue('config-numTurns')
    config.aliensRespawn = configForm['config-aliensRespawn']?.checked;

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

    config.activeRoles = {
        "Captain": getConfigValue('config-numCaptain'),
        "Pilot": getConfigValue('config-numPilot'),
        "Copilot": getConfigValue('config-numCopilot'),
        "Soldier": getConfigValue('config-numSoldier'),
        "Psychologist": getConfigValue('config-numPsychologist'),
        "Executive Officer": getConfigValue('config-numEO'),
        "Medic": getConfigValue('config-numMedic'),
        "Engineer": getConfigValue('config-numEngineer'),

        "Fast": getConfigValue('config-numFast'),
        "Surge": getConfigValue('config-numSurge'),
        "Blink": getConfigValue('config-numBlink'),
        "Silent": getConfigValue('config-numSilent'),
        "Brute": getConfigValue('config-numBrute'),
        "Invisible": getConfigValue('config-numInvisible'),
        "Lurking": getConfigValue('config-numLurking'),
        "Psychic": getConfigValue('config-numPsychic'),
    }

    config.requiredRoles = {
        "Captain": getConfigValue('config-numCaptainRequired'),
        "Pilot": getConfigValue('config-numPilotRequired'),
        "Copilot": getConfigValue('config-numCopilotRequired'),
        "Soldier": getConfigValue('config-numSoldierRequired'),
        "Psychologist": getConfigValue('config-numPsychologistRequired'),
        "Executive Officer": getConfigValue('config-numEORequired'),
        "Medic": getConfigValue('config-numMedicRequired'),
        "Engineer": getConfigValue('config-numEngineerRequired'),

        "Fast": getConfigValue('config-numFastRequired'),
        "Surge": getConfigValue('config-numSurgeRequired'),
        "Blink": getConfigValue('config-numBlinkRequired'),
        "Silent": getConfigValue('config-numSilentRequired'),
        "Brute": getConfigValue('config-numBruteRequired'),
        "Invisible": getConfigValue('config-numInvisibleRequired'),
        "Lurking": getConfigValue('config-numLurkingRequired'),
        "Psychic": getConfigValue('config-numPsychicRequired'),
    }

    function getConfigValue(inputKey) {
        return configForm[inputKey]?.value ? parseInt(configForm[inputKey].value) : 0;
    }

    return config;
}

function updatePossible(inputName) {
    let configForm = document.getElementById("lobby-gameConfig")
    let possible = configForm[`config-${inputName}`]
    let required = configForm[`config-${inputName}Required`]

    possible.min = required.value ? parseInt(required.value) : 0
    possible.value = Math.max(possible.value ? parseInt(possible.value) : 0, possible.min ? parseInt(possible.min) : 0)
}

function checkPossible(inputName) {
    let configForm = document.getElementById("lobby-gameConfig")
    let possible = configForm[`config-${inputName}`]
    possible.value = Math.max(possible.value ? parseInt(possible.value) : 0, possible.min ? parseInt(possible.min) : 0)
}

function endTurn() {
    const actionToSend = {
        gameId: thisGameStateId,
        action: {
            type: 'EndTurn',
            turn: {}
        }
    }
    sendWsMessage(ws, 'submitAction', actionToSend)
    document.getElementById("endTurn-button").style.display = 'none'
    playerHasMoved = false;
}