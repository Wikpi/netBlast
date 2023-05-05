bld:
	go build -o build/server ./internal/server/server.go
	go build -o build/client ./internal/client/client.go

server:
	go run ./internal/server/server.go
client:
	go run ./internal/client/client.go	

test:
	go run ./test/unitTests