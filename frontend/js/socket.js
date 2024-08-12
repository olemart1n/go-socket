import { state } from "./state.js";
import { elements } from "./elements.js";
export class SocketEvents {
    constructor(type, payload) {
        this.type = type;
        this.payload = payload;
    }

    sendMessage(type, message, from) {
        return JSON.stringify({ type: type, payload: { message: message, from: from } });
    }

    receiveMessage(message, from, sent) {
        return { message, from, sent };
    }
}

/**
 *
 * @param {WebSocket} socket
 */
export function addSocketEventListeners(socket) {
    socket.addEventListener("open", () => {
        state.connectionStatusMessage = "Connected to websocket: True";
    });
    socket.addEventListener("close", () => {
        state.connectionStatusMessage = "Connected to websocket: false";
    });
    socket.addEventListener("message", (e) => {
        handleEvent(JSON.parse(e.data));
    });
}

function handleEvent(e) {
    if (e.type === undefined) {
        alert("no type field in the event");
    }
    switch (e.type) {
        case "new_message":
            console.log(e.Payload);

            appendChatMessage(e.Payload);
            break;

        default:
            alert("undupported message");
            break;
    }
}

export function appendChatMessage(payload) {
    elements.textArea.innerHTML = elements.textArea.innerHTML =
        payload.from + ": " + payload.message + "\n";
    elements.textArea.scrollTop = elements.textArea.scrollHeight;
}
