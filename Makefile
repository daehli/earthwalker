build:
	go build
<<<<<<< HEAD
=======
	cd frontend-svelte; npm install; npm run build
>>>>>>> features/domainrefactor
test:
	go fmt $(go list ./...)
	go vet $(go list ./...)
	go test ./...

