# Add Restore Files After Pull

## Status: completed 20260207125541

## Context
Currently `pull` only does `git pull` - files stay in `~/.config-sync/synced-files/`. Users need to manually copy files to their destinations, which is tedious.

## Value Proposition
- After pull, automatically copy files/dirs to their original locations
- Full sync experience: pull restores everything to correct paths

## Alternatives considered
1. **Modify pull command** (chosen) - Automatic restore after pull
2. **Separate restore command** - More explicit but extra step
3. **Ask user each time** - Safer but annoying

## Todos
- [x] Add RestoreFiles method to JsonConfig
- [x] Update pullCmd to call RestoreFiles after git pull
- [x] Test restore with files and directories

## Notes
Need to reverse the copy logic: copy FROM hash folder TO original path
