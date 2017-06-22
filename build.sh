#!/bin/bash
cd "$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd scripts

export GITHUB_TOKEN

# To release a new version, create a new tag first, and push with git push --tags
# The release publishing requires a github token to be set in the var GITHUB_TOKEN
# Add any dependencies to install-dependencies.sh

./clean.sh || exit 1
./install-dependencies.sh || exit 1
./build-all.sh || exit 1
./release-notes.sh || exit 1
./publish-release.sh || exit 1
