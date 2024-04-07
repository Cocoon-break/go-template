.PHONY: build

NAME=go-template
VERSION=0.0.1
RELEASE_NUM=1
# DIR="$(shell cd "$( dirname "${BASH_SOURCE[0]}")" && pwd)"
DATE:="$(shell date +'%Y-%m-%dT%H:%M:%S')"
GITHASH="$(shell git rev-parse HEAD)"

lint:
	golangci-lint run --fix
unittest: 
	go test -count=1 ./... -v -buildvcs=false -short
build: 
	GOARCH=amd64 GOOS=linux CGO_ENABLED=0 go build -buildvcs=false -ldflags "-extldflags '-static' \
	-X '$(NAME)/pkg/env.Built=$(DATE)' \
	-X '$(NAME)/pkg/env.Version=$(VERSION)-$(RELEASE_NUM)' \
	-X '$(NAME)/pkg/env.App=$(NAME)' \
	-X '$(NAME)/pkg/env.GitHash=$(GITHASH)'" \
	-tags=jsoniter -mod=vendor \
	-o cmd/$(NAME) cmd/main.go 

all: lint unittest