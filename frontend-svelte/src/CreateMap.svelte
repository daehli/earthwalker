<script>
    // TODO: svelteify this file
    import {onMount} from 'svelte';

    const NOMINATIM_URL = (locStringEncoded) => `https://nominatim.openstreetmap.org/search?q=${locStringEncoded}&polygon_geojson=1&limit=5&polygon_threshold=0.005&format=json`;

    let ewapi = new EarthwalkerAPI();

    let mapSettings = {
        Name: "",
        Polygon: null,
        Area: 0,
        NumRounds: 5,
        TimeLimit: 0,
        GraceDistance: 10,
        MinDensity: 15,
        MaxDensity: 100,
        Connectedness: 1,
        Copyright: 0,
        Source: 1,
        ShowLabels: true
    };
    // extra bindings (handleFormSubmit converts these to mapSettings fields)
    let timeLimitMinutes = 0;
    let timeLimitSeconds = 0;

    let locString = "";
    let oldLocString = "";
    let previewMap;
    let previewPolyGroup;
    let advancedHidden = true;

    onMount(async () => {
        previewMap = L.map("bounds-map", {center: [0, 0], zoom: 1});
        let tileServer = (await ewapi.getTileServer()).tileserver;
        L.tileLayer(tileServer, {
            attribution: '&copy; <a href="https://www.openstreetmap.org/copyright">OSM</a> contributors, <a href="https://wikitech.wikimedia.org/wiki/Wikitech:Cloud_Services_Terms_of_use">Wikimedia Cloud Services</a>'
        }).addTo(previewMap);
        previewPolyGroup = L.layerGroup().addTo(previewMap);
    });

    // collates createmap form data into a JSON object, 
    // then sends a newmap request to the server
    function handleFormSubmit() {
        // calculate total TimeLimit
        mapSettings.TimeLimit = 60 * timeLimitMinutes + timeLimitSeconds;

        // sanity check density fields
        // TODO: nicer error messages than alerts
        // TODO: check that population density in Polygon overlaps with the
        //       specified range (otherwise we'll never be able to find good
        //       panos.)
        if (mapSettings.MinDensity > mapSettings.MaxDensity) {
            alert("Max density must be greater than min density.");
            return;
        }
        // TODO: evaluate challenge generation (to make sure mapSettings aren't so
        //       specific that it takes a huge number of API requests to find good
        //       panos)
        // send new map to server
        ewapi.postMap(mapSettings)
            .then( (response) => {
                if (response && response.MapID) {
                    console.log("mapSettings sent to server");
                    window.location.replace("/createchallenge?mapid="+response.MapID);
                } else {
                    alert("Failed to submit map?!");
                }
            });
    }

    function handleLocStringUpdate() {
        if (locString != oldLocString) {
            oldLocString = locString;
            updatePolygonFromLocString();
        }
    }

    function showPolygonOnMap() {
        previewPolyGroup.clearLayers();
        let map_poly = L.geoJSON(mapSettings.Polygon).addTo(previewPolyGroup);
        previewMap.fitBounds(map_poly.getBounds());
    }

    function updatePolygonFromLocString() {
        if (locString === "" || !locString) {
            mapSettings.Polygon = null;
            return;
        }
        
        fetch(NOMINATIM_URL(encodeURI(locString.replace(" ", "+"))))
            .then(response => response.json())
            .then(data => {
                mapSettings.Polygon = geojsonFromNominatim(data);
                mapSettings.Area = turf.area(mapSettings.Polygon);
                showPolygonOnMap();
            });
    }

    // given Nominatim results, takes the most significant one with a polygon or
    // multipolygon and returns it as a turf.multiPolygon
    function geojsonFromNominatim(data) {
        for (let i = 0; i < data.length; i++) {
            let type = data[i].geojson.type.toLowerCase();
            if (type === "multipolygon") {
                return turf.multiPolygon(data[i].geojson.coordinates);
            } else if (type === "polygon") {
                return turf.multiPolygon([data[i].geojson.coordinates]);
            }
        }
        console.log("No matching polygon!");
        return null;
    }
</script>

<main>
    <div class="container">

    <br>

    <h2>Create a New Map</h2>

    <br>

    <form on:submit|preventDefault={handleFormSubmit} method="post">

        <div class="form-group">
            <div class="input-group">
                <div class="input-group-prepend">
                    <div class="input-group-text">Map Name</div>
                </div>
                <input type="text" class="form-control" id="Name" required bind:value={mapSettings.Name}/>
            </div>
        </div>

        <div class="form-group">
            <div class="input-group">
                <div class="input-group-prepend">
                    <div class="input-group-text">Number of Rounds</div>
                </div>
                <input type="number" class="form-control" id="NumRounds" bind:value={mapSettings.NumRounds} min="1" max="100"/>
            </div>
        </div>

        <div class="form-row">
            <div class="col">
                <div class="input-group">
                    <div class="input-group-prepend">
                        <div class="input-group-text">Round Time, Minutes</div>
                    </div>
                    <input type="number" min="0" class="form-control mr-sm-3" id="TimeLimit_minutes" bind:value={timeLimitMinutes}/>
                </div>
            </div>
            <div class="col">
                <div class="input-group">
                    <div class="input-group-prepend">
                        <div class="input-group-text">Seconds</div>
                    </div>
                    <input type="number" min="0" class="form-control" id="TimeLimit_seconds" bind:value={timeLimitSeconds}/>
                </div>
            </div>
        </div>
        <small class="form-text text-muted">
            Leave zero for no time limit.
        </small>

        <br/>

        <div class="card border-info">
            <div class="card-header">
                <button class="btn btn-info" type="button" on:click={() => {advancedHidden = !advancedHidden; setTimeout(function() {previewMap.invalidateSize()}, 400)}}>
                    Show advanced settings
                </button>
            </div>

            <div class="card-body" id="advanced-settings" hidden={advancedHidden}>
                <div class="form-group">
                    <div class="input-group">
                        <div class="input-group-prepend">
                            <div class="input-group-text">Grace Distance (m)</div>
                        </div>
                        <input type="number" class="form-control" id="GraceDistance" bind:value={mapSettings.GraceDistance} min="0"/>
                    </div>
                </div>
                <small class="form-text text-muted">
                    Guesses within this distance (in meters) will be awarded full points.
                </small>
                <hr/>
                <!-- TODO: it would be nice if this was a double range slider -->
                <div class="form-row">
                    <div class="col">
                        <div class="input-group">
                            <div class="input-group-prepend">
                                <div class="input-group-text">Population Density %, Minimum</div>
                            </div>
                            <input type="number" class="form-control mr-sm-3" id="MinDensity" bind:value={mapSettings.MinDensity} min="0" max="100"/>
                        </div>
                    </div>
                    <div class="col">
                        <div class="input-group">
                            <div class="input-group-prepend">
                                <div class="input-group-text">Maximum</div>
                            </div>
                            <input type="number" class="form-control mr-sm-3" id="MaxDensity" bind:value={mapSettings.MaxDensity} min="0" max="100"/>
                        </div>
                    </div>
                </div>
                <small class="form-text text-muted">
                    0% is ocean. 10% is barren road. With 20%, you will find signs of civilization. Anything above 50% is already very populated.
                </small>

                <hr/>

                <div class="form-group">
                    <div class="input-group">
                        <div class="input-group-prepend">
                            <div class="input-group-text">Panorama connectedness</div>
                        </div>
                        <!-- note: select values are Object.  Wrapping them in
                                   brackets takes advantage of object init
                                   shorthand to give us ints instead of strings.
                                   However! The resulting binding is not 
                                   bidirectional, so make sure your mapSettings
                                   defaults match the select defaults. -->
                        <select class="form-control" id="Connectedness" bind:value={mapSettings.Connectedness}>
                            <option value={1} selected="selected">always</option>
                            <option value={2} >never</option>
                            <option value={0} >any</option>
                        </select>
                    </div>
                </div>
                <small class="form-text text-muted">
                    If you want to be able to always walk somewhere or if you want single-image ones. 
                </small>

                <hr/>

                <div class="form-group">
                    <div class="input-group">
                        <div class="input-group-prepend">
                            <div class="input-group-text">Copyright</div>
                        </div>
                        <select class="form-control" id="Copyright" bind:value={mapSettings.Copyright}>
                            <option value={0} selected="selected">any</option>
                            <option value={1}>Google only</option>
                            <option value={2}>third party only</option>
                        </select>
                    </div>
                </div>
                <small class="form-text text-muted">
                    If you want to see only Google panos or also include third party panos.
                </small>

                <hr/>

                <div class="form-group">
                    <div class="input-group">
                        <div class="input-group-prepend">
                            <div class="input-group-text">Source</div>
                        </div>
                        <select class="form-control" id="Source" bind:value={mapSettings.Source}>
                            <option value={1} selected="selected">outdoors only</option>
                            <option value={0}>any</option>
                        </select>
                    </div>
                </div>
                <small class="form-text text-muted">
                    If you want to exclude panoramas inside businesses.
                </small>

                <hr/>

                <div class="form-group">
                    <div class="form-check form-check-inline">
                        <input class="form-check-input" type="checkbox" id="ShowLabels" bind:checked={mapSettings.ShowLabels}>
                        <label class="form-check-label" for="label">Show labels on map</label>
                    </div>
                </div>
                <small class="form-text text-muted">
                    Check this if the map should tell you how places are called.
                </small>

                <hr/>
                
                <div class="form-group">
                    <div class="input-group">
                        <div class="input-group-prepend">
                            <div class="input-group-text">Location string </div>
                        </div>
                        <input type="text" class="form-control mr-sm-3" id="locString" placeholder="Location" bind:value={locString} on:change={handleLocStringUpdate}/>
                    </div>
                    <small class="form-text text-muted">
                        Constrain the game to a specified area - enter a country, state, city, neighborhood, lake, or any other bounded area.  Does not yet affect scoring.
                    </small>
                    <div class="card bg-danger text-white mt-1" id="error-dialog" hidden>
                        <p class="card-text">Sorry, that does not seem like a valid bounding box on OSM Nominatim.</p>
                    </div>
                </div>
                <div id="bounds-map" style="width: 80%; height: 50vh; margin-left: 10%; margin-right: 10%;"></div>
            </div>
        </div>

        <br/>

        <input id="hidden-input" type="hidden" name="result" value=""/>

        <button id="submit-button" type="submit" class="btn btn-primary" style="margin-bottom: 2em;">Create Map</button>

    </form>
    <link rel="stylesheet" href="static/leaflet/leaflet.css"/>
    </div>
</main>