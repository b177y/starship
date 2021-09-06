import { writable } from 'svelte/store';

const storedLoggedIn = localStorage.getItem("loggedIn");
const storedAuthToken = localStorage.getItem("authToken");
console.log("Getting localStorage loggedIn: ", storedLoggedIn)

var sl = (storedLoggedIn === null ? false : JSON.parse(storedLoggedIn))
var at = (storedAuthToken === null ? "" : storedAuthToken)

console.log("writable using", sl)

export const LoggedIn = writable(sl);
export const AuthToken = writable(at);

LoggedIn.subscribe(value => {
    console.log("Setting localStorage loggedIn to", value)
    localStorage.setItem("loggedIn", value)
})

AuthToken.subscribe(value => {
    localStorage.setItem("authToken", value)
})

