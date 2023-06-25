build:
	@go build -o bin/smartpay-be

run: build
	@./bin/smartpay-be

test:
	@go test -v ./...