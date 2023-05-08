.PHONY: build dev generate

build: generate
	mkdir -p build
	go build -o build/onlyserve ./cmd/onlyserve

dev:
	saq -- bash -c 'make && ./build/onlyserve --http localhost:8081'

generate:
	go generate ./...
