# Singleton Config Init

## Status: completed 20260207110903

## Context
Config is initialized multiple times across commands. Should initialize once at startup, pass to commands, and fail if methods called before initialization.

## Value Proposition
- Single initialization at startup
- Methods fail if config not initialized
- Pass config to commands instead of re-initializing

## Alternatives considered
- Keep current: Re-init on every command - wasteful
- Global variable: Simple but makes testing harder
- **Initialized flag + preflight check (chosen)**: Explicit, fails fast

## Todos
- [x] Add initialized flag to JsonConfig
- [x] Add preflight checks to methods
- [x] Create global config instance
- [x] Initialize in main, pass to commands
- [x] Build and test

## Notes
