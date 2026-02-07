# Add Merge Conflict Handling

## Status: completed 20260207131215

## Context
Git pull/push can fail with merge conflicts. Currently the tool just exits with error. User needs to be guided to resolve conflicts in the config folder.

## Value Proposition
- Detect merge conflicts during pull/push
- Guide user to resolve conflicts in ~/.config-sync
- Clear error messages with next steps

## Alternatives considered
1. **Abort on conflict** (chosen) - Stop and tell user to fix manually
2. **Auto-merge with strategy** - Risky, could lose data
3. **Interactive resolution** - Too complex for a simple tool

## Todos
- [x] Detect git conflict exit code
- [ ] Add helpful error message with folder path
- [x] Test conflict detection

## Notes
Git returns exit code 1 on merge conflicts
