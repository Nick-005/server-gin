all:
	swag init -g cmd/server/main.go --output docs --parseDependency
	go build ./cmd/server/main.go
	./main