VERSION=0.0.3
LDFLAGS=-ldflags "-X main.Version=${VERSION}"
GO111MODULE=on

all: check

check: consulclient.go consulclient_test.go
	go test -v ./...

deps:
	go get -d
	go mod tidy

deps-update:
	go get -u -d
	go mod tidy

tag:
	git tag v${VERSION}
	git push origin v${VERSION}
	git push origin master
