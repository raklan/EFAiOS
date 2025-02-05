//Heavily inspired by https://github.com/gojko/hexgridwidget, but altered to not require JQuery

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

const PlayerTeams = {
    Human: 'Human',
    Alien: 'Alien',
    Spectator: 'Spectator'
}

const ClickModes = {
    Moving: 'Moving',
    Noise: 'Noise',
    None: 'None'
}

var clickMode = ClickModes.None;

var cssClass = 'hexfield';//If you change this, change it in hexClick() too

var MAP = null

function createGrid(rows, columns) {
    let radius = Math.min(rows * columns / (rows + columns) * (window.screen.width / 150), 55)
    var grid = document.getElementById("gameplay-gridParent");

    var createSVG = function (tag) {
        var newElement = document.createElementNS('http://www.w3.org/2000/svg', tag || 'svg');
        if(tag !== 'svg') //Only add to the polygons
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

function drawMapOnPage(){
    if(!MAP){
        return;
    }

    Object.values(MAP.spaces).forEach(space => {
        var el = document.getElementById(`hex-${space.row}-${space.col}`)
        if(el){
            var spaceClass = 'safe'
            switch (space.type){
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

function drawMap(map){
    if (map){
        MAP = map;
        clearGrid();
        createGrid(MAP.rows, MAP.cols);
        drawMapOnPage();
    }else{
        console.error("No map given")
    }
}

function hexClick(event) {
    if(thisPlayer.team == PlayerTeams.Spectator || !isThisPlayersTurn){
        return
    }

    var row = event.target.getAttribute('hex-row') ?? "!";
    var col = parseInt(event.target.getAttribute('hex-column') ?? -99);

    var actionToSend = {}
    if(clickMode == ClickModes.Moving){
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
    } else if(clickMode == ClickModes.Noise){
        selectedSpace = {
            row: row,
            col: col
        }
        document.querySelectorAll('.hexfield.selected').forEach(x => x.classList.remove('selected'))
        event.target.classList.add('selected')

        document.getElementById("greenCard-confirm").style.display = ''
    }

    
    
}

function clearGrid() {
    var polycontainer = document.getElementById("polycontainer")
    polycontainer?.remove();
}

function showPlayerChoicePopup(mode){
    let popup = document.getElementById("playerChoice-popup");
    let title = document.getElementById("playerChoice-title");
    let content = document.getElementById("playerChoice-content");

    title.innerHTML = '';

    for(let child of content.children){
        child.style.display = 'none'
    }

    if(mode == 'greenCard'){
        document.getElementById("greenCard-confirm").style.display = 'none'
        document.getElementById("playerChoice-greenCard").style.display = '';

        popup.style.color = 'lime'
        popup.style.border = '2px solid lime'

        let content_info = document.getElementById('playerChoice-greenCard-content')
        content_info.innerHTML = ''

        typeWord(title, 'Green Card Drawn')
        typeWord(content_info, 'Choose a space to make noise in')
    }else if(mode == 'redCard'){
        document.getElementById("playerChoice-redCard").style.display = '';
        typeWord(title, 'Red Card Drawn')

        popup.style.color = 'red'
        popup.style.border = '2px solid red'

        let content_info = document.getElementById("playerChoice-redCard-content")
        content_info.innerHTML = ''

        typeWord(content_info, "You're about to make noise in your space")
    }else if(mode == 'attack'){
        document.getElementById("playerChoice-attack").style.display = '';
        typeWord(title, 'Attack Space?')

        popup.style.color = 'white'
        popup.style.border = '2px solid white'

        let content_info = document.getElementById("playerChoice-attack-content")
        content_info.innerHTML = ''

        typeWord(content_info, 'Would you like to attack this space?')
    }

    popup.classList.add('notification-displayed')
}

function hidePlayerChoicePopup(){
    var popup = document.getElementById("playerChoice-popup");
    popup.classList.remove('notification-displayed')
}

function redCardConfirm(){
    const actionToSend = {
        gameId: thisGameStateId,
        action: {
            type: 'Noise',
            turn: {
                row: thisPlayer.row,
                col: thisPlayer.col
            }
        }
    }
    sendWsMessage(ws, 'submitAction', actionToSend)
    hidePlayerChoicePopup();
}

function greenCardConfirm(){
    document.querySelectorAll('.hexfield.selected').forEach(x => x.classList.remove('selected'))
    clickMode = ClickModes.Moving
    const actionToSend = {
        gameId: thisGameStateId,
        action: {
            type: 'Noise',
            turn: {
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

function attack(isAttacking){
    var actionToSend = {
        gameId: thisGameStateId,
        action: {
            type: 'Attack',
            turn: {
                row: isAttacking? thisPlayer.row : "!",
                col: isAttacking? thisPlayer.col : -99
            }
        }
    }

    sendWsMessage(ws, 'submitAction', actionToSend);
    hidePlayerChoicePopup();
}

function renderPlayerHand(){
    let hand = document.getElementById("cards")

    if(thisPlayer?.hand?.cards?.length > 0){
        hand.replaceChildren()
        for(let card of thisPlayer?.hand?.cards){
            let node = document.createElement("div")
            node.classList = 'card'
            node.innerHTML = `${card.name}`
            node.onclick = () => cardClick(card)
    
            hand.appendChild(node)
        }
    }
}

function cardClick(card){
    console.log(card)
}