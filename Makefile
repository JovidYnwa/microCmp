build:
	@go build - o bin/gomicro

run : build
	@./bin/gomicro

test:
	@go test -v ./...