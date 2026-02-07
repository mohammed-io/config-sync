# Add Directory Support

## Status: completed 20260207124541

## Context
Currently only files are supported. When trying to track a directory, it fails with "is a directory" error because the code uses `os.ReadFile` which doesn't work on directories.

## Value Proposition
- Support tracking entire directories (e.g., ~/.agents, ~/.config)
- Recursively copy directory contents
- Maintain directory structure in synced-files

## Alternatives considered
1. **Recursive file walking** (chosen) - Copies all files, maintains structure
2. **Archive/tar the directory** - Simpler but harder to restore
3. **Symlinks only** - Won't work across different machines

## Todos
- [x] Check if path is file or directory
- [x] Implement recursive directory copy
- [x] Update SyncFiles to use copyDir for directories
- [x] Test with directory

## Notes
Use filepath.Walk or io.Copy for directories
