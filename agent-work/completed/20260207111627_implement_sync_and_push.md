# Implement Sync and Push

## Status: completed 20260207111843

## Context
Push command needs to copy tracked files to synced-files folder (in MD5-hashed subfolders), clean state each time, create manifest, then git add/commit/push.

## Value Proposition
- Files organized by MD5 hash of tilde path (collision resistance)
- Clean slate on each push (no stale files)
- Auto-commit with descriptive message
- Full git sync workflow

## Alternatives considered
- Keep files flat: Could have name collisions
- Don't clean: Stale files accumulate
- **MD5 folders + clean (chosen)**: Organized, fresh each time

## Todos
- [x] Add MD5 helper function
- [x] Add SyncFiles method to JsonConfig
- [x] Add CleanSyncedFolder method
- [x] Add CreateManifest method
- [x] Add Add/Commit/Push to GitRunner
- [x] Update push command to use Sync then push
- [x] Build and test

## Notes
