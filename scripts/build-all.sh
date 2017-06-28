#1/bin/bash
cd "$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
source common.sh

echo "Building"

source scripts/common-versioninfo.sh

echo "--> Creating build directory"
mkdir -vp build  2>&1 | indent

echo "--> Generating bindata"
go generate

echo "--> Building BookBrowser for Linux 64bit"
env GOOS=linux GOARCH=amd64 go build -ldflags "-X main.curversion=$APP_VERSION" -o build/BookBrowser-$APP_VERSION-linux-64bit  2>&1 | indent

echo "--> Building BookBrowser for Linux 32bit"
env GOOS=linux GOARCH=386 go build -ldflags "-X main.curversion=$APP_VERSION" -o build/BookBrowser-$APP_VERSION-linux-32bit  2>&1 | indent

echo "--> Building BookBrowser for Linux arm"
env GOOS=linux GOARCH=arm go build -ldflags "-X main.curversion=$APP_VERSION" -o build/BookBrowser-$APP_VERSION-linux-32bit  2>&1 | indent

echo "--> Building BookBrowser for Windows 64bit"
env GOOS=windows GOARCH=amd64 go build -ldflags "-X main.curversion=$APP_VERSION" -o build/BookBrowser-$APP_VERSION-windows-64bit.exe  2>&1 | indent

echo "--> Building BookBrowser for Windows 32bit"
env GOOS=windows GOARCH=386 go build -ldflags "-X main.curversion=$APP_VERSION" -o build/BookBrowser-$APP_VERSION-windows-32bit.exe  2>&1 | indent

echo "--> Building BookBrowser for Darwin 64bit"
env GOOS=darwin GOARCH=amd64 go build -ldflags "-X main.curversion=$APP_VERSION" -o build/BookBrowser-$APP_VERSION-darwin-64bit  2>&1 | indent

echo "--> Building BookBrowser for Darwin 32bit"
env GOOS=darwin GOARCH=386 go build -ldflags "-X main.curversion=$APP_VERSION" -o build/BookBrowser-$APP_VERSION-darwin-32bit  2>&1 | indent