build:
	@go build -o bin/microCmp

run : build
	@./bin/microCmp

test:
	@go test -v ./...