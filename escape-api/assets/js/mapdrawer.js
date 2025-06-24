var MAP = null

function createGrid(rows, columns, windowWidthDenominator) {
    const byMap = 700/(columns * 0.5 + rows * 0.5);
    const byWindow = window.innerWidth/windowWidthDenominator;
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
            poly.setAttribute('hex-row', row+1);
            poly.setAttribute('hex-column', numberToLetter(column));
            poly.setAttribute('hex-type', SpaceTypes.Dangerous);
            poly.setAttribute('id', `hex-${numberToLetter(column)}-${row+1}`)
            svgParent.appendChild(poly);

            var polyText = document.createElementNS("http://www.w3.org/2000/svg", "text")
            polyText.setAttribute('x', `${center.x}`)
            polyText.setAttribute('y', `${center.y}`)
            polyText.setAttribute('fill', 'black')
            polyText.setAttribute('text-anchor', 'middle')
            polyText.setAttribute('font-size', `${radius/2.25}px`)
            //polyText.setAttribute('textLength', `${radius/2.25}px`)
            polyText.innerHTML = `[${numberToLetter(column)}-${row+1}]`
            polyText.style.pointerEvents = 'none'
            svgParent.appendChild(polyText)
        }
    }
}

function showSpaceTooltip(event){
    let tooltip = document.getElementById("spaceHoverTooltip")
    tooltip.style.display = 'block'
    tooltip.style.left = event.pageX + 10 + 'px'
    tooltip.style.top = event.pageY + 10 + 'px'
    tooltip.innerHTML = event.target.getAttribute('tooltip-text')
    tooltip.style.setProperty('--tooltip-color', event.target.getAttribute('tooltip-color') )
}

function hideSpaceTooltip(event){
    let tooltip = document.getElementById("spaceHoverTooltip")
    tooltip.style.display = 'none'
}