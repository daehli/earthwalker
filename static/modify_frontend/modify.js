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

