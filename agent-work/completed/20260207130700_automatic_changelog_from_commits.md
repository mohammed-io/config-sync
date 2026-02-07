# Auto-Generate Changelog from Git Commits

## Status: completed 20260207130742

## Context
Currently release script prompts for manual changelog entry. User wants automatic parsing of git commits since last tag.

## Value Proposition
- Automatic changelog from commits
- No manual entry needed
- Enforces conventional commit format

## Alternatives considered
1. **Parse git log** (chosen) - Get commits since last tag, categorize by type
2. **Keep manual** - Flexible but tedious
3. **Hybrid** - Show parsed commits, allow editing

## Todos
- [x] Get latest tag to determine commit range
- [x] Parse git log for conventional commits
- [x] Categorize by type (feat, fix, chore, etc.)
- [x] Update release.sh with automatic parsing

## Notes
Conventional commits: feat:, fix:, chore:, docs:, refactor:, test:
