import { elements } from "./elements.js";

/**
 * @typedef {Object} State
 * @property {string} _connectionStatusMessage - The internal status message.
 * @property {string} connectionStatusMessage - The public status message (defined via Object.defineProperty).
 */

/**
 * @type {State}
 */
export const state = {
    _connectionStatusMessage: "Connected to websocket: False",
};

Object.defineProperty(state, "connectionStatusMessage", {
    get: function () {
        return this._connectionStatusMessage;
    },
    set: function (value) {
        this._connectionStatusMessage = value;
        elements.connectionStatus.innerHTML = value;
        console.log("connectionStatusMessage updated:", value);
    },
});
