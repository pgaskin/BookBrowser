#1/bin/bash
cd "$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
source common.sh

export GITHUB_TOKEN

echo
echo "Publishing GitHub release"
source scripts/common-versioninfo.sh

if [[ "$IS_DEV" == "false" ]]; then
echo "--> Creating release from current tag"
github-release release \
    --user geek1011 \
    --repo BookBrowser \
    --tag $APP_VERSION \
    --name "BookBrowser $APP_VERSION" \
    --description "$(cat build/release-notes.md)" | indent

echo "--> Uploading files"
for f in build/BookBrowser-*;do 
    fn="$(basename $f)"
    echo "$f > $fn" | indent
    github-release upload \
        --user geek1011 \
        --repo BookBrowser \
        --tag $APP_VERSION \
        --name "$fn" \
        --file "$f" | indent | indent
done
else
echo "--> Skipping because not on a tag (probably development version)"
fi