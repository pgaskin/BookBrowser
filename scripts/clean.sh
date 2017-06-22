#1/bin/bash
cd "$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
source common.sh

echo "Cleaning up old build"

echo "--> Removing build directory"
rm -rfv build  2>&1 | indent
echo "--> Removing bindata"
rm -rfv bindata_assetfs.go 2>&1 | indent

echo