# Add check-for-updates Command

## Status: completed 20260209174436

## Context
Users need a quick way to check if their config is out of sync - either they have local unpushed changes or there are remote changes they haven't pulled. This should be lightweight and not require initialization.

## Value Proposition
- Fast check with no side effects
- Works even if not initialized (just exits silently)
- Clear messaging about what needs to be done (pull/push/both)
- Useful for prompts/shell startup scripts

## Alternatives considered
- Full sync status: Too heavy, requires network
- Git status parsing: Complex, this wraps it cleanly
- Always show in prompt: Annoying, opt-in better

## Todos
- [x] Add GitRunner methods: HasUnpushedChanges(), HasUnpulledChanges()
- [x] Add check-updates command to main.go
- [x] Skip init check for this command (works regardless of init state)
- [x] Update README with usage example
- [x] Build and test

## Notes
