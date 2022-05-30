let sock = null;
let recentMessage = "";

let isDrawing = false;

const wsuri = `ws://${window.location.host}/ws`;

let dpr = window.devicePixelRatio;

const c = document.querySelector("canvas");
c.width = window.innerWidth * window.devicePixelRatio;
c.height = window.innerHeight * window.devicePixelRatio;
c.style.width = window.innerWidth + 'px';
c.style.height = window.innerHeight + 'px';

let ctx = c.getContext("2d");
ctx.lineWidth = 0.7 * dpr;
ctx.strokeStyle = "#000000";
ctx.lineJoin = 'round';
ctx.lineCap = 'round';

sock = new WebSocket(wsuri);

sock.onopen = () => {
    console.log("connected to " + wsuri);
};

sock.onclose = (e) => {
    console.log("connection closed (" + e.code + ")");
};

sock.onmessage = (e) => {
    const msg = e.data;
    console.log(`received: ${msg}`)
};

//canvas part

const paintCanvas = document.querySelector( '.js-paint' );
const context = paintCanvas.getContext( '2d' );
context.lineCap = 'round';

const colorPicker = document.querySelector( '.js-color-picker');

colorPicker.addEventListener( 'change', event => {
    context.strokeStyle = event.target.value;
} );

const lineWidthRange = document.querySelector( '.js-line-range' );
const lineWidthLabel = document.querySelector( '.js-range-value' );

lineWidthRange.addEventListener( 'input', event => {
    const width = event.target.value;
    lineWidthLabel.innerHTML = width;
    context.lineWidth = width;
} );

let x = 0, y = 0;
let isMouseDown = false;

const stopDrawing = () => { isMouseDown = false; }
const startDrawing = event => {
    isMouseDown = true;
   [x, y] = [event.offsetX, event.offsetY];
}
const drawLine = event => {
    if ( isMouseDown ) {
        const newX = event.offsetX;
        const newY = event.offsetY;
        context.beginPath();
        context.moveTo( x, y );
        context.lineTo( newX, newY );
        context.stroke();
        //[x, y] = [newX, newY];
        x = newX;
        y = newY;
    }
}

paintCanvas.addEventListener( 'mousedown', startDrawing );
paintCanvas.addEventListener( 'mousemove', drawLine );
paintCanvas.addEventListener( 'mouseup', stopDrawing );
paintCanvas.addEventListener( 'mouseout', stopDrawing );