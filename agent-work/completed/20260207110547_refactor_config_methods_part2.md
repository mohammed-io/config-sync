# Refactor Config Methods Part 2

## Status: completed 20260207110646

## Context
trackFiles/untrackFiles operate on JsonConfig but are standalone functions. initializeConfig should also be a method that receives the config folder path.

## Value Proposition
- All config operations as methods on JsonConfig
- Consistent API: config.Track(), config.Untrack()
- Initialize as method: JsonConfig{}.Initialize(folder)

## Alternatives considered
- Keep as functions: Simpler but less cohesive
- **Methods (chosen)**: Better encapsulation, OOP-style

## Todos
- [x] Make Initialize a method on JsonConfig
- [x] Move trackFiles to Track() method
- [x] Move untrackFiles to Untrack() method
- [x] Update all callers
- [x] Build and test

## Notes
