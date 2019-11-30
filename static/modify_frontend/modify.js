function injectStylesheet() {
	var node = document.createElement("link");
	node.href = "/static/modify_frontend/modify.css";
	node.rel = "stylesheet";
	document.body.appendChild(node);
}

window.onload = injectStylesheet;

var replaceStateLocal = history.replaceState;
history.replaceState = function() {
	document.title = "earthwalker";
}

var pushStateLocal = history.pushState;
history.pushState = function() {
	document.title = "earthwalker";
}


