var canvas, ctx, flag = false,
    prevX = 0,
    currX = 0,
    prevY = 0,
    currY = 0,
    dot_flag = false;

let isPainting = false;

var selectedColor = "green",
    strokeWidth = 2;

function initializeCanvas() {
    console.log('initializing canvas')
    canvas = document.getElementById('can');
    ctx = canvas.getContext("2d");
    resizeCanvasToDisplaySize(canvas, document.getElementById("gridParent"));
    canvas.style.pointerEvents = 'none';
    w = canvas.width;
    h = canvas.height;

    canvas.addEventListener("mousemove", function (e) {
        findxy('move', e)
    }, false);
    canvas.addEventListener("mousedown", function (e) {
        findxy('down', e)
    }, false);
    canvas.addEventListener("mouseup", function (e) {
        findxy('up', e)
    }, false);
    canvas.addEventListener("mouseout", function (e) {
        findxy('out', e)
    }, false);

    canvas.addEventListener("touchmove", function (e) {
        findxy('move', e)
    }, false);
    canvas.addEventListener("touchstart", function (e) {
        findxy('down', e)
    }, false);
    canvas.addEventListener("touchend", function (e) {
        findxy('up', e)
    }, false);
    canvas.addEventListener("touchcancel", function (e) {
        findxy('up', e)
    }, false);
}

function color(obj) {
    selectedColor = obj.getAttribute('color-value')
    document.querySelectorAll('.color-picker.selected').forEach(el => el.classList.remove('selected'))
    obj.classList.add('selected')
}

function updateStrokeWidth(control) {
    strokeWidth = control.value
}

function colorPicker(control) {
    control.setAttribute('color-value', control.value)
    control.style.setProperty('--color-input-bg', control.value)
    color(control)
}

function draw() {
    ctx.beginPath();
    ctx.strokeStyle = selectedColor;
    ctx.lineWidth = strokeWidth;
    if (selectedColor == "eraser") {
        ctx.globalCompositeOperation = "destination-out";
    } else {
        ctx.globalCompositeOperation = "source-over";
    }
    ctx.moveTo(prevX, prevY);
    ctx.lineTo(currX, currY);
    ctx.stroke();
    ctx.arc(prevX, prevY, strokeWidth / 2, 0, Math.PI * 2, false);
    ctx.fill();
    ctx.closePath();
}

function eraseCanvas() {
    ctx.clearRect(0, 0, w, h);
}

function findxy(res, e) {
    if (res == 'down') {

        let clientX, clientY;

        clientX = e.clientX ? e.clientX : e.targetTouches[0].clientX;
        clientY = e.clientY ? e.clientY : e.targetTouches[0].clientY;

        prevX = currX;
        prevY = currY;
        currX = clientX - canvas.offsetLeft;
        currY = clientY - canvas.offsetTop;

        flag = true;
        dot_flag = true;
        if (dot_flag) {
            ctx.beginPath();
            ctx.fillStyle = selectedColor;
            ctx.fillRect(currX, currY, 2, 2);
            ctx.closePath();
            dot_flag = false;
        }
    }
    if (res == 'up' || res == "out") {
        flag = false;
    }
    if (res == 'move') {
        if (flag) {
            let clientX, clientY;

            clientX = e.clientX ? e.clientX : e.targetTouches[0].clientX;
            clientY = e.clientY ? e.clientY : e.targetTouches[0].clientY;

            prevX = currX;
            prevY = currY;
            currX = clientX - canvas.offsetLeft;
            currY = clientY - canvas.offsetTop;
            draw();
        }
    }
}

function togglePaint() {
    isPainting = !isPainting;
    let paintControls = document.getElementById('paint-controls');
    let pencilIcon = document.getElementById('paint-toggle-pencil');
    let closeIcon = document.getElementById('paint-toggle-close');
    let paintToggle = document.getElementById('paint-toggle');

    if (isPainting) {
        document.getElementById("can").style.pointerEvents = '';
        paintToggle.title = 'Close Paint Tray';
        pencilIcon.style.display = 'none';
        closeIcon.style.display = '';
        paintControls.classList.add('displayed');
    } else {
        document.getElementById("can").style.pointerEvents = 'none';
        paintToggle.title = 'Open Paint Tray';
        pencilIcon.style.display = '';
        closeIcon.style.display = 'none';
        paintControls.classList.remove('displayed');
    }
}

function resizeCanvasToDisplaySize(canvas, sizeTo) {
    // look up the size the canvas is being displayed
    const width = sizeTo.clientWidth;
    const height = sizeTo.clientHeight;
    console.log('resizing canvas', canvas, width, height)

    canvas.style.width = `${width}px`;
    canvas.style.height = `${height}px`;

    // If it's resolution does not match change it
    if (canvas.width !== width || canvas.height !== height) {
        canvas.width = width;
        canvas.height = height;
        return true;
    }

    return false;
}