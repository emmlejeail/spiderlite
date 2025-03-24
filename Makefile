APP_NAME=spiderlite
MODULE=github.com/emmelejail/$(APP_NAME)
VERSION=latest

.PHONY: all build run docker docker-run lint tidy clean

all: build

build:
	go build -o $(APP_NAME) ./main.go

run:
	./$(APP_NAME) https://example.com

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
	rm -f $(APP_NAME)
