build:
	go build
	cd frontend-svelte; npm install; npm run build

test:
	go fmt $(go list ./...)
	go vet $(go list ./...)
	go test ./...

