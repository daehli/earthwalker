import { writable } from 'svelte/store';

function createLoc() {
    const {subscribe, set, update} = writable("");
    return {
        subscribe,
        set: (href) => {
            history.pushState({}, "", href);
            set(href);
        },
        write: (href) => {
            set(href);
        }
    }
}
export const loc = createLoc();

export const globalMap = writable(null);
export const globalChallenge = writable(null);
export const globalResult = writable(null);