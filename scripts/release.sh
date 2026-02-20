#!/usr/bin/env bash

# Release helper script for versioning
# Usage: ./scripts/release.sh [major|minor|patch]

set -e

RELEASE_TYPE="${1:-patch}"

# Verify in git repo
if ! git rev-parse --git-dir >/dev/null 2>&1; then
    echo "✗ Not in a git repository"
    exit 1
fi

# Get current version
CURRENT_VERSION=$(git describe --tags --abbrev=0 2>/dev/null || echo "v0.0.0")
CURRENT_VERSION=${CURRENT_VERSION#v}

# Parse version
IFS='.' read -r MAJOR MINOR PATCH <<< "$CURRENT_VERSION"
MAJOR=${MAJOR:-0}
MINOR=${MINOR:-0}
PATCH=${PATCH:-0}

# Calculate new version
case "$RELEASE_TYPE" in
    major)
        MAJOR=$((MAJOR + 1))
        MINOR=0
        PATCH=0
        ;;
    minor)
        MINOR=$((MINOR + 1))
        PATCH=0
        ;;
    patch)
        PATCH=$((PATCH + 1))
        ;;
    *)
        echo "✗ Invalid release type: $RELEASE_TYPE"
        echo "  Usage: ./scripts/release.sh [major|minor|patch]"
        exit 1
        ;;
esac

NEW_VERSION="v${MAJOR}.${MINOR}.${PATCH}"

echo "════════════════════════════════════════"
echo "  Release: ${CURRENT_VERSION} → ${NEW_VERSION}"
echo "════════════════════════════════════════"
echo ""

# Check working directory is clean
if ! git diff-index --quiet HEAD --; then
    echo "✗ Working directory has uncommitted changes"
    echo "  Commit or stash changes before releasing"
    exit 1
fi

# Get current branch
CURRENT_BRANCH=$(git rev-parse --abbrev-ref HEAD)
if [ "$CURRENT_BRANCH" != "main" ] && [ "$CURRENT_BRANCH" != "master" ]; then
    echo "⚠ Warning: Not on main/master branch (current: $CURRENT_BRANCH)"
    read -p "  Continue? (y/N) " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        exit 1
    fi
fi

# Confirm release
echo "Release details:"
echo "  Version: $NEW_VERSION"
echo "  Type: $RELEASE_TYPE"
echo "  Branch: $CURRENT_BRANCH"
echo ""
read -p "Create release $NEW_VERSION? (y/N) " -n 1 -r
echo
if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    echo "✗ Cancelled"
    exit 1
fi

# Create and push tag
echo ""
echo "→ Creating git tag..."
git tag -a "$NEW_VERSION" -m "Release $NEW_VERSION"

echo "→ Pushing tag..."
git push origin "$NEW_VERSION"

echo ""
echo "════════════════════════════════════════"
echo "  ✓ Release created: $NEW_VERSION"
echo "════════════════════════════════════════"
echo ""
echo "  GitHub Actions will:"
echo "  1. Build binaries for all platforms"
echo "  2. Create checksums"
echo "  3. Generate release notes"
echo "  4. Create GitHub Release"
echo ""
echo "  Monitor progress:"
echo "  https://github.com/murtadanazar/go-api-template/actions"
echo ""
