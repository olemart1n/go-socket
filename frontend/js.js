const chatroomForm = document.querySelector("#chatroom-selection");
const messageForm = document.querySelector("#chatroom-message");
const socket = new WebSocket("ws://" + "localhost:8080" + "/ws");

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

    sendMessage();
});

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

function sendMessage() {
    const newMessage = document.querySelector("#message");
    sendEvent("send_message", newMessage.value);
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

function sendEvent(eName, payload) {
    const event = new Event(eName, payload);
    socket.send(JSON.stringify(event));
}
