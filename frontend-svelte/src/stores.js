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