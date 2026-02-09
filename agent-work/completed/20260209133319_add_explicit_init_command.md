# Add Explicit Init Command

## Status: completed 20260209133610

## Context
Currently config auto-initializes silently on first command run. This is unclear and can create unexpected side effects. Users should explicitly initialize to understand what's being created.

## Value Proposition
- Clear user intent - no silent folder/git repo creation
- Better error messages guide users to proper initialization
- Distinguishes between "fresh start" vs "clone existing"
- More secure - users are aware of initialization

## Alternatives considered
- Keep auto-init: Convenient but unclear
- Auto-init with prompt: Still implicit, adds friction
- Explicit init (chosen): Clear intent, best UX

## Todos
- [x] Add `init` command that creates local git repo
- [x] Remove auto-init from PersistentPreRunE
- [x] Add isInitialized check with helpful error messages
- [x] Update init-from to handle pre-existing config
- [x] Update README with new workflow
- [x] Build and test

## Notes
