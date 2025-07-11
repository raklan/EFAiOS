const WALL_TOOL = 'Walls';
const POD_TOOL = 'Pods';
const SAFE_TOOL = 'Safe Sector';
const DANGER_TOOL = 'Dangerous Sector';
const HUMAN_TOOL = 'Human Start';
const ALIEN_TOOL = 'Alien Start';

var currentTool = DANGER_TOOL;

function drawMapOnPage() {
    if (!MAP) {
        return;
    }

    document.getElementById('name').value = MAP.name;
    document.getElementById('columns').value = MAP.cols;
    document.getElementById('rows').value = MAP.rows;
    document.getElementById('description').value = MAP.description;
    Object.values(MAP.spaces).forEach(space => {
        var el = document.getElementById(`hex-${space.col}-${space.row}`)
        if (el) {
            var spaceClass = 'safe'
            switch (space.type) {
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
            if (spaceClass == 'wall') {
                el.nextElementSibling.setAttribute('fill', 'white')
            }
        }
    });
}

function hexClick(event) {
    event.target.classList = [cssClass]
    event.target.nextElementSibling.setAttribute('fill', 'black')
    switch (currentTool) {
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
        default:
            break;
    }
}

function rebuildGrid() {
    var
        columns = parseInt(document.getElementById('columns').value),
        rows = parseInt(document.getElementById('rows').value);
    clearGrid();
    createGrid(rows, columns, document.getElementById("gridParent"));
};

function setTool(newTool) {
    currentTool = newTool;
    document.getElementById("current-tool").innerText = `Current Tool: ${newTool}`;
}

async function initializePage() {
    const urlParams = new URLSearchParams(window.location.search);
    if (urlParams.has("id")) {
        var map = await loadMap(urlParams.get("id"));
        MAP = map;
        createGrid(map.rows, map.cols, document.getElementById("gridParent"));
        drawMapOnPage();
        setConfigFormFromObject(map.gameConfig);
    } else {
        rebuildGrid();
    }
}

async function exportMap() {
    var map = {
        id: MAP?.id ?? '',
        spaces: {},
        name: document.getElementById('name')?.value ?? 'No Name',
        cols: parseInt(document.getElementById('columns')?.value ?? 0),
        rows: parseInt(document.getElementById('rows')?.value ?? 0),
        description: document.getElementById('description')?.value ?? "No Description Given",
        gameConfig: getGameConfig()
    };
    var polycontainer = document.getElementById("polycontainer")

    for (child of polycontainer.children) {
        var row = child.getAttribute('hex-row')
        col = child.getAttribute('hex-column')
        type = child.getAttribute('hex-type')
        if(row && col){
            map.spaces[`${col}-${row}`] = {
                row: parseInt(row),
                col: col,
                type: parseInt(type)
            }
        }
    }

    if (!validateMap(map)) {
        return;
    }

    fetch('/api/map', {
        method: "POST",
        body: JSON.stringify(map)
    }).then(resp => resp.json()).then(apiObj => showNotification(`Map \"${apiObj.name}\" saved`, 'Success'))
}

function validateMap(map) {
    let keys = Object.keys(map.spaces)
    if(!keys.some(key => map.spaces[key].type == SpaceTypes.Pod)){
        showNotification("Map must have at least 1 Escape Pod", "Error")
        return false;
    }
    if(!keys.some(key => map.spaces[key].type == SpaceTypes.HumanStart)){
        showNotification("Map must have at least 1 Human Start", "Error")
        return false;
    }
    if(!keys.some(key => map.spaces[key].type == SpaceTypes.AlienStart)){
        showNotification("Map must have at least 1 Alien Start", "Error")
        return false
    }
    if(!map.name?.length > 0){
        showNotification("Map must be given a Name", "Error")
        return false;
    }
    return true
}

async function loadMap(id) {
    var map = null;
    await fetch(`/api/map?id=${id}`)
        .then(resp => resp.json())
        .then(apiObj => {
            map = apiObj;
        })
    return map;
}