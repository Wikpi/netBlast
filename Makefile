all:
	mkdir ./build
	go build -o ./build ./cmd/main.go

server:
	go run ./cmd/main.go server
client:
	go run ./cmd/main.go client	

test:
	go -test ./cmd/client/client_test.go
	go -test ./cmd/client/server_test.go