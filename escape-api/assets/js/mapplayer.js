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

const PlayerTeams = {
    Human: 'Human',
    Alien: 'Alien',
    Spectator: 'Spectator'
}



var currentTool = 'None';

var radius = 20;
var cssClass = 'hexfield';//If you change this, change it in hexClick() too

var MAP = null

function createGrid(rows, columns) {

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
            poly.setAttribute('hex-row', row);
            poly.setAttribute('hex-column', column);
            poly.setAttribute('hex-type', SpaceTypes.Safe);
            poly.setAttribute('id', `hex-${row}-${column}`)
            poly.innerHTML = `<title>[${row},${column}]</title>`; //TODO: This won't work for mobile
            svgParent.appendChild(poly);

            var polyText = document.createElementNS("http://www.w3.org/2000/svg", "text")
            polyText.setAttribute('x', `${center.x}`)
            polyText.setAttribute('y', `${center.y}`)
            polyText.setAttribute('fill', '#303030')
            polyText.setAttribute('text-anchor', 'middle')
            polyText.setAttribute('font-size', 'small')
            polyText.innerHTML = `[${row},${column}]`
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
    if(thisPlayer.team == PlayerTeams.Spectator){
        return
    }

    var row = parseInt(event.target.getAttribute('hex-row') ?? -99);
    var col = parseInt(event.target.getAttribute('hex-column') ?? -99);

    var actionToSend = {
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
}

function clearGrid() {
    var polycontainer = document.getElementById("polycontainer")
    polycontainer?.remove();
}