build:
	@go build -o bin/product-service

run: build
	@./bin/product-service

test:
	 @go test -v ./...