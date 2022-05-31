let sock = null;

const connected = document.getElementById("connected");

const userCount = document.getElementById("userCount");


const wsuri =  `${location.protocol === "https:" ? "wss" : "ws"}://${window.location.host}/ws`;

sock = new WebSocket(wsuri);

sock.onopen = () => {
    console.log("connected to " + wsuri);
    connected.innerText = "🟢";
};

sock.onclose = (e) => {
    console.log("connection closed (" + e.code + ")");
    connected.innerText = "🔴";
};

sock.onmessage = (e) => {
    const drawData = JSON.parse(e.data);
    try {
        if (drawData.type === "canvasData") {
            drawExternal(drawData.data);
        }else if (drawData.type === "userCount"){
            userCount.innerText = drawData.data;
        }
    } catch (error) {
        console.log("Unable to execute draw data")
    }
};

//canvas part
let mouseDown = false;

const canvas = document.getElementById("canvas1");
const colourSelect = document.getElementById("colourSelect");

canvas.addEventListener("mousedown", () => {
    mouseDown = true;
});
canvas.addEventListener("mouseup", () => {
    mouseDown = false;
    [old.x, old.y] = [null, null];
});

colourSelect.addEventListener("change", e => strokeColour = e.target.value);


const strokeColourDefault = "#000000";
let strokeColour = strokeColourDefault;

const ctx = canvas.getContext("2d")
ctx.strokeStyle = strokeColourDefault;


function resize() {
    canvas.width = window.innerWidth;
    canvas.height = window.innerHeight;
}

let old = {
    x: null,
    y: null,
}

function drawExternal(data) {
    ctx.beginPath();
    ctx.strokeStyle = data.colour;
    ctx.lineWidth = 0.7 * window.devicePixelRatio;
    ctx.moveTo(data.old.x, data.old.y);
    ctx.lineTo(data.new.x, data.new.y);
    ctx.stroke();
    ctx.closePath();
}

function draw(event) {
    if (old.x === null && old.y === null) {
        old.x = event.clientX
        old.y = event.clientY
    }
    ctx.beginPath();
    ctx.strokeStyle = strokeColour;
    ctx.lineWidth = 0.7 * window.devicePixelRatio;
    ctx.moveTo(old.x, old.y);
    ctx.lineTo(event.clientX, event.clientY);
    ctx.stroke();
    ctx.closePath();
    const drawingData = {
        type: "canvasData",
        data: {
            old,
            new: {
                x: event.clientX,
                y: event.clientY
            },
            colour: strokeColour
        }
    }
    sock.send(JSON.stringify(drawingData));
    old.x = event.clientX
    old.y = event.clientY
}
resize();

window.addEventListener("resize", () => {
    resize();
});

canvas.addEventListener("mousemove", (event) => {
    if(mouseDown){
        draw(event);
    }
});