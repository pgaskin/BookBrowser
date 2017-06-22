#!/bin/bash
cd "$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd scripts

./clean.sh || exit 1
./install-dependencies.sh || exit 1
./build-all.sh || exit 1
./release-notes.sh || exit 1
./publish-release.sh || exit 1
