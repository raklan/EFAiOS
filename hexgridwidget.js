function createGrid(radius, columns, rows, cssClass) {

    var grid = document.getElementById("gridParent");

    var createSVG = function (tag) {
        var newElement = document.createElementNS('http://www.w3.org/2000/svg', tag || 'svg');
        if(tag !== 'svg') //Only add to the polygons
            newElement.addEventListener('click', hexClick);
        return $(newElement);
    };
    var toPoint = function (dx, dy) {
        return Math.round(dx + center.x) + ',' + Math.round(dy + center.y);
    };

    var height = Math.sqrt(3) / 2 * radius;
    svgParent = createSVG('svg').attr({ 'tabindex': 1, 'id': 'polycontainer' }).appendTo(grid).css({
        width: (1.5 * columns + 0.5) * radius,
        height: (2 * rows + 1) * height
    });

    for (row = 0; row < rows; row++) {
        for (column = 0; column < columns; column++) {
            center = { x: Math.round((1 + 1.5 * column) * radius), y: Math.round(height * (1 + row * 2 + (column % 2))) };
            createSVG('polygon').attr({
                points: [
                    toPoint(-1 * radius / 2, -1 * height),
                    toPoint(radius / 2, -1 * height),
                    toPoint(radius, 0),
                    toPoint(radius / 2, height),
                    toPoint(-1 * radius / 2, height),
                    toPoint(-1 * radius, 0)
                ].join(' '),
                'class': cssClass,
                tabindex: 1
            })
                .appendTo(svgParent).data({ center: center, row: row, column: column }).attr({ 'hex-row': row, 'hex-column': column });
        }
    }
}

function hexClick(event) {
    if(event.target.classList.contains('clicked')){
        event.target.classList.remove('clicked')
    } else{
        event.target.classList.add('clicked')
    }
}

function clearGrid() {
    var polycontainer = document.getElementById("polycontainer")
    polycontainer?.remove();
}

function rebuildGrid() {
    var
        radius = parseInt($('#radius').val()),
        columns = parseInt($('#columns').val()),
        rows = parseInt($('#rows').val()),
        cssClass = 'hexfield';
    clearGrid();
    createGrid(radius, columns, rows, cssClass);
};