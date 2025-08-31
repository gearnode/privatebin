#!/bin/bash

# Extract changelog for a specific version from CHANGELOG.md
# Usage: ./scripts/extract-changelog.sh [version]

VERSION="${1:-${GORELEASER_CURRENT_TAG}}"
VERSION="${VERSION#v}" # Remove 'v' prefix if present

# Extract the section for the specific version
awk -v version="$VERSION" '
    /^## \[/ {
        if (found) exit
        if ($2 == "["version"]") {
            found = 1
            next
        }
    }
    found && /^## \[/ { exit }
    found { print }
' CHANGELOG.md | sed '/^$/d' # Remove empty lines
