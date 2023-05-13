all:
	go build -o build/server ./cmd/server/server.go
	go build -o build/client ./cmd/client/client.go

server:
	go run ./cmd/server/server.go
client:
	go run ./cmd/client/client.go	

test:
	go run ./test/unitTests