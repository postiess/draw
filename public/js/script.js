const connected = document.getElementById("connected");
const userCount = document.getElementById("userCount");
const canvas = document.getElementById("canvas1");
const colourSelect = document.getElementById("colourSelect");

const wsuri = `${location.protocol === "https:" ? "wss" : "ws"}://${window.location.host}/ws`;
const sock = new WebSocket(wsuri);

sock.onopen = () => {
    console.log("connected to " + wsuri);
    connected.innerText = "🟢";
};

sock.onclose = (event) => {
    console.log("connection closed (" + event.code + ")");
    connected.innerText = "🔴";
};

sock.onmessage = (event) => {
    try {
        console.log(event.data)
        const drawData = JSON.parse(event.data);
        if (drawData.type === "canvasData") {
            draw(null, drawData.data);
        } else if (drawData.type === "userCount") {
            userCount.innerText = drawData.data;
        }
    } catch (error) {
        console.log(error)
        console.log("Unable to parse draw data");
    }
};

//canvas part

const resize = () => [canvas.width, canvas.height] = [window.innerWidth, window.innerHeight];

const config = {
    mouseDown: false,
    strokeColourDefault: "#000000",
    old: {
        x: null,
        y: null,
    },
    strokeColour: this.strokeColourDefault
}

const ctx = canvas.getContext("2d");
ctx.strokeStyle = config.strokeColourDefault;

const draw = (event, data) => {
    const isExternal = event === null; //not local drawing data

    if (!isExternal && config.old.x === null && config.old.y === null)[config.old.x, config.old.y] = [event.clientX, event.clientY];

    ctx.beginPath();
    ctx.strokeStyle = (isExternal ? data.colour : config.strokeColour);
    ctx.moveTo((isExternal ? data.old.x : config.old.x), (isExternal ? data.old.y : config.old.y));
    ctx.lineTo((isExternal ? data.new.x : event.clientX), (isExternal ? data.new.y : event.clientY));
    ctx.stroke();
    ctx.closePath();

    if (!isExternal) {
        const drawingData = {
            type: "canvasData",
            data: {
                old: config.old,
                new: {
                    x: event.clientX,
                    y: event.clientY
                },
                colour: config.strokeColour
            }
        }
        sock.send(JSON.stringify(drawingData));
        [config.old.x, config.old.y] = [event.clientX, event.clientY];
    }
}

resize();

window.addEventListener("resize", resize);
canvas.addEventListener("mousemove", (event) => config.mouseDown ? draw(event, null) : false);
canvas.addEventListener("mousedown", () => config.mouseDown = true);
canvas.addEventListener("mouseup", () => {
    config.mouseDown = false;
    [config.old.x, config.old.y] = [null, null];
});
colourSelect.addEventListener("change", (event) => config.strokeColour = event.target.value);