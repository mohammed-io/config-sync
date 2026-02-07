#!/bin/bash
set -e

# Release script for config-sync
# Usage: ./scripts/release.sh <version>
# Example: ./scripts/release.sh 0.0.2

if [ -z "$1" ]; then
  echo "Usage: $0 <version>"
  echo "Example: $0 0.0.2"
  exit 1
fi

VERSION="v$1"
VERSION_NO_V="$1"

echo "Releasing $VERSION..."

# Update version in version.go
sed -i.bak "s/var Version = \".*\"/var Version = \"$VERSION\"/" version.go
rm -f version.go.bak

# Update README install command with new version
sed -i.bak "s/go install github.com\/mohammed-io\/config-sync@v[0-9.]*/go install github.com\/mohammed-io\/config-sync@$VERSION/" README.md
rm -f README.md.bak

# Commit changes
git add version.go README.md
git commit -m "chore: release $VERSION"

# Create and push tag
git tag "$VERSION"
git push origin main
git push origin "$VERSION"

echo "✓ Released $VERSION"
