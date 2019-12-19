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

let service = new google.maps.StreetViewService();

let searchingForResults = false;
let results = [];
let numDesiredResults = 5;
let connectedOnly = false;

let map = null;
let markerGroup = null; // DEBUGGING: map layer group for place markers
let polygonGroup = null; // map layer group for polygon regions
let locString = null;
let placesPolygon = null;

// given a turf.polygon or turf.multiPolygon,
// display it on the map, and fit the map to its bounds
function showPolygonOnMap(map, polygon) {
	let map_poly = L.geoJSON(polygon).addTo(polygonGroup);
	map.fitBounds(map_poly.getBounds());
}

// given a location string, request a polygon from nominatim
// then, update from the form inputs and start looking for places TODO: this isn't great
function getPolygonFromLocString(locString) {
	// don't update the polygon if locString is falsey/empty string
	// that's handled in queryPosition()
	if (locString === "" || !locString) {
		numberOfRoundsUpdated();
		connectedOnlyUpdated();
		queryPosition();
		return;
	}
	const Http = new XMLHttpRequest();
	const url = "https://nominatim.openstreetmap.org/search?q=" + encodeURI(locString.replace(" ", "+")) + "&polygon_geojson=1&limit=1&format=json";
	Http.open("GET", url);
	Http.send();

	// TODO: this is insane, improve async flow
	Http.onreadystatechange = (event) => {
		if (Http.readyState == 4) {
			response = JSON.parse(Http.responseText)[0];
			console.log("Response received, display name: " + response["display_name"]);
			if (response["geojson"]["type"].toLowerCase() === "multipolygon") {
				placesPolygon = turf.multiPolygon(response["geojson"]["coordinates"]);
			} else {
				placesPolygon = turf.polygon(response["geojson"]["coordinates"]);
			}
			showPolygonOnMap(map, placesPolygon);
			numberOfRoundsUpdated();
			connectedOnlyUpdated();
			queryPosition();
		}
	}
}

function getRandomLatLng() {
	randomLng = (Math.random() * 360 - 180);
	randomLat = (Math.random() * 180 - 90);
	return new google.maps.LatLng(randomLat, randomLng);
}

// get a random google.maps.LatLng within the specified turf.polygon or turf.multiPolygon
function getRandomLatLngInPolygon(polygon) {
	bounds = turf.bbox(polygon);
	// TODO: not exactly the height of efficiency, but suffices for the small number of points needed
	do { 
		randomLng = (Math.random() * (bounds[2] - bounds[0]) + bounds[0]);
		randomLat = (Math.random() * (bounds[3] - bounds[1]) + bounds[1]);
		lnglat = turf.point([randomLng, randomLat]);
	} while (!turf.booleanPointInPolygon(lnglat, polygon))
	//L.marker([randomLat, randomLng]).addTo(markerGroup); // DEBUGGING: show random points on map
	return new google.maps.LatLng(randomLat, randomLng);
}

function queryPosition() {
	searchingForResults = true;
	let point;
	if (locString === "" || !locString) {
		point = getRandomLatLng();
	} else {
		point = getRandomLatLngInPolygon(placesPolygon);
	}
	let radius = 10000;
	service.getPanoramaByLocation(point, radius, function(result, status) {
		if (status == google.maps.StreetViewStatus.OK) {
			let nearestLatLng = result.location.latLng;
			// There seems to be a panorama graveyard at the top and bottom
			// of the earth of incorrectly positioned paranoramas.
			// Do not take these incorrect panorams into account.
			// Of course, there is some sacrifice of actually interesting panoramas here.
			console.log(nearestLatLng.lat());
			if (nearestLatLng.lat() < 85 && nearestLatLng.lat() > -85
				&& (!connectedOnly || result.links.length > 0) // exclude unconnected/orphan panos
				// && result.copyright.includes("Google") // For now
			) {
				console.log("num links: " + result.links.length);
				L.marker([nearestLatLng.lat(), nearestLatLng.lng()]).addTo(markerGroup); // DEBUGGING: show selected places on map
				results.push(nearestLatLng);
			}
		} else {
			console.log("Failed to get location: " + status.toString());
		}
		document.getElementById("loading-progress").setAttribute("style", "width: " + ((100 * results.length) / numDesiredResults) + "%;");
		if (results.length < numDesiredResults) {
			queryPosition();
		} else {
			// Yea, this is probably incorrect and hacky. But I'm writing JavaScript, so it's
			// incorrect and hacky anyways.
			let location = window.location.href;
			let topLevel = location.substring(0, location.indexOf("/", 3));
			// Insert endpoint (hidden) into the form and add the submit button
			let input = document.getElementById("hidden-input");
			let button = document.getElementById("submit-button");

			input.setAttribute("value", JSON.stringify(results));
			button.removeAttribute("disabled");
			searchingForResults = false;
		}
	});
}

function numberOfRoundsUpdated() {
	let num = document.getElementById("rounds").value;
	if (!num) {
		return;
	}
	numDesiredResults = num;	
	if (num > results.length) {
		let button = document.getElementById("submit-button");
		button.setAttribute("disabled", "disabled");
		if (!searchingForResults) {
			queryPosition();
		}
	}  else if (num < results.length) {
		results = results.slice(num);
	}
}

function connectedOnlyUpdated() {
	let old = connectedOnly;
	connectedOnly = document.getElementById("connectedOnly").value.toLowerCase().includes("only");
	if (old !== connectedOnly) {
		let button = document.getElementById("submit-button");
		button.setAttribute("disabled", "disabled");
		results = [];
		markerGroup.clearLayers();
		queryPosition();
	}
}

function locStringUpdated() {
	let old = locString;
	locString = document.getElementById("locString").value;
	if (old !== locString) {
		map.setView([0, 0], 1);
		let button = document.getElementById("submit-button");
		button.setAttribute("disabled", "disabled");
		results = [];
		polygonGroup.clearLayers();
		markerGroup.clearLayers();
		placesPolygon = getPolygonFromLocString(locString);
	}
}

// settings may have been cached by the browser (wouldn't trigger the onchange),
// so check them once the DOM has loaded
window.addEventListener("DOMContentLoaded", (event) => {
	// TODO: stick map stuff in a function
	map = L.map("bounds-map", {center: [0, 0], zoom: 1});
	L.tileLayer("https://maps.wikimedia.org/osm-intl/{z}/{x}/{y}.png", {
		attribution: "&copy; <a href=\"https://www.openstreetmap.org/copyright\">OSM</a> contributors, <a href=\"https://foundation.wikimedia.org/wiki/Maps_Terms_of_Use\">Wikimedia Maps</a>"
	}).addTo(map);
	markerGroup = L.layerGroup().addTo(map);
	polygonGroup = L.layerGroup().addTo(map);
	locStringUpdated();
});
