# Add GitRunner Interface

## Status: completed 20260207110022

## Context
Git commands are scattered functions. Need to group them behind an interface for better organization and testability (mocking).

## Value Proposition
- Single interface for all git operations
- Easy to mock for tests
- Cleaner separation of concerns

## Alternatives considered
- Struct with methods: Harder to mock, fixed implementation
- Keep functions: Simple but not testable
- **Interface (chosen)**: Idiomatic Go, mockable, flexible

## Todos
- [x] Create GitRunner interface
- [x] Create RealGitRunner implementation
- [x] Update pullRepo/pushRepo to use interface
- [x] Build and test

## Notes
