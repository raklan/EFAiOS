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

const WALL_TOOL = 'Walls';
const POD_TOOL = 'Pods';
const SAFE_TOOL = 'Safe Sector';
const DANGER_TOOL = 'Dangerous Sector';
const HUMAN_TOOL = 'Human Start';
const ALIEN_TOOL = 'Alien Start';

var currentTool = 'None';


var cssClass = 'hexfield';//If you change this, change it in hexClick() too

var MAP = null

function createGrid(rows, columns) {
    const byMap = 700/(columns * 0.5 + rows * 0.5);
    const byWindow = window.innerWidth/50;
    let radius = Math.min(byMap, byWindow)
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
            poly.setAttribute('class', [cssClass, 'dangerous'].join(' '));
            poly.setAttribute('tabindex', 1);
            poly.setAttribute('hex-row', row);
            poly.setAttribute('hex-column', numberToLetter(column));
            poly.setAttribute('hex-type', SpaceTypes.Dangerous);
            poly.setAttribute('id', `hex-${numberToLetter(column)}-${row}`)
            svgParent.appendChild(poly);

            var polyText = document.createElementNS("http://www.w3.org/2000/svg", "text")
            polyText.setAttribute('x', `${center.x}`)
            polyText.setAttribute('y', `${center.y}`)
            polyText.setAttribute('fill', 'black')
            polyText.setAttribute('text-anchor', 'middle')
            polyText.setAttribute('font-size', `${radius/2.25}px`)
            polyText.innerHTML = `[${numberToLetter(column)}-${row}]`
            polyText.style.pointerEvents = 'none'
            svgParent.appendChild(polyText)
        }
    }
}

function drawMapOnPage(){
    if(!MAP){
        return;
    }

    document.getElementById('name').value = MAP.name;
    document.getElementById('columns').value = MAP.cols;
    document.getElementById('rows').value = MAP.rows;
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

            el.nextElementSibling.setAttribute('fill', 'black')
            if(spaceClass == 'wall'){
                el.nextElementSibling.setAttribute('fill', 'white')
            }
        }
    });
}

function hexClick(event) {
    event.target.classList = [cssClass]
    event.target.removeAttribute('hex-type')
    event.target.nextElementSibling.setAttribute('fill', 'black')
    switch(currentTool){
        case WALL_TOOL:
            event.target.classList.add('wall');
            event.target.setAttribute('hex-type', SpaceTypes.Wall);
            event.target.nextElementSibling.setAttribute('fill', 'white')
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
    columns = parseInt(document.getElementById('columns').value),
    rows = parseInt(document.getElementById('rows').value);    
    clearGrid();
    createGrid(rows, columns);
};

function setTool(newTool){
    currentTool = newTool;
    document.getElementById("current-tool").innerText = `Current Tool: ${newTool}`;
}

async function initializePage(){
    const urlParams = new URLSearchParams(window.location.search);
    if (urlParams.has("id")){
        var map = await loadMap(urlParams.get("id"));
        MAP = map;
        createGrid(map.rows, map.cols);
        drawMapOnPage();
    }else{
        rebuildGrid();
    }
}

async function exportMap(){
    var map = {
        id: MAP?.id ?? '',
        spaces: {},
        name: document.getElementById('name')?.value ?? 'No Name',
        cols: parseInt(document.getElementById('columns')?.value ?? 0),
        rows: parseInt(document.getElementById('rows')?.value ?? 0)
    };
    var polycontainer = document.getElementById("polycontainer")

    for(child of polycontainer.children){
        var row = child.getAttribute('hex-row')
            col = child.getAttribute('hex-column')
            type = child.getAttribute('hex-type')
        map.spaces[`${col}-${row}`] = {
            row: parseInt(row),
            col: col,
            type: parseInt(type)
        }
    }

    fetch('/api/map', {
        method: "POST",
        body: JSON.stringify(map)
    }).then(resp => resp.json()).then(apiObj => console.log(apiObj))
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