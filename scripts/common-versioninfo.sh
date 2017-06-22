#1/bin/bash
cd "$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
source common.sh

echo "--> Getting version info"
# Gets current tag, or else previous-tag+commit-hash~dev
export APP_VERSION="$(git name-rev --name-only --tags HEAD | sed "s/undefined/$(git describe --abbrev=0 --tags)+$(git rev-parse --short HEAD)-dev/")"
echo "APP_VERSION=$APP_VERSION" | indent
export IS_DEV=$(if [[ $(git name-rev --name-only --tags HEAD) == "undefine
d" ]];then echo true;else echo false;fi)
echo "IS_DEV=$IS_DEV" | indent