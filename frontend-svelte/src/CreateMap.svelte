<script>
    // TODO: svelteify this file
    import {onMount} from 'svelte';
    const NOMINATIM_URL = (locStringEncoded) => `https://nominatim.openstreetmap.org/search?q=${locStringEncoded}&polygon_geojson=1&limit=5&polygon_threshold=0.005&format=json`;

    let mapSettings = {
        Name: "",
        Polygon: null,
        Area: 0,
        NumRounds: 0,
        TimeLimit: 0,
        GraceDistance: 10,
        MinDensity: 0,
        MaxDensity: 100,
        Connectedness: 0,
        Copyright: 0,
        Source: 0,
        ShowLabels: true
    };
    let locString = "";
    let previewMap;
    let previewPolyGroup;
    let advancedHidden = true;
    $: window.globalMap = previewMap;

    onMount(async () => {
        previewMap = L.map("bounds-map", {center: [0, 0], zoom: 1});
        let tileServer = await getTileServer();
        L.tileLayer(tileServer, {
            attribution: '&copy; <a href="https://www.openstreetmap.org/copyright">OSM</a> contributors, <a href="https://wikitech.wikimedia.org/wiki/Wikitech:Cloud_Services_Terms_of_use">Wikimedia Cloud Services</a>'
        }).addTo(previewMap);
        previewPolyGroup = L.layerGroup().addTo(previewMap);
    });

    // collates createmap form data into a JSON object, 
    // then sends a newmap request to the server
    function handleFormSubmit() {
        mapSettings.Name          = strById("Name");
        mapSettings.NumRounds     = intById("NumRounds");
        mapSettings.GraceDistance = intById("GraceDistance");
        mapSettings.MinDensity    = intById("MinDensity");
        mapSettings.MaxDensity    = intById("MaxDensity");
        mapSettings.Connectedness = intById("Connectedness");
        mapSettings.Copyright     = intById("Copyright");
        mapSettings.Source        = intById("Source");
        
        let showLabelsInput = document.getElementById("ShowLabels");
        if (showLabelsInput) {
            mapSettings.ShowLabels = showLabelsInput.checked;
        }

        // read total TimeLimit
        mapSettings.TimeLimit = 0;
        mapSettings.TimeLimit += 60 * intById("TimeLimit_minutes");
        mapSettings.TimeLimit += intById("TimeLimit_seconds");

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
        fetch("/api/maps", {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify(mapSettings),
        }).then(console.log("mapSettings sent to server"));
        // TODO: redirect to createchallenge
    }

    function intById(id, fallback=0) {
        let input = document.getElementById(id);
        if (input) {
            if (!input.value) {
                return fallback;
            }
            return parseInt(input.value, 10);
        } else {
            console.log("Couldn't find input '" + id + "', using fallback.")
            return fallback;
        }
    }

    function strById(id, fallback="") {
        let input = document.getElementById(id);
        if (input) {
            return input.value;
        } else {
            console.log("Couldn't find input '" + id + "', using fallback.")
            return fallback;
        }
    }

    async function getTileServer() {
        let response = await fetch("/api/config/tileserver", {
            method: "GET",
        });
        let data = await response.json();
        return data.tileserver;
    }

    function locStringUpdated() {
        let old = locString;
        let locStringInput = document.getElementById("locString");
        if (locStringInput) {
            locString = document.getElementById("locString").value;
        }
        if (old !== locString) {
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
        console.log("getting geojson...");
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
                <input type="text" class="form-control" id="Name"/>
            </div>
        </div>

        <div class="form-group">
            <div class="input-group">
                <div class="input-group-prepend">
                    <div class="input-group-text">Number of Rounds</div>
                </div>
                <input type="number" class="form-control" id="NumRounds" value="5" min="1" max="100"/>
            </div>
        </div>

        <div class="form-row">
            <div class="col">
                <div class="input-group">
                    <div class="input-group-prepend">
                        <div class="input-group-text">Round Time, Minutes</div>
                    </div>
                    <input type="number" min="0" class="form-control mr-sm-3" id="TimeLimit_minutes"/>
                </div>
            </div>
            <div class="col">
                <div class="input-group">
                    <div class="input-group-prepend">
                        <div class="input-group-text">Seconds</div>
                    </div>
                    <input type="number" min="0" class="form-control" id="TimeLimit_seconds"/>
                </div>
            </div>
        </div>
        <small class="form-text text-muted">
            Leave empty or zero for no time limit.
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
                        <input type="number" class="form-control" id="GraceDistance" value="10" min="0"/>
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
                            <input type="number" class="form-control mr-sm-3" id="MinDensity" value="15" min="0" max="100"/>
                        </div>
                    </div>
                    <div class="col">
                        <div class="input-group">
                            <div class="input-group-prepend">
                                <div class="input-group-text">Maximum</div>
                            </div>
                            <input type="number" class="form-control mr-sm-3" id="MaxDensity" value="100" min="0" max="100"/>
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
                        <select class="form-control" id="Connectedness">
                            <option value=1 selected="selected">always</option>
                            <option value=2 >never</option>
                            <option value=0 >any</option>
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
                        <select class="form-control" id="Copyright">
                            <option value=0 selected="selected">any</option>
                            <option value=1>Google only</option>
                            <option value=2>third party only</option>
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
                        <select class="form-control" id="Source">
                            <option value=1 selected="selected">outdoors only</option>
                            <option value=0 >any</option>
                        </select>
                    </div>
                </div>
                <small class="form-text text-muted">
                    If you want to exclude panoramas inside businesses.
                </small>

                <hr/>

                <div class="form-group">
                    <div class="form-check form-check-inline">
                        <input class="form-check-input" type="checkbox" id="ShowLabels" checked>
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
                        <input type="text" class="form-control mr-sm-3" id="locString" placeholder="Location" on:change={locStringUpdated}/>
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