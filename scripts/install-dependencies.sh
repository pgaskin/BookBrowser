#1/bin/bash
cd "$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
source common.sh

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

echo "--> github-release"
go get -u -v github.com/aktau/github-release 2>&1 | indent

echo