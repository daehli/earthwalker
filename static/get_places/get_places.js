let service = new google.maps.StreetViewService();

let results = [];

function queryPosition() {
	let randomLat = (Math.random() * 180.) - 90.;
	let randomLon = (Math.random() * 360.) - 180.;
	let point = new google.maps.LatLng(randomLat, randomLon);
	let radius = 10000;
	service.getPanoramaByLocation(point, radius, function(result, status) {
		if (status == google.maps.StreetViewStatus.OK) {
			let nearestLatLng = result.location.latLng;
			console.log(nearestLatLng.toString());
			results.push(nearestLatLng);
		} else {
			console.log("Failed to get location.");
			console.log(status.toString());
		}
		if (results.length < 1) {
			document.getElementById("counter").innerHTML = results.length;
			queryPosition();
		} else {
			// Yea, this is probably incorrect and hacky. But I'm writing JavaScript, so it's
			// incorrect and hacky anyways.
			let location = window.location.href;
			let topLevel = location.substring(0, location.indexOf("/", 3));
			// Insert a form containing results and send it off to the endpoint
			document.body.innerHTML = "<form id='resultForm' action='" + topLevel + "/found_points' method='post'>" +
				"<input type='hidden' name='result' value='" + JSON.stringify(results) + "'/></form>"
			document.getElementById("resultForm").submit();
		}
	});
}

queryPosition();
