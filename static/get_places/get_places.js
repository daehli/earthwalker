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

function showPolygonOnMap(map, polygon) {
	L.tileLayer("https://maps.wikimedia.org/osm-intl/{z}/{x}/{y}.png", {
		attribution: "&copy; <a href=\"https://www.openstreetmap.org/copyright\">OSM</a> contributors, <a href=\"https://foundation.wikimedia.org/wiki/Maps_Terms_of_Use\">Wikimedia Maps</a>"
	}).addTo(map);
	let map_poly = L.geoJSON(polygon).addTo(map);
	map.fitBounds(map_poly.getBounds());
}

function getPolygonFromLocString(loc) {
	const Http = new XMLHttpRequest();
	const url = "https://nominatim.openstreetmap.org/search?q=" + encodeURI(loc.replace(" ", "+")) + "&polygon_geojson=1&limit=1&format=json";
	Http.open("GET", url);
	Http.send();

	Http.onreadystatechange = (event) => {
		if (Http.readyState == 4) {
			response = JSON.parse(Http.responseText)[0];
			console.log("Response received, display name: " + response["display_name"]);
			if (response["geojson"]["type"].toLowerCase() === "multipolygon") {
				polygon = turf.multiPolygon(response["geojson"]["coordinates"])
			} else {
				polygon = turf.polygon(response["geojson"]["coordinates"]);
			}
			showPolygonOnMap(map, polygon);
			return polygon;
		}
	}
}

function getRandomLngLat(polygon) {
	bounds = turf.bbox(polygon);
	do {
		randomLng = (Math.random() * (bounds[2] - bounds[0]) + bounds[0]);
		randomLat = (Math.random() * (bounds[3] - bounds[1]) + bounds[1]);
		lnglat = turf.point([randomLng, randomLat]);
	} while (!turf.booleanPointInPolygon(lnglat, polygon))
	L.marker([randomLat, randomLng]).addTo(map);
	return lnglat;
}

function queryPosition() {
	searchingForResults = true;
	let randomLat = (Math.random() * 180.) - 90.;
	let randomLon = (Math.random() * 360.) - 180.;
	let point = new google.maps.LatLng(randomLat, randomLon);
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
		queryPosition();
	}
}

// settings may have been cached by the browser (wouldn't trigger the onchange),
// so check them once the DOM has loaded
window.addEventListener("DOMContentLoaded", (event) => {
	map = L.map("bounds-map");
	getPolygonFromLocString("powell wyoming");
	numberOfRoundsUpdated();
	connectedOnlyUpdated();
	queryPosition();
});
