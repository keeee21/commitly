run:
	go run main.go

build:
	go build -o bin/main main.go

start:
	./bin/main

swagger:
	swag init -g main.go -o ./swagger --parseInternal --parseDependency
	rm -f ./swagger/docs.go

swagger-gen: swagger
	cd ../../packages/openapi && npx swagger2openapi ../../apps/api/swagger/swagger.json -o schema.yaml --yaml
