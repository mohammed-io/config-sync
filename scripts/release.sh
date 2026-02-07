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

# Get repo URL for links
REPO_URL=$(git config --get remote.origin.url)
if [[ "$REPO_URL" =~ git@github.com:(.+)\.git$ ]]; then
  REPO_URL="https://github.com/${BASH_REMATCH[1]}"
fi
COMMIT_RANGE_LINK="${REPO_URL}/compare/${LATEST_TAG}...${VERSION_TAG}"

# Get commits and categorize them
echo "Parsing commits from ${COMMIT_RANGE}..."
echo ""

ADDED=()
FIXED=()
CHANGED=()
OTHER=()

# Get commits into a temp file
TEMP_COMMITS=$(mktemp)
git log "${COMMIT_RANGE}" --pretty=format:"%s" > "$TEMP_COMMITS"

# Get commit hashes into another temp file
TEMP_HASHES=$(mktemp)
git log "${COMMIT_RANGE}" --pretty=format:"%s|%H" > "$TEMP_HASHES"

while IFS= read -r commit || [[ -n "$commit" ]]; do
  # Check if commit has a colon (conventional commit format: type: message)
  if [[ "$commit" == *": "* ]]; then
    # Extract type (before first colon) and message (after colon and space)
    TYPE="${commit%%:*}"
    MESSAGE="${commit#*: }"
    # Remove optional (scope) from type if present
    TYPE="${TYPE%%\(*}"

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
done < "$TEMP_COMMITS"

rm -f "$TEMP_COMMITS"

# Helper function to get hash for a message (must be called before TEMP_HASHES is deleted)
# Also stores TEMP_HASHES path globally so subshells can find it
_get_hash_file="$TEMP_HASHES"

get_hash() {
  local msg="$1"
  # Use grep with fixed string matching for the message
  if [[ -f "$_get_hash_file" ]]; then
    local hash=$(grep -F "$msg" "$_get_hash_file" | cut -d'|' -f2 | head -1)
    echo "$hash"
  else
    echo ""
  fi
}

# Build changelog entries
CHANGELOG_ENTRIES=()

if [[ ${#ADDED[@]} -gt 0 ]]; then
  CHANGELOG_ENTRIES+=("### Added")
  for entry in "${ADDED[@]}"; do
    HASH=$(get_hash "$entry")
    SHORT_HASH="${HASH:0:7}"
    CHANGELOG_ENTRIES+=("- $entry ([${SHORT_HASH}](${REPO_URL}/commit/${HASH}))")
  done
  CHANGELOG_ENTRIES+=("")
fi

if [[ ${#FIXED[@]} -gt 0 ]]; then
  CHANGELOG_ENTRIES+=("### Fixed")
  for entry in "${FIXED[@]}"; do
    HASH=$(get_hash "$entry")
    SHORT_HASH="${HASH:0:7}"
    CHANGELOG_ENTRIES+=("- $entry ([${SHORT_HASH}](${REPO_URL}/commit/${HASH}))")
  done
  CHANGELOG_ENTRIES+=("")
fi

if [[ ${#CHANGED[@]} -gt 0 ]]; then
  CHANGELOG_ENTRIES+=("### Changed")
  for entry in "${CHANGED[@]}"; do
    HASH=$(get_hash "$entry")
    SHORT_HASH="${HASH:0:7}"
    CHANGELOG_ENTRIES+=("- $entry ([${SHORT_HASH}](${REPO_URL}/commit/${HASH}))")
  done
  CHANGELOG_ENTRIES+=("")
fi

# "Other" section as expandable details
if [[ ${#OTHER[@]} -gt 0 ]]; then
  CHANGELOG_ENTRIES+=("<details>")
  CHANGELOG_ENTRIES+=("<summary>Other</summary>")
  CHANGELOG_ENTRIES+=("")
  for entry in "${OTHER[@]}"; do
    HASH=$(get_hash "$entry")
    SHORT_HASH="${HASH:0:7}"
    CHANGELOG_ENTRIES+=("- $entry ([${SHORT_HASH}](${REPO_URL}/commit/${HASH}))")
  done
  CHANGELOG_ENTRIES+=("")
  CHANGELOG_ENTRIES+=("</details>")
  CHANGELOG_ENTRIES+=("")
fi

# Clean up temp hash file (must be after all get_hash calls)
rm -f "$TEMP_HASHES"

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
  if [[ -n "$COMMIT_RANGE_LINK" ]]; then
    echo "[Full Changelog](${COMMIT_RANGE_LINK})"
    echo ""
  fi
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
