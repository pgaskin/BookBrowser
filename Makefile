.PHONY: default
default: clean build-deps deps test build

.PHONY: clean
clean:
	rm -rf build BookBrowser

.PHONY: build-deps
build-deps:
	GO111MODULE=off go get -v "github.com/aktau/github-release"
	GO111MODULE=off go get -v "github.com/gobuffalo/packr/..."

.PHONY: deps
deps:
	GO111MODULE=on go mod download

.PHONY: generate
generate:
	GO111MODULE=on go generate ./...

.PHONY: test
test:
	GO111MODULE=on go test ./...

.PHONY: build
build:
	mkdir -p build
	GO111MODULE=on go build -ldflags "-X main.curversion=dev" -o "build/BookBrowser"

.PHONY: install
install:
	GO111MODULE=on go install