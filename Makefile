APP_NAME=spiderlite
MODULE=github.com/emmelejail/$(APP_NAME)
VERSION=latest

.PHONY: all build run docker docker-run lint tidy clean

all: build

build:
	go build -o spiderlite ./cmd/crawler
	go build -o spiderlite-server ./cmd/server

run: build
	./spiderlite https://example.com

server: build
	./spiderlite-server

docker:
	docker build -t $(APP_NAME):$(VERSION) .

docker-run:
	docker run --rm $(APP_NAME):$(VERSION) https://example.com

lint:
	go vet ./...
	golangci-lint run || true

tidy:
	go mod tidy

clean:
	rm -f spiderlite spiderlite-server
