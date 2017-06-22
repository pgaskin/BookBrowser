#!/bin/bash
cd "$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

set -e

if [ -z "$1" ]; then
    echo "Usage: $0 v#.#.#"
    exit 1
fi

read -r -p "Create a new tag for version $1, push it, and build a release? [y/N] " response
case "$response" in
    [yY][eE][sS]|[yY]) 
        echo "Creating release"
        ;;
    *)
        echo Aborted
        exit 1
        ;;
esac

git tag $1
git push --tags
./build.sh
./scripts/release-notes.sh
./scripts/publish-release.sh
