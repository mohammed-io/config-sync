# Add file comparison to check-updates command

## Status: completed 20260209180440

## Context
The `check-updates` command only checks git status on the `synced-files/` directory. If a source file (like ~/.vimrc) is modified after the last sync, the check won't detect it because `synced-files/` still contains the old version. This creates a stale-cache problem.

## Value Proposition
- `check-updates` will detect changes to tracked source files since last sync
- Non-destructive: doesn't modify `synced-files/` until user runs `push`
- Fast: uses mtime+size comparison (microseconds per file)
- Accurate: warns user when their config is out of sync

## Alternatives considered
- Copy files during check: Destructive, slower, confusing if user doesn't push
- Full hash comparison: Slower, unnecessary when mtime+size is sufficient
- **mtime+size comparison (chosen)**: Fast, accurate enough, non-destructive

## Todos
- [x] Add HasUnsyncedChanges() method to JsonConfig
- [x] Update checkUpdatesCmd to call HasUnsyncedChanges()
- [x] Test the implementation

## Notes
