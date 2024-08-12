// export function changeChatroom() {
//     const newChat = document.querySelector("#chatroom");
//     if (newChat == !null && newChat.value !== selectedChat) {
//         console.log(newChat);
//     }
//     return false;
// }

chatroomForm.addEventListener("submit", (e) => {
    e.preventDefault();
    changeChatroom();
});
messageForm.addEventListener("submit", (e) => {
    e.preventDefault();
    sendMessage(socket);
});
