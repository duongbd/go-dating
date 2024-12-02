.PHONY: build run test docker-build docker-run

build:
	go build -o datingapp ./cmd/app/main.go 

run: build
	./datingapp 

test:
	go test -race ./...

docker-compose-up:
	docker-compose up -d 
	
docker-compose-down:
	docker-compose down