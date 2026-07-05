# DiffForge

Universal Diff & Merge Toolkit — a single CLI tool that diffs anything: text, JSON, YAML, CSV, binary files, and entire directories.

## Features

- **Text diffing** — LCS-based line diff with unified, side-by-side, context, and minimal output formats
- **JSON structural diff** — Path-based comparison showing exactly what changed in nested objects
- **YAML structural diff** — Same structural diff engine for YAML files
- **CSV row-level diff** — Cell-by-cell comparison with column name tracking
- **Binary diff** — Hex dump comparison with offset tracking
- **Directory diff** — Recursive comparison using content hashing (SHA-256)
- **Three-way merge** — Merge changes from two branches with conflict detection
- **Ignore rules** — Skip whitespace, blank lines, comments, and custom patterns
- **Word-level diff** — Fine-grained diff within changed lines

## Installation

```bash
go install github.com/EdgarOrtegaRamirez/diffforge/cmd/diffforge@latest
```

Or build from source:

```bash
git clone https://github.com/EdgarOrtegaRamirez/diffforge
cd diffforge
go build -o diffforge ./cmd/diffforge
```

## Quick Start

### Text diff

```bash
# Unified diff (default)
diffforge text old.txt new.txt

# Side-by-side
diffforge text old.txt new.txt -f side-by-side -w 120

# Ignore whitespace changes
diffforge text old.txt new.txt -i whitespace
```

### JSON diff

```bash
# Structural diff
diffforge json old.json new.json

# Shows paths to changed values
diffforge json config-v1.json config-v2.json
```

### YAML diff

```bash
diffforge yaml old.yaml new.yaml
```

### CSV diff

```bash
# Row-level diff with column tracking
diffforge csv data-v1.csv data-v2.csv
```

### Binary diff

```bash
# Hex dump comparison
diffforge binary old.bin new.bin
```

### Directory diff

```bash
# Recursive directory comparison
diffforge dir ./v1/ ./v2/

# Ignore whitespace in text files
diffforge dir ./v1/ ./v2/ -i whitespace
```

### Three-way merge

```bash
# Merge changes from two branches
diffforge merge base.txt branch-a.txt branch-b.txt
```

## Output Formats

| Format | Flag | Description |
|--------|------|-------------|
| `unified` | `-f unified` | Standard unified diff (default) |
| `side-by-side` | `-f side-by-side` | Two-column comparison |
| `context` | `-f context` | Context diff format |
| `minimal` | `-f minimal` | Compact change summary |
| `json` | `-f json` | Machine-readable diff (for programmatic use) |

## Ignore Rules

| Rule | Flag | Description |
|------|------|-------------|
| `whitespace` | `-i whitespace` | Ignore leading/trailing whitespace |
| `blank-lines` | `-i blank-lines` | Ignore blank lines |
| `comments` | `-i comments` | Ignore single-line comments (`//`, `#`, `--`, `;`) |
| `code` | `-i code` | Ignore whitespace + blank lines + comments |

## Architecture

```
diffforge/
├── cmd/diffforge/          # CLI entry point
│   ├── main.go             # Command routing and flag parsing
│   └── main_test.go        # Integration tests
├── internal/
│   ├── diff/               # Core diff engine
│   │   ├── lcs.go          # LCS algorithm + text diffing
│   │   ├── json.go         # JSON structural diff
│   │   ├── yaml.go         # YAML structural diff
│   │   ├── csv.go          # CSV row-level diff
│   │   ├── binary.go       # Binary byte-level diff
│   │   └── dir.go          # Directory recursive diff
│   ├── format/             # Output formatters
│   │   └── format.go       # Unified, side-by-side, context, minimal
│   ├── merge/              # Three-way merge
│   │   └── merge.go        # Merge algorithm with conflict detection
│   └── rules/              # Ignore rules
│       └── rules.go        # Whitespace, comment, blank line rules
├── go.mod
├── go.sum
├── LICENSE                 # MIT License
├── README.md
└── AGENTS.md               # AI agent instructions
```

## Algorithms

- **LCS (Longest Common Subsequence)** — Core algorithm for text diffing, finding the minimum edit distance between two sequences
- **SHA-256 Content Hashing** — Used for directory diff to detect file modifications without content comparison
- **Recursive Structural Comparison** — JSON/YAML diffs traverse the document tree, comparing values at each path
- **Three-Way Merge** — Compares base against both old and new, accepting changes from one side or flagging conflicts when both changed

## Testing

```bash
# Run all tests
go test ./...

# Verbose output
go test ./... -v

# Run specific package tests
go test ./internal/diff/ -run TestTextDiff
go test ./internal/merge/ -run TestThreeWayMerge
```

## License

MIT License — see [LICENSE](LICENSE) for details.
