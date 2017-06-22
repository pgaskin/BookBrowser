#!/bin/bash
set -e

indent() {
  while read l; do echo "    $l"; done
}

echo "Cleaning up old build"

echo "--> Removing build directory"
rm -rfv build  2>&1 | indent

echo
echo "Installing dependencies"

echo "--> go-bindata"
go get -u -v github.com/jteeuwen/go-bindata/...  2>&1 | indent

echo "--> go-bindata-assetfs"
go get -u -v github.com/elazarl/go-bindata-assetfs/...  2>&1 | indent

echo "--> etree"
go get -u -v github.com/beevik/etree  2>&1 | indent

echo "--> go-zglob"
go get -u -v github.com/mattn/go-zglob  2>&1 | indent

echo "--> resize"
go get -u -v github.com/nfnt/resize  2>&1 | indent

echo "--> zipfs"
go get -u -v golang.org/x/tools/godoc/vfs/zipfs 2>&1 | indent

echo
echo "Building"

echo "--> Creating build directory"
mkdir -vp build  2>&1 | indent

echo "--> Generating bindata"
go generate

echo "--> Building BookBrowser for Linux 64bit"
env GOOS=linux GOARCH=amd64 go build -o build/BookBrowser-linux-64bit  2>&1 | indent

if [[ "$1" == "all" ]]; then
    echo "--> Building BookBrowser for Linux 32bit"
    env GOOS=linux GOARCH=386 go build -o build/BookBrowser-linux-32bit  2>&1 | indent

    echo "--> Building BookBrowser for Windows 64bit"
    env GOOS=windows GOARCH=amd64 go build -o build/BookBrowser-windows-64bit.exe  2>&1 | indent

    echo "--> Building BookBrowser for Windows 32bit"
    env GOOS=windows GOARCH=386 go build -o build/BookBrowser-windows-32bit.exe  2>&1 | indent

    echo "--> Building BookBrowser for Darwin 64bit"
    env GOOS=darwin GOARCH=amd64 go build -o build/BookBrowser-darwin-64bit  2>&1 | indent

    echo "--> Building BookBrowser for Darwin 32bit"
    env GOOS=darwin GOARCH=386 go build -o build/BookBrowser-darwin-32bit  2>&1 | indent
fi


echo
echo "Generating release notes"

echo "--> Changelog"
echo "## Changes for $(git describe --tags --abbrev=0 HEAD^)" | tee -a build/release-notes.md | indent
echo "$(git log $(git describe --tags --abbrev=0 HEAD^)..HEAD --oneline)" | tee -a build/release-notes.md | indent
echo "" | tee -a build/release-notes.md | indent
echo "--> Usage"
echo "## Usage" | tee -a build/release-notes.md | indent
echo "1. Download the binary for your platform below" | tee -a build/release-notes.md | indent
echo "2. Copy it to the directory with your books" | tee -a build/release-notes.md | indent
echo "3. Run it" | tee -a build/release-notes.md | indent


echo

echo "Cleaning up"
echo "--> Removing bindata"
rm -rfv bindata_assetfs.go 2>&1 | indent

echo
echo "Built to $PWD/build/BookBrowser"
