#1/bin/bash
cd "$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
source common.sh

echo "Generating release notes"

source scripts/common-versioninfo.sh

echo "--> Changelog"
echo "## Changes for $APP_VERSION" | tee -a build/release-notes.md | indent
echo "$(git log $(git describe --tags --abbrev=0 HEAD^)..HEAD --oneline)" | tee -a build/release-notes.md | indent
echo "" | tee -a build/release-notes.md | indent
echo "--> Usage"
echo "## Usage" | tee -a build/release-notes.md | indent
echo "1. Download the binary for your platform below" | tee -a build/release-notes.md | indent
echo "2. Copy it to the directory with your books" | tee -a build/release-notes.md | indent
echo "3. Run it" | tee -a build/release-notes.md | indent