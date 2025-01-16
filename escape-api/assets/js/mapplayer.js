//Heavily inspired by https://github.com/gojko/hexgridwidget, but altered to not require JQuery

const SpaceTypes = {
    Wall: 0,
    Safe: 1,
    Dangerous: 2,
    Pod: 3,
    HumanStart: 4,
    AlienStart: 5
}

const WALL_TOOL = 'Walls';
const POD_TOOL = 'Pods';
const SAFE_TOOL = 'Safe Sector';
const DANGER_TOOL = 'Dangerous Sector';
const HUMAN_TOOL = 'Human Start';
const ALIEN_TOOL = 'Alien Start';

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
            svgParent.appendChild(poly);
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
                case SpaceTypes.Dangerous: spaceClass = 'dangerous'; break;
                case SpaceTypes.HumanStart: spaceClass = 'humanstart'; break;
                case SpaceTypes.AlienStart: spaceClass = 'alienstart'; break;
            }

            el.classList = [cssClass, spaceClass].join(' ');
            el.setAttribute('hex-type', space.type);
        }
    });
}

function hexClick(event) {
    event.target.classList = [cssClass]
    event.target.removeAttribute('hex-type')
    switch(currentTool){
        case WALL_TOOL:
            event.target.classList.add('wall');
            event.target.setAttribute('hex-type', SpaceTypes.Wall);
            break;
        case POD_TOOL:
            event.target.classList.add('pod');
            event.target.setAttribute('hex-type', SpaceTypes.Pod);
            break;
        case SAFE_TOOL:
            event.target.classList.add('safe');
            event.target.setAttribute('hex-type', SpaceTypes.Safe);
            break;
        case DANGER_TOOL:
            event.target.classList.add('dangerous');
            event.target.setAttribute('hex-type', SpaceTypes.Dangerous);
            break;
        case ALIEN_TOOL:
            event.target.classList.add('alienstart')
            event.target.setAttribute('hex-type', SpaceTypes.AlienStart);
            break;
        case HUMAN_TOOL:
            event.target.classList.add('humanstart')
            event.target.setAttribute('hex-type', SpaceTypes.HumanStart);
            break;
    }
    event.target.classList.add('hexfield')
}

function clearGrid() {
    var polycontainer = document.getElementById("polycontainer")
    polycontainer?.remove();
}

async function initializeMap(mapId){
    if (mapId){
        var map = await loadMap(mapId);
        MAP = map;
        clearGrid();
        createGrid(map.rows, map.cols);
        drawMapOnPage();
    }else{
        console.error("No map id given")
    }
}

async function loadMap(id){
    var map = null;
    await fetch(`/api/map?id=${id}`)
    .then(resp => resp.json())
    .then(apiObj => {
        map = apiObj;
    })
    return map;
}