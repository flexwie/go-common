deps:
	go mod download

test: deps
	go test ./...

proto:
	protoc --go_out=. --go-grpc_out=. --go_opt=paths=source_relative --go-grpc_opt=paths=source_relative auth/*.proto
