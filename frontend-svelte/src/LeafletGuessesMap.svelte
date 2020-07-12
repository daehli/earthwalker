<script>
    import {onMount} from 'svelte';

    export let ewapi, curMap, curChallenge, displayedResult, showAll;

    let tileServer;

    let mapDiv;

    let lMap;
    let polyGroup;
    let guessGroup;

    $: if (guessGroup && displayedResult) {
        if (showAll) {
            showGuesses(guessGroup, displayedResult.Guesses);
        } else {
            showGuesses(guessGroup, displayedResult.Guesses.slice(-1));
        }
        lMap.fitBounds(guessGroup.getBounds());
    };

    onMount(async () => {
        lMap = new L.Map(mapDiv);
        lMap.setView([0.0, 0.0], 1);

        tileServer = (await ewapi.getTileServer()).tileserver;
        L.tileLayer(tileServer, {
            attribution: '&copy; <a href="https://www.openstreetmap.org/copyright">OSM</a> contributors, <a href="https://wikitech.wikimedia.org/wiki/Wikitech:Cloud_Services_Terms_of_use">Wikimedia Cloud Servides</a>'
        }).addTo(lMap);

        polyGroup = L.layerGroup().addTo(lMap);
        if (curMap.Polygon) {
            showPolygonOnMap(polyGroup, curMap.Polygon);
        }

        guessGroup = L.featureGroup().addTo(lMap);
    });

    function showGuesses(layer, guesses) {
        layer.clearLayers();
        guesses.forEach(guess => {
            showGuessOnMap(layer, guess, curChallenge.Places[guess.RoundNum], guess.RoundNum, displayedResult.Nickname, displayedResult.Icon);
        });
    }
</script>

<style>
    div {
        width: 100%;
        height: 50vh;
    }
</style>

<div bind:this={mapDiv}></div>