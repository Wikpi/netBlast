all:
	mkdir ./buildr
	go build [-o ./build] [-arg server] ./cmd/main.go
	go build [-o ./build] [-arg client] ./cmd/main.go

server:
	go run ./cmd/main.go -arg server
client:
	go run ./cmd/main.go -arg client	

test:
	go -test ./cmd/client/client_test.go
	go -test ./cmd/client/server_test.go