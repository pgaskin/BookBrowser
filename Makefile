ifndef $$GOPATH
  export GOPATH=${HOME}/go
endif

.PHONY: default
default: clean build-deps deps generate test build

.PHONY: clean
clean:
	rm -rf build BookBrowser

.PHONY: build-deps
build-deps:
	go get -v "github.com/kardianos/govendor"
	go get -v "github.com/aktau/github-release"
	go get -v github.com/jteeuwen/go-bindata/...
	go get -v github.com/elazarl/go-bindata-assetfs/...

.PHONY: deps
deps:
	govendor sync

.PHONY: generate
generate:
	go generate ./...

.PHONY: test
test:
	go test ./...

.PHONY: build
build:
	mkdir -p build
	go build -ldflags "-X main.curversion=dev" -o "build/BookBrowser"

.PHONY: install
install:
	go install
