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
CHANGELOG_FILE="CHANGELOG.md"

echo "Releasing $VERSION_TAG..."

# Get the latest tag to determine commit range
LATEST_TAG=$(git describe --tags --abbrev=0 2>/dev/null || echo "")

# Parse commits since last tag
if [[ -n "$LATEST_TAG" ]]; then
  COMMIT_RANGE="$LATEST_TAG..HEAD"
else
  COMMIT_RANGE="HEAD"
fi

# Get commits and categorize them
echo "Parsing commits from ${COMMIT_RANGE}..."
echo ""

ADDED=()
FIXED=()
CHANGED=()
OTHER=()

while IFS= read -r commit; do
  # Try conventional commit format first: type: description or type(scope): description
  if [[ "$commit" =~ ^([a-z]+)(\(.+\)):\s+(.+)$ ]]; then
    TYPE="${BASH_REMATCH[1]}"
    MESSAGE="${BASH_REMATCH[3]}"

    case "$TYPE" in
      feat|add)
        ADDED+=("$MESSAGE")
        ;;
      fix|bugfix)
        FIXED+=("$MESSAGE")
        ;;
      chore|refactor|perf|style|test|ci|build|docs)
        CHANGED+=("$MESSAGE")
        ;;
      *)
        OTHER+=("$MESSAGE")
        ;;
    esac
  else
    # Non-conventional commit - add as-is to Other
    OTHER+=("$commit")
  fi
done < <(git log "${COMMIT_RANGE}" --pretty=format:"%s" --reverse)

# Track all commits for the "All Changes" section
ALL_COMMITS=()

# Build changelog entries
CHANGELOG_ENTRIES=()

if [[ ${#ADDED[@]} -gt 0 ]]; then
  CHANGELOG_ENTRIES+=("### Added")
  for entry in "${ADDED[@]}"; do
    CHANGELOG_ENTRIES+=("- $entry")
    ALL_COMMITS+=("✨ $entry")
  done
  CHANGELOG_ENTRIES+=("")
fi

if [[ ${#FIXED[@]} -gt 0 ]]; then
  CHANGELOG_ENTRIES+=("### Fixed")
  for entry in "${FIXED[@]}"; do
    CHANGELOG_ENTRIES+=("- $entry")
    ALL_COMMITS+=("🐛 $entry")
  done
  CHANGELOG_ENTRIES+=("")
fi

if [[ ${#CHANGED[@]} -gt 0 ]]; then
  CHANGELOG_ENTRIES+=("### Changed")
  for entry in "${CHANGED[@]}"; do
    CHANGELOG_ENTRIES+=("- $entry")
    ALL_COMMITS+=("♻️ $entry")
  done
  CHANGELOG_ENTRIES+=("")
fi

if [[ ${#OTHER[@]} -gt 0 ]]; then
  CHANGELOG_ENTRIES+=("### Other")
  for entry in "${OTHER[@]}"; do
    CHANGELOG_ENTRIES+=("- $entry")
    ALL_COMMITS+=("📝 $entry")
  done
  CHANGELOG_ENTRIES+=("")
fi

# Add expandable "All Changes" section at the end
if [[ ${#ALL_COMMITS[@]} -gt 0 ]]; then
  CHANGELOG_ENTRIES+=("<details>")
  CHANGELOG_ENTRIES+=("<summary>All Changes</summary>")
  CHANGELOG_ENTRIES+=("")
  for entry in "${ALL_COMMITS[@]}"; do
    CHANGELOG_ENTRIES+=("- $entry")
  done
  CHANGELOG_ENTRIES+=("")
  CHANGELOG_ENTRIES+=("</details>")
  CHANGELOG_ENTRIES+=("")
fi

# Get current date
DATE=$(date -u +"%Y-%m-%d")

# Create/update CHANGELOG.md
if [[ ! -f "$CHANGELOG_FILE" ]]; then
  # Create new CHANGELOG.md with header
  cat > "$CHANGELOG_FILE" << EOF
# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

EOF
fi

# Prepend new release to CHANGELOG.md
TEMP_FILE=$(mktemp)
{
  echo "## [$VERSION_TAG] - $DATE"
  echo ""
  for entry in "${CHANGELOG_ENTRIES[@]}"; do
    echo "$entry"
  done
  echo ""
  cat "$CHANGELOG_FILE"
} > "$TEMP_FILE"
mv "$TEMP_FILE" "$CHANGELOG_FILE"

echo "Changelog updated:"
echo "## [$VERSION_TAG] - $DATE"
for entry in "${CHANGELOG_ENTRIES[@]}"; do
  echo "$entry"
done
echo ""

# Update version in version.go
sed -i.bak "s/var Version = \".*\"/var Version = \"$VERSION_TAG\"/" version.go
rm -f version.go.bak

# Update README install command with new version
sed -i.bak "s/go install github.com\/mohammed-io\/config-sync@v[0-9.]*/go install github.com\/mohammed-io\/config-sync@$VERSION_TAG/" README.md
rm -f README.md.bak

# Commit changes
git add version.go README.md "$CHANGELOG_FILE"
git commit -m "chore: release $VERSION_TAG"

# Create and push tag
git tag "$VERSION_TAG"
git push origin main --tags
git push origin "$VERSION_TAG"

echo "✓ Released $VERSION_TAG"
