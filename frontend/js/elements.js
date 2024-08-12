/**
 * @type {HTMLFormElement}
 */
const chatroomForm = document.querySelector("#chatroom-selection");

/**
 * @type {HTMLFormElement}
 */
const messageForm = document.querySelector("#chatroom-message");

/**
 * @type {HTMLFormElement}
 */
const loginForm = document.querySelector("#login-form");

/**
 * @type {HTMLDivElement}
 */
const connectionStatus = document.querySelector("#connection-status");

/**
 * @type {HTMLTextareaElement}
 */
const textArea = document.querySelector("#message-area");
/**
 * @type {HTMLButtonElement}
 */
const helperButton = document.querySelector("#helper-btn");

export const elements = {
    chatroomForm,
    messageForm,
    connectionStatus,
    textArea,
    loginForm,
    helperButton,
};
