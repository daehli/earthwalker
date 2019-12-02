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

	// The leaflet minimap!
	let leafletMap = document.createElement("div");
	leafletMap.id = "leaflet-map";
	leafletMap.class = "leaflet-map";
	document.body.appendChild(leafletMap);

	let map = L.map("leaflet-map").setView([0.0, 0.0], 1);

	L.tileLayer("https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png", {
		attribution: "&copy; <a href=\"https://www.openstreetmap.org/copyright\">OpenStreetMap</a> contributors"
	}).addTo(map);

	let marker = null;
	function onMapClick(event) {
		if (marker != null) {
			map.removeControl(marker);
		}
		marker = L.marker(event.latlng);
		marker.addTo(map);
	}

	map.on("click", onMapClick);

	setTimeout(function() {
		leafletMap.style = "";
		map.invalidateSize();
	}, 100);
}

window.onload = injectStylesheet;
// Sometimes, the google scripts crash on startup. Just reload the page if that happens.
window.onerror = function() {
	location.reload(false);
};

let replaceStateLocal = history.replaceState;
history.replaceState = function() {
}

let pushStateLocal = history.pushState;
history.pushState = function() {
}
