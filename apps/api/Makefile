run:
	go run main.go

build:
	go build -o bin/main main.go

start:
	./bin/main

generate:
	oapi-codegen -config .oapi-codegen.yaml ../../packages/openapi/schema.yaml