lint:
	golangci-lint run ./... 

server-run:
	go run cmd/server/server.go

client-run:
	go run cmd/client/client.go

compose:
	docker-compose up

test:
	go test --cover --race ./...