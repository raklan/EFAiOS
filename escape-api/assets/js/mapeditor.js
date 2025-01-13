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

function createGrid(radius, columns, rows, cssClass) {

    var grid = document.getElementById("gridParent");

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
            svgParent.appendChild(poly);
        }
    }
}

function hexClick(event) {
    event.target.classList = []
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

function rebuildGrid() {
    var
        radius = parseInt(document.getElementById('radius').value),
        columns = parseInt(document.getElementById('columns').value),
        rows = parseInt(document.getElementById('rows').value),
        cssClass = 'hexfield';//If you change this, change it in hexClick() too
    clearGrid();
    createGrid(radius, columns, rows, cssClass);
};

function setTool(newTool){
    currentTool = newTool;
    document.getElementById("current-tool").innerText = `Current Tool: ${newTool}`;
}

function exportMap(){
    var game = {
        spaces: {}
    };
    var polycontainer = document.getElementById("polycontainer")
    console.log('exporting', polycontainer.children)

    for(child of polycontainer.children){
        var row = child.getAttribute('hex-row')
            col = child.getAttribute('hex-column')
            type = child.getAttribute('hex-type')
        game.spaces[`${row},${col}`] = {
            row: parseInt(row),
            col: parseInt(col),
            type: parseInt(type)
        }
    }

    fetch('/api/map', {
        method: "POST",
        body: JSON.stringify(game)
    }).then(resp => resp.json()).then(apiObj => console.log(apiObj))
}