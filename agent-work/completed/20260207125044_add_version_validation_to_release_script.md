# Add Version Validation to Release Script

## Status: completed 20260207125222

## Context
Release script accepts any version string, which can lead to invalid tags. User should only be able to pass semantic versions like "0.0.1", not "v0.0.1" or "0.0" or invalid formats.

## Value Proposition
- Validate version format before running release
- Prevent invalid git tags
- Clear error message for invalid formats

## Alternatives considered
1. **Regex validation** (chosen) - Simple, reliable, standard semver check
2. **Go validation** - Overkill to call Go binary for validation
3. **No validation** - User's manual error, not ideal UX

## Todos
- [x] Add semver regex validation to release.sh
- [x] Reject 'v' prefix strictly
- [x] Test validation

## Notes
Semver regex: ^[0-9]+\.[0-9]+\.[0-9]+$
