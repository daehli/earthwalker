// TODO: This file is getting out of hand.
//
// StreetViewService return service:
// {
//   "location": {
//     "latLng": LatLng,
//     "description": string,
//     "pano": string
//   },
//   "copyright": string,
//   "links": [{
//       "heading": number,
//       "description": string,
//       "pano": string,
//       "roadColor": string,
//       "roadOpacity": number
//     }],
//   "tiles": {
//     "worldSize": Size,
//     "tileSize": Size,
//     "centerHeading": number
//   }
// }

let debug = false;

const PANO_SEARCH_RADIUS = 50000;
const LAT_LIMIT = 85; // polar panos are discarded, they're usually garbage
// string formatting in javascript...
const NOMINATIM_URL = (locStringEncoded) => `https://nominatim.openstreetmap.org/search?q=${locStringEncoded}&polygon_geojson=1&limit=5&format=json`;

let streetViewService = new google.maps.StreetViewService();

// loadGeoTiff gets awaited at the start, so geoImage is always present.
let tiff = undefined;
// let geoImage = undefined;
// let geoData = undefined;
async function loadGeoTiff() {
	const response = await fetch("/static/nasa_pop_data.tif");
	const arrayBuffer = await response.arrayBuffer();
	tiff = await GeoTIFF.fromArrayBuffer(arrayBuffer);
	// geoImage = await tiff.getImage();
	// [geoData] = await geoImage.readRasters();
}
async function getLocationPopulation(lat, lng) {
	const delta = 0.1;
	let value = await tiff.readRasters({
		bbox: [lng, lat, lng + 10 * delta, lat + 10 * delta],
		resX: delta,
		resY: delta,
	});
	let actualValue = value[0][0];
	// 255 means ocean
	if (actualValue == 255) {
		actualValue = 0;
	}
	return actualValue / 255;
}

let pageMapInfo = {
	"numRounds": 5,
	"locStrings": [],
	"locPolygon": null,
	"panoReqs": {
		"panoConnectedness": "always",
		"populationMin": 0.15,
		"populationMax": 1,
	},
	"panoCoords": []
}

let numPanoFetchesInProgress = 0;

let previewMap = null;
let markerGroup = null; // DEBUGGING: map layer group for place markers
let polygonGroup = null; // map layer group for polygon regions

// given a turf.polygon or turf.multiPolygon,
// display it on the map, and fit the map to its bounds
function showPolygonOnMap(map, polygon) {
	let map_poly = L.geoJSON(polygon).addTo(polygonGroup);
	if (debug) {
		console.log(map_poly.getBounds());
	}
	map.fitBounds(map_poly.getBounds());
}

// given a location string, request a polygon from nominatim
// then, update from the form inputs and start looking for places 
// TODO: this isn't great
// TODO: support multiple char separated strings (combine into a multipoly)
function fetchPolygonFromLocString(mapInfo) {
	locString = mapInfo["locStrings"][0]; // TODO: multiple strings (see above)
	// return null if locString is falsey/empty string
	// (handled in getRandomLatLngInPolygon())
	if (locString === "" || !locString) {
		mapInfo["locPolygon"] = null;
		numberOfRoundsUpdated();
		connectedOnlyUpdated();
		fetchPanos(mapInfo);
		return;
	}

	const Http = new XMLHttpRequest();
	const url = NOMINATIM_URL(encodeURI(locString.replace(" ", "+")));
	Http.open("GET", url);
	Http.send();

	// TODO: this is insane, improve async flow
	Http.onreadystatechange = (event) => {
		if (Http.readyState == 4) {
			let placesPolygon;
			let response = JSON.parse(Http.responseText);
			let errorDialog = document.getElementById("error-dialog");

			for (let i = 0; i < response.length; i++) {
				let type = response[i]["geojson"]["type"].toLowerCase();
				if (type === "multipolygon") {
					if (debug) {
						console.log(response[i]);
					}
					errorDialog.setAttribute("hidden", "hidden");
					placesPolygon = turf.multiPolygon(response[i]["geojson"]["coordinates"]);
					break;
				} else if (type === "polygon") {
					if (debug) {
						console.log(response[i]);
					}
					errorDialog.setAttribute("hidden", "hidden");
					placesPolygon = turf.multiPolygon([response[i]["geojson"]["coordinates"]]);
					break;
				}
			}

			if (!placesPolygon) {
				console.log("No polygon response recieved");
				errorDialog.removeAttribute("hidden");
				return;
			}

			showPolygonOnMap(previewMap, placesPolygon);
			mapInfo["locPolygon"] = placesPolygon;
			numberOfRoundsUpdated();
			connectedOnlyUpdated();
			fetchPanos(mapInfo);
		}
	}
}

// ===== API/Panorama Fetching =====

// mapInfo object format (type and default in parens):
/*
{
	"numRounds": (int > 0, default to 5)
	"locStrings": (array of strings, default to [])
	"locPolygon": (turf.multiPolygon, default to null TODO: convert polygon to multiPolygon)
	"panoReqs": {
		"panoConnectedness": (string, one of ["always", "never", "any"], default to "any")
		"populationMin": (between 0 and 1, default 0),
		"populationMax": (between 0 and 1, default 1),
	}
	// TODO: consider storing addition pano information, such as connectedness
	"panoCoords": (array of google.maps.LatLng, default to [])
}
*/
function fetchPanos(mapInfo) {
	disableSubmitButton();

	if (mapInfo["panoCoords"] == null) {
		mapInfo["panoCoords"] = [];
	}

	if (mapInfo["panoCoords"].length + numPanoFetchesInProgress < mapInfo["numRounds"]) {
		for (let i = mapInfo["panoCoords"].length + numPanoFetchesInProgress; i < mapInfo["numRounds"]; i++) {
			numPanoFetchesInProgress += 1; // TODO: still a race condition here
			fetchPano(mapInfo);
		}
	} else {
		// re-enables the submit button if fetchPano never needed to be called (ugh)
		updateSecretForm(mapInfo["panoCoords"], mapInfo["numRounds"]);
	}
}

// fetch a pano and add it to mapInfo["panoCoords"]
// api query is repeated until a good pano is found
// TODO: I think we've ended up with excessive numRounds checks here, try to clean it up
async function fetchPano(mapInfo) {
	let randomLatLng = getRandomLatLngInPolygon(mapInfo["locPolygon"]);

	let handleFail = function() {
		// user may have decreased numRounds, if so don't make another request
		if (mapInfo["panoCoords"].length < mapInfo["numRounds"]) {
			setTimeout(function() {
				fetchPano(mapInfo);
			}, 0);
		} else {
			updateSecretForm(mapInfo["panoCoords"], mapInfo["numRounds"]);
		}
	};

	function handlePanoResponse(result, status) {
		if (status == google.maps.StreetViewStatus.OK && resultPanoIsGood(result, mapInfo["panoReqs"], mapInfo["locPolygon"])) {
			if (debug) {
				L.marker([result.location.latLng.lat(), result.location.latLng.lng()]).addTo(markerGroup); // DEBUGGING: show selected places on map
			}

			// in case the user has decreased numRounds while the request was running, don't add the pano
			if (mapInfo["panoCoords"].length < mapInfo["numRounds"]) {
				mapInfo["panoCoords"].push(result.location.latLng);
			}
			numPanoFetchesInProgress -= 1; // TODO: still a race condition here
			updateFetchingBar(mapInfo["panoCoords"], mapInfo["numRounds"]);
			updateSecretForm(mapInfo["panoCoords"], mapInfo["numRounds"]);
		} else {
			console.log("Failed to get location; api request: " + status.toString());
			handleFail();
		}
	}

	let population = await getLocationPopulation(randomLatLng.lat(), randomLatLng.lng());
	let min = mapInfo["panoReqs"]["populationMin"];
	let max = mapInfo["panoReqs"]["populationMax"];
	if (population >= min && population <= max) {
		streetViewService.getPanoramaByLocation(randomLatLng, PANO_SEARCH_RADIUS, handlePanoResponse);
	} else {
		console.log("Failed, population not right: " + population);
		handleFail();
	}
}

// returns whether result (pano) meets the requirements of mapInfo
function resultPanoIsGood(result, panoReqs, polygon) {
	if (result.location.latLng.lat() > LAT_LIMIT || result.location.latLng.lat() < -1 * LAT_LIMIT) {return false;}

	if (panoReqs["panoConnectedness"] === "always" && result.links.length == 0) {
		return false;
	}
	if (panoReqs["panoConnectedness"] === "never" && result.links.length > 0) {
		return false;
	}

	let locationTurfPoint = turf.point([result.location.latLng.lng(), result.location.latLng.lat()]);
	if (polygon != null && !turf.booleanPointInPolygon(locationTurfPoint, polygon)) {
		return false;
	}

	return true;
}

// =====

function disableSubmitButton() {
	let button = document.getElementById("submit-button");
	button.setAttribute("disabled", "disabled");
}

// update loading/fetching progress bar with number of panoCoords found
function updateFetchingBar(panoCoords, numRounds) {
	document.getElementById("loading-progress").setAttribute("style", "width: " + ((100 * panoCoords.length) / numRounds) + "%;");
}

// put panoCoords into the hidden form input
// TODO: this is a hack
// re-enables the submit button
function updateSecretForm(panoCoords, numRounds) {
	if (panoCoords.length >= numRounds) {
		if (panoCoords.length > numRounds) {
			console.warn("Too many panoCoords?! mapInfo:");
			console.log(pageMapInfo); // DEBUGGING: should probably remove use of this global
		}
		let input = document.getElementById("hidden-input");
		let button = document.getElementById("submit-button");
		input.setAttribute("value", JSON.stringify(panoCoords));
		button.removeAttribute("disabled");
	}
}

// get a random google.maps.LatLng, anywhere
function getRandomLatLng() {
	let randomLng = (Math.random() * 360 - 180);
	let randomLat = (Math.random() * 180 - 90);
	return new google.maps.LatLng(randomLat, randomLng);
}

// get a random google.maps.LatLng within the specified turf.polygon or turf.multiPolygon
function getRandomLatLngInPolygon(polygon) {
	if (polygon == null) {
		// fall back to global
		return getRandomLatLng();
	}
	bounds = turf.bbox(polygon);
	let randomLng;
	let randomLat;
	let lnglat;
	// TODO: more efficient algorithm? - suffices for the small number of points needed
	do { 
		randomLng = (Math.random() * (bounds[2] - bounds[0]) + bounds[0]);
		randomLat = (Math.random() * (bounds[3] - bounds[1]) + bounds[1]);
		lnglat = turf.point([randomLng, randomLat]);
	} while (!turf.booleanPointInPolygon(lnglat, polygon))
	//L.marker([randomLat, randomLng]).addTo(markerGroup); // DEBUGGING: show _all_ random points on map
	return new google.maps.LatLng(randomLat, randomLng);
}

// ===== Form Change Handlers =====

function numberOfRoundsUpdated() {
	let newNumRounds = document.getElementById("rounds").value;
	if (!newNumRounds) {
		return;
	}
	if (newNumRounds < pageMapInfo["numRounds"]) {
		// note: can't decrease length of panoCoords beyond 0, so any excess requests are handled in fetchPano()
		pageMapInfo["panoCoords"] = pageMapInfo["panoCoords"].slice(0, newNumRounds);
	}
	pageMapInfo["numRounds"] = newNumRounds;
	fetchPanos(pageMapInfo);
}

function connectedOnlyUpdated() {
	// TODO: improve user-friendliness of these values
	let newConnectedOnly = document.getElementById("connectedOnly").value;
	if (pageMapInfo["panoReqs"]["panoConnectedness"] !== newConnectedOnly) {
		disableSubmitButton();
		pageMapInfo["panoReqs"]["panoConnectedness"] = newConnectedOnly;
		pageMapInfo["panoCoords"] = []; // TODO: considering storing pano connectedness and only removing as necessary
		markerGroup.clearLayers(); // DEBUGGING: clear markers
		fetchPanos(pageMapInfo);
	}
}

// TODO: support multiple loc strings
function locStringUpdated() {
	let old = pageMapInfo["locStrings"][0];
	let newLocString = document.getElementById("locString").value;
	if (old !== newLocString) {
		pageMapInfo["locStrings"][0] = newLocString;
		disableSubmitButton();
		previewMap.setView([0, 0], 1);
		pageMapInfo["panoCoords"] = [];
		polygonGroup.clearLayers();
		markerGroup.clearLayers();
		pageMapInfo["locPolygon"] = fetchPolygonFromLocString(pageMapInfo);
	}
}

function popDensityUpdated() {
	let newMin = document.getElementById("minDensity").value / 100;
	let newMax = document.getElementById("maxDensity").value / 100;
	let oldMin = pageMapInfo["panoReqs"]["populationMin"];
	let oldMax = pageMapInfo["panoReqs"]["populationMax"];
	if (newMin !== oldMin || newMax !== oldMax) {
		pageMapInfo["panoReqs"]["populationMin"] = newMin;
		pageMapInfo["panoReqs"]["populationMax"] = newMax;
		disableSubmitButton();
		previewMap.setView([0, 0], 1);
		pageMapInfo["panoCoords"] = [];
		polygonGroup.clearLayers();
		markerGroup.clearLayers();
		fetchPanos(pageMapInfo);
	}
}

// settings may have been cached by the browser (wouldn't trigger the onchange),
// so check them once the DOM has loaded
window.addEventListener("DOMContentLoaded", (event) => {
	// TODO: stick map stuff in a function
	let load = async function() {
		await loadGeoTiff();
		previewMap = L.map("bounds-map", {center: [0, 0], zoom: 1});
		L.tileLayer("https://tiles.wmflabs.org/osm/{z}/{x}/{y}.png", {
			attribution: '&copy; <a href="https://www.openstreetmap.org/copyright">OSM</a> contributors, <a href="https://wikitech.wikimedia.org/wiki/Wikitech:Cloud_Services_Terms_of_use">Wikimedia Cloud Services</a>'
		}).addTo(previewMap);
		markerGroup = L.layerGroup().addTo(previewMap);
		polygonGroup = L.layerGroup().addTo(previewMap);
		numberOfRoundsUpdated();
		connectedOnlyUpdated();
		locStringUpdated();
		popDensityUpdated();
	};
	load();
});
