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

export const curMap = writable({});
export const curChallenge = writable({});
export const curResult = writable({});