build:
	@go build -o bin/data cmd/data/main.go

data: build
	@./bin/data