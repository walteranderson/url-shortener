BINARY_NAME=bin/url-shortener-server

build:
	go build -o $(BINARY_NAME) -v ./cmd/http
run: build
	./$(BINARY_NAME)
clean:
	rm -f $(BINARY_NAME)
up:
	docker-compose up -d
down:
	docker-compose down
