#!/bin/bash
cd "$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd scripts

export GITHUB_TOKEN

./clean.sh || exit 1
./install-dependencies.sh || exit 1
./build-all.sh || exit 1
