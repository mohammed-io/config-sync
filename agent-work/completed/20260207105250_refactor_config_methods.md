# Refactor Config Methods

## Status: completed 20260207105340

## Context
Config-related functions are currently package-level. Some operations like `saveConfig` are better as methods since they operate on existing JsonConfig instances.

## Value Proposition
- More idiomatic Go: methods for operations on instances
- Better encapsulation: `config.Save()` vs `saveConfig(&config)`
- Cleaner API

## Alternatives considered
- Keep all functions: Consistent but less idiomatic for instance operations
- Make everything methods: Overkill, constructors don't make sense as methods
- **Hybrid (chosen)**: `Save()` as method, `initializeConfig()` stays as function (constructor pattern)

## Todos
- [x] Convert saveConfig to Save() method
- [x] Update all callers
- [x] Make loadConfigFromJson private (loadConfig)
- [x] Build and test

## Notes
