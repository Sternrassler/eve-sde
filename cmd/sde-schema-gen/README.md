# sde-schema-gen

Go CLI tool to generate type-safe Go structs from EVE Online SDE JSONL files.

## Features

- **Multi-line analysis**: Scans up to 100 JSONL lines (configurable) to infer complete schemas
- **LocalizedText detection**: Automatically recognizes EVE's 8-language text objects
- **Smart type inference**: Detects int64, float64, bool, string, maps, slices
- **CamelCase conversion**: Handles snake_case, camelCase, ID/NPC/CEO abbreviations
- **Template-based**: Uses Go's text/template for clean code generation
- **Nested structures**: Handles nesting with `map[string]interface{}` for maximum compatibility

## Installation

```bash
go build -o bin/sde-schema-gen ./cmd/sde-schema-gen
```

## Usage

```bash
# Generate schemas from JSONL data
./bin/sde-schema-gen \
  -input data/jsonl \
  -output internal/schema/types \
  -lines 100 \
  -v

# Or via wrapper script
./scripts/fetch-schemas.sh --refresh --verbose
```

### Options

- `-input DIR`: JSONL input directory (default: `data/jsonl`)
- `-output DIR`: Go output directory (default: `internal/schema/types`)
- `-lines N`: Max lines to analyze per file (default: 100)
- `-v`: Verbose logging

## Output

Generates two types of files:

1. **common.go**: Shared `LocalizedText` type
2. **{schema}.go**: One file per JSONL schema (e.g., `blueprints.go`)

### Example Output

```go
// blueprints.go
package types

type Blueprints struct {
    Key                int64                  `json:"_key,omitempty"`
    Activities         map[string]interface{} `json:"activities,omitempty"`
    BlueprintTypeID    int64                  `json:"blueprintTypeID,omitempty"`
    MaxProductionLimit int64                  `json:"maxProductionLimit,omitempty"`
}
```

## Architecture

- **analyzer.go**: JSONL parsing & schema extraction
- **types.go**: CamelCase conversion & naming utilities
- **writer.go**: Template-based Go code generation
- **main.go**: CLI entry point

## Type Inference Rules

| JSON Type | Go Type | Notes |
|-----------|---------|-------|
| `null` | `interface{}` | Unknown type |
| `true`/`false` | `bool` | Boolean |
| `123` | `int64` | Integer (float64 with no fractional part) |
| `123.45` | `float64` | Floating point |
| `"text"` | `string` | String |
| `{"de":"...", "en":"..."}` | `LocalizedText` | 8-language EVE text |
| `{"key": "value"}` | `map[string]interface{}` | Generic object |
| `[1, 2, 3]` | `[]int64` | Typed array |
| `[{...}, {...}]` | `[]map[string]interface{}` | Array of objects |

## LocalizedText Recognition

Automatically detects EVE's multilingual text format:

```json
{
  "de": "German text",
  "en": "English text",
  "es": "Spanish text",
  "fr": "French text",
  "ja": "Japanese text",
  "ko": "Korean text",
  "ru": "Russian text",
  "zh": "Chinese text"
}
```

Converted to:

```go
type LocalizedText struct {
    De string `json:"de,omitempty"`
    En string `json:"en"`
    Es string `json:"es,omitempty"`
    Fr string `json:"fr,omitempty"`
    Ja string `json:"ja,omitempty"`
    Ko string `json:"ko,omitempty"`
    Ru string `json:"ru,omitempty"`
    Zh string `json:"zh,omitempty"`
}
```

## Why Not Python?

Previous approach used embedded Python in bash script - this was:

- ❌ Hard to debug (Python-in-bash)
- ❌ Incomplete recursion (missing nested struct definitions)
- ❌ Complex type handling
- ❌ No proper testing

New Go-based tool provides:

- ✅ Type-safe code generation
- ✅ Debuggable with standard Go tools
- ✅ Template-based approach
- ✅ Reusable as library
- ✅ Unit testable
