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

VERSION="$1"

# Reject if user included 'v' prefix
if [[ "$VERSION" == v* ]]; then
  echo "Error: Version should not include 'v' prefix"
  echo "Use: $0 0.0.1"
  echo "Not: $0 v0.0.1"
  exit 1
fi

# Validate semver format (X.Y.Z where X,Y,Z are numbers)
if ! [[ "$VERSION" =~ ^[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
  echo "Error: Invalid version format '$VERSION'"
  echo "Version must be in semantic versioning format: X.Y.Z"
  echo "Example: 0.0.1, 1.2.3, 2.0.0"
  exit 1
fi

VERSION_TAG="v$VERSION"

echo "Releasing $VERSION_TAG..."

# Update version in version.go
sed -i.bak "s/var Version = \".*\"/var Version = \"$VERSION_TAG\"/" version.go
rm -f version.go.bak

# Update README install command with new version
sed -i.bak "s/go install github.com\/mohammed-io\/config-sync@v[0-9.]*/go install github.com\/mohammed-io\/config-sync@$VERSION_TAG/" README.md
rm -f README.md.bak

# Commit changes
git add version.go README.md
git commit -m "chore: release $VERSION_TAG"

# Create and push tag
git tag "$VERSION_TAG"
git push origin main
git push origin "$VERSION_TAG"

echo "✓ Released $VERSION_TAG"
