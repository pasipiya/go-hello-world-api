.PHONY: run build test cover fmt vet

run:
	go run ./cmd/server

build:
	go build -o bin/server ./cmd/server

test:
	go test -race ./...

cover:
	go test -cover ./...

fmt:
	gofmt -l -w .

vet:
	go vet ./...
