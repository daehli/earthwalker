.PHONY: build

build:
	go build

test:
	go fmt $(go list ./...)
	go vet $(go list ./...)
	go test ./...
