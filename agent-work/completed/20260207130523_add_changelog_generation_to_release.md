# Add Changelog Generation to Release Script

## Status: completed 20260207130553

## Context
Release script needs to generate/update CHANGELOG.md with each release. This helps users track what changed between versions.

## Value Proposition
- Automatic changelog updates on each release
- Consistent changelog format
- Track features, fixes, and changes over time

## Alternatives considered
1. **Manual entry** (chosen) - Script prompts for changes, appends to CHANGELOG.md
2. **Parse git commits** - Automatic but requires conventional commits
3. **No changelog** - Users have to check git log

## Todos
- [x] Add CHANGELOG.md generation to release.sh
- [x] Prompt for changelog entries
- [x] Create initial CHANGELOG.md structure

## Notes
Keep it simple: append new version with user-provided notes
