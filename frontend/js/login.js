/**
 * @async
 * @function
 * @param {any} credentials
 * @returns {Promise<string>} OTP
 */
export const login = async (credentials) => {
    const res = await fetch("/login", {
        method: "POST",
        body: JSON.stringify(credentials),
        mode: "cors",
    });
    if (!res.ok) {
        alert("unauthorized");
        return;
    }
    try {
        const json = await res.json();
        return json.OTP;
    } catch (error) {
        console.error(error);
    }
};
