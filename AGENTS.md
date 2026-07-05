# DiffForge — Agent Instructions

## Project Overview
DiffForge is a universal diff & merge toolkit in Go. It diffs text, JSON, YAML, CSV, binary files, and directories with multiple output formats.

## Build & Test
```bash
go build ./cmd/diffforge    # Build CLI
go test ./...               # Run all tests (42 tests)
go vet ./...                # Static analysis
```

## Architecture
- `internal/diff/` — Core diff algorithms (LCS, JSON/YAML/CSV/Binary/Directory)
- `internal/format/` — Output formatters (unified, side-by-side, context, minimal)
- `internal/merge/` — Three-way merge engine
- `internal/rules/` — Ignore rules (whitespace, comments, blank lines)
- `cmd/diffforge/` — CLI entry point

## Key Files
- `internal/diff/lcs.go` — LCS algorithm, text diff, word diff
- `internal/diff/json.go` — JSON structural diff with path tracking
- `internal/diff/csv.go` — CSV row-level diff
- `internal/diff/binary.go` — Binary hex dump diff
- `internal/diff/dir.go` — Directory recursive diff with SHA-256 hashing
- `internal/merge/merge.go` — Three-way merge with conflict detection
- `internal/rules/rules.go` — Configurable ignore rules

## Adding New Diff Types
1. Create `internal/diff/<type>.go` with a `<Type>Diff()` function
2. Return a result type with ops and stats
3. Add a formatter function
4. Add CLI command in `cmd/diffforge/main.go`
5. Add tests in `internal/diff/diff_test.go`

## Dependencies
- `gopkg.in/yaml.v3` — YAML parsing (only external dependency)

## Conventions
- All public functions must have tests
- Use table-driven tests
- Keep functions small and focused
- Error messages should be descriptive
