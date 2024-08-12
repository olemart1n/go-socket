import { elements } from "./elements.js";
import { login } from "./login.js";
import { state } from "./state.js";
import { addSocketEventListeners } from "./socket.js";
let socket;
elements.loginForm.addEventListener("submit", async (e) => {
    e.preventDefault();
    const data = Object.fromEntries(new FormData(e.target));
    const otp = await login(data);
    console.log("index.js: " + otp);
    socket = new WebSocket("wss://" + document.location.host + "/ws?otp=" + otp);
    addSocketEventListeners(socket);
});

elements.messageForm.addEventListener("submit", (e) => {
    e.preventDefault();
    const formdata = new FormData(e.target);
    const data = Object.fromEntries(formdata.entries());
    data.from = "percy";
    socket.send(JSON.stringify({ type: "send_message", payload: data }));
    elements.messageForm.reset();
});
elements.connectionStatus.innerHTML = "Connected to websocket: False";
