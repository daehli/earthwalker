// Be warned, traveller. You are entering the domain of some very dodgy javascript
// hacks. Maybe that is what you like. If so, please look around.

function injectStylesheet() {
	var node = document.createElement("link");
	node.href = "/static/modify_frontend/modify.css";
	node.rel = "stylesheet";
	document.body.appendChild(node);

	// This MutationObserver always resets the title to earthwalker.
	let interval = setInterval(function() {
		try {
			new MutationObserver(function(mutations) {
				if (document.title != "earthwalker") {
					document.title = "earthwalker";
				}
			}).observe(
				document.querySelector('title'),
				{ childList: true }
			);
			clearInterval(interval);
		} catch (e) {
			// Title element is not there yet.
			// Wait for the next poll...
		}
	}, 50);

	createMinimap();
}

window.onload = injectStylesheet;
// Sometimes, the google scripts crash on startup. Just reload the page if that happens.
window.onerror = function(e) {
	if (e.includes("Timer")) {
		location.reload(false);
	}
};

let replaceStateLocal = history.replaceState;
history.replaceState = function() {
}

let pushStateLocal = history.pushState;
history.pushState = function() {
}

let leafletMap = null;

// The leaflet minimap!
function createMinimap() {
	let floatingContainer = document.createElement("div");
	floatingContainer.id = "leaflet-container";
	document.body.appendChild(floatingContainer);

	let barDiv = document.createElement("div");
	barDiv.id = "navigation-bar";
	barDiv.className = "btn-group btn-group-sm";
	floatingContainer.appendChild(barDiv)

	let zoomInButton = document.createElement("button");
	zoomInButton.type = "button";
	zoomInButton.className = "btn btn-light";
	zoomInButton.innerHTML = "⬉";
	barDiv.appendChild(zoomInButton);
	zoomInButton.addEventListener("click", function() {
		scaleMap(true);
	});

	let zoomOutButton = document.createElement("button");
	zoomOutButton.type = "button";
	zoomOutButton.className = "btn btn-light";
	zoomOutButton.innerHTML = "⬊";
	barDiv.appendChild(zoomOutButton);
	zoomOutButton.addEventListener("click", function() {
		scaleMap(false);
	});

	let marker = null;
	let guessButton = document.createElement("button");
	guessButton.type = "button";
	guessButton.className = "btn btn-primary btn-sm float-right disabled";
	guessButton.innerHTML = "Guess!";
	floatingContainer.appendChild(guessButton);
	guessButton.addEventListener("click", function() {
		if (marker == null) {
			alert("You have to add a marker first! Do this by clicking the map.");
			return;
		}
		// Post data back to earthwalker.
		let location = window.location.href;
		let topLevel = location.substring(0, location.indexOf("/", 3));
		let xhr = new XMLHttpRequest();
		xhr.open("POST", topLevel + "/guess", true);
		xhr.setRequestHeader('Content-Type', 'application/json');
		xhr.send(JSON.stringify(marker.getLatLng()));
		window.location.replace(topLevel + "/scores");
	});

	let leafletMapDiv = document.createElement("div");
	leafletMapDiv.id = "leaflet-map";
	floatingContainer.appendChild(leafletMapDiv);

	leafletMap = L.map("leaflet-map").setView([0.0, 0.0], 1);

	L.tileLayer("https://maps.wikimedia.org/osm-intl/{z}/{x}/{y}.png", {
		attribution: "&copy; <a href=\"https://www.openstreetmap.org/copyright\">OSM</a> contributors, <a href=\"https://foundation.wikimedia.org/wiki/Maps_Terms_of_Use\">Wikimedia Maps</a>"
	}).addTo(leafletMap);

	function onMapClick(event) {
		if (marker != null) {
			leafletMap.removeControl(marker);
		}
		marker = L.marker(event.latlng);
		marker.addTo(leafletMap);
		guessButton.className = guessButton.className.replace("disabled", "");
	}

	leafletMap.on("click", onMapClick);

	setTimeout(function() {
		leafletMap.invalidateSize();
	}, 100);
}

let sizes = [
	[150, 150],
	[300, 300],
	[500, 500],
	[800, 800],
];

function scaleMap(bigger) {
	let map = document.getElementById("leaflet-container");

	let size = [map.scrollWidth, map.scrollHeight];
	let nextSize = null;

	let index = -1;
	for (el in sizes) {
		index++;
		if (sizes[el][0] == size[0]) {
			break;
		}
	}

	if (bigger) {
		index++;
	} else {
		index--;
	}

	if (index < 0) {
		index = 0;
	}
	if (index > sizes.length) {
		index = sizes.length;
	}

	map.style.width = sizes[index][0] + "px";
	map.style.height = sizes[index][1] + "px";

	leafletMap.invalidateSize();
}
