# Add Versioning and Release Script

## Status: completed 20260207123730

## Context
Users installing with `go install @latest` don't always get the newest version due to Go's module caching. Without proper semantic version tags, `@latest` behavior is undefined.

## Value Proposition
- Version flag displays current version
- Release script automates tagging and version updates
- README updated with install command using latest tag
- Proper semver tags ensure `@latest` works reliably

## Alternatives considered
1. **Manual tagging** - Error-prone, easy to forget steps
2. **GitHub Actions** - Overkill for simple releases
3. **Shell script** (chosen) - Simple, transparent, easy to modify

## Todos
- [x] Add version variable to main.go
- [x] Add --version flag to CLI
- [x] Display version in help header
- [x] Create scripts/release.sh script
- [ ] Update README with v0.0.1 tag
- [x] Tag and push v0.0.1

## Notes
Version stored in single source of truth (version.go)
