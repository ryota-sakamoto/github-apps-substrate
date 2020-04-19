.PHONY: erorr dev build

error:
	exit 1

dev:
	go run cmd/api/api.go

build:
	GOOS=linux GOARCH=amd64 go build -o api cmd/api/api.go

rename:
	make error
