const chatroomForm = document.querySelector("#chatroom-selection");
const messageForm = document.querySelector("#chatroom-message");
const loginForm = document.querySelector("#login-form");
const connectionStatus = document.querySelector("#connection-status");

class Event {
    constructor(type, payload) {
        this.type = type;
        this.payload = payload;
    }
}

function changeChatroom() {
    const newChat = document.querySelector("#chatroom");
    if (newChat == !null && newChat.value !== selectedChat) {
        console.log(newChat);
    }
    return false;
}

function sendMessage(socket) {
    const newMessage = document.querySelector("#message");
    sendEvent("send_message", newMessage.value, socket);
}

function routeEvent(e) {
    if (e.type === undefined) {
        alert("no type field in the event");
    }
    switch (e.type) {
        case "new_message":
            console.log("new message");
            break;

        default:
            alert("undupported message");
            break;
    }
}

function sendEvent(eName, payload, socket) {
    const event = new Event(eName, payload);
    socket.send(JSON.stringify(event));
}

loginForm.addEventListener("submit", async (e) => {
    e.preventDefault();
    const data = Object.fromEntries(new FormData(e.target));

    const res = await fetch("/login", {
        method: "POST",
        body: JSON.stringify(data),
        mode: "cors",
    });
    if (!res.ok) {
        alert("unauthorized");
        return;
    }
    try {
        const json = await res.json();
        console.log(json);
        connectWebsocket(json.OTP);
    } catch (error) {
        console.error(error);
    }
});

function connectWebsocket(otp) {
    // IN PRODUCTION THE SOCKET PROTOCOL SHOULD BE SECURE (WSS)
    console.log(document.location.host);
    console.log(window.location.hostname);
    const socket = new WebSocket("wss://" + document.location.host + "/ws?otp=" + otp);

    socket.addEventListener("open", () => {
        connectionStatus.innerHTML = "Connected to websocket: True";
    });
    socket.addEventListener("close", () => {
        connectionStatus.innerHTML = "Connected to websocket: false";
    });
    socket.addEventListener("message", (e) => {
        const eventData = JSON.parse(e.data);
        const event = Object.assign(new Event(), eventData);
        routeEvent(event);
    });

    chatroomForm.addEventListener("submit", (e) => {
        e.preventDefault();
        changeChatroom();
    });
    messageForm.addEventListener("submit", (e) => {
        e.preventDefault();
        sendMessage(socket);
    });
}
