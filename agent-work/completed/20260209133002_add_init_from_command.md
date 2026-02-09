# Add init-from Command

## Status: completed 20260209133148

## Context
Users need to clone an existing config-sync repository on new machines. Currently they must manually clone and set up. Adding `init-from` command streamlines onboarding.

## Value Proposition
- One command to clone repo and set up config-sync
- Creates ~/.config-sync from existing git repository
- Preserves existing security checks (public repo warnings)
- Ready to use immediately after clone

## Alternatives considered
- Manual clone + set-origin-repo: More steps, error-prone
- Auto-detect repo from ssh config: Too complex, invasive
- init-from with --force flag: Same pattern as existing commands

## Todos
- [x] Add `init-from` cobra command to main.go
- [x] Implement Clone() method in GitRunner
- [x] Add public repo check before clone (same as set-origin-repo)
- [x] Update README with usage example
- [x] Test the command

## Notes
