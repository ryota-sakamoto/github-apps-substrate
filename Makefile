.PHONY: erorr dev build

error:
	exit 1

dev:
	go run cmd/api/api.go

build:
	GOOS=linux GOARCH=amd64 go build -o api cmd/api/api.go

rename:
	$(eval BEFORE := $(shell head -n 1 go.mod | sed -e "s/module //g"))
	git ls-files | grep -E "(.go|go.mod)" | xargs sed -i "" -e "s#$(BEFORE)#$(NAME)#g"
