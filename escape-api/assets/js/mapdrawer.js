//Heavily inspired by https://github.com/gojko/hexgridwidget, but altered to not require JQuery and to support my use case better

var MAP = null;
var cssClass = 'hexfield';//If you change this, change it in hexClick() too

const SpaceTypes = {
    Wall: 0,
    Safe: 1,
    Dangerous: 2,
    Pod: 3,
    UsedPod: 4,
    HumanStart: 5,
    AlienStart: 6
};

function createGrid(rows, columns, container, withText = true) {
    let radius = Math.min(container.clientWidth * 0.6 / columns, container.clientHeight * 0.5 / rows)
    var grid = container;
    let tooltip = document.createElement('div')
    tooltip.id = 'spaceHoverTooltip'
    tooltip.style.textWrap = 'nowrap'
    grid.appendChild(tooltip)

    var createSVG = function (tag) {
        var newElement = document.createElementNS('http://www.w3.org/2000/svg', tag || 'svg');
        if(tag !== 'svg') //Only add to the polygons
            newElement.addEventListener('click', hexClick);
            newElement.addEventListener('dblclick', null);
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
            poly.setAttribute('hex-row', row+1);
            poly.setAttribute('hex-column', numberToLetter(column));
            poly.setAttribute('hex-type', SpaceTypes.Dangerous);
            poly.setAttribute('id', `hex-${numberToLetter(column)}-${row+1}`)
            svgParent.appendChild(poly);

            if(withText){
                var polyText = document.createElementNS("http://www.w3.org/2000/svg", "text")
                polyText.setAttribute('x', `${center.x}`)
                polyText.setAttribute('y', `${center.y}`)
                polyText.setAttribute('fill', 'black');
                polyText.setAttribute('text-anchor', 'middle');
                polyText.setAttribute('font-size', `${radius/2.25}px`);
                //polyText.setAttribute('textLength', `${radius/2.25}px`)
                polyText.innerHTML = `[${numberToLetter(column)}-${row+1}]`
                polyText.style.pointerEvents = 'none';
                polyText.style.userSelect = 'none';
                svgParent.appendChild(polyText)
            }
        }
    }
}

function clearGrid() {
    var polycontainer = document.getElementById("polycontainer")
    polycontainer?.remove();
}

function drawMapOnPage() {
    if (!MAP) {
        return;
    }

    Object.values(MAP.spaces).forEach(space => {
        var el = document.getElementById(`hex-${space.col}-${space.row}`)
        if (el) {
            var spaceClass = 'safe'
            var tooltipText = ''
            switch (space.type) {
                case SpaceTypes.Wall:
                    spaceClass = 'wall';
                    tooltipText = '';
                    break;
                case SpaceTypes.Safe:
                    spaceClass = 'safe';
                    tooltipText = ''
                    break;
                case SpaceTypes.Pod:
                    spaceClass = 'pod';
                    tooltipText = 'Escape Pod';
                    break;
                case SpaceTypes.UsedPod:
                    spaceClass = 'pod-used';
                    tooltipText = 'Used Escape Pod';
                    break;
                case SpaceTypes.Dangerous:
                    spaceClass = 'dangerous';
                    tooltipText = '';
                    break;
                case SpaceTypes.HumanStart:
                    spaceClass = 'humanstart';
                    tooltipText = 'Human Start Sector';
                    break;
                case SpaceTypes.AlienStart:
                    spaceClass = 'alienstart';
                    tooltipText = 'Alien Start Sector'
                    break;
            }

            el.classList = [cssClass, spaceClass].join(' ');
            el.setAttribute('hex-type', space.type);
            el.setAttribute('tooltip-text', tooltipText)
            el.setAttribute('tooltip-color', `var(--space-${spaceClass})`)
            if (tooltipText.length > 0) {
                el.onmousemove = (event) => showSpaceTooltip(event)
                el.onmouseleave = (event) => hideSpaceTooltip(event)
            }
        }
    });
}

function showSpaceTooltip(event){
    let tooltip = document.getElementById("spaceHoverTooltip")
    tooltip.style.display = 'block'
    tooltip.style.left = event.layerX + 10 + 'px'
    tooltip.style.top = event.layerY + 10 + 'px'
    tooltip.innerHTML = event.target.getAttribute('tooltip-text')
    tooltip.style.setProperty('--tooltip-color', event.target.getAttribute('tooltip-color') )
}

function hideSpaceTooltip(event){
    let tooltip = document.getElementById("spaceHoverTooltip")
    tooltip.style.display = 'none'
}