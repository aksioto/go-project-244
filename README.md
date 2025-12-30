# gendiff

gendiff is a CLI utility and Go package for comparing two configuration files and showing the differences between them.

[![Actions Status](https://github.com/aksioto/go-project-244/actions/workflows/hexlet-check.yml/badge.svg)](https://github.com/aksioto/go-project-244/actions)
[![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=aksioto_go-project-244&metric=alert_status)](https://sonarcloud.io/summary/new_code?id=aksioto_go-project-244)
[![Coverage](https://sonarcloud.io/api/project_badges/measure?project=aksioto_go-project-244&metric=coverage)](https://sonarcloud.io/summary/new_code?id=aksioto_go-project-244)
[![CI](https://github.com/aksioto/go-project-244/actions/workflows/ci.yml/badge.svg)](https://github.com/aksioto/go-project-244/actions/workflows/ci.yml)

## Features

- Supported input formats: **JSON**, **YAML**
- Works with deeply nested data structures
- Output formats: **stylish** (default), **plain**, **json**

## Installation

```bash
git clone https://github.com/aksioto/go-project-244.git
cd go-project-244
make build
```

## Usage

### Help

```bash
./bin/gendiff --help

NAME:
   gendiff - Compares two configuration files and shows a difference.

USAGE:
   gendiff [global options]

GLOBAL OPTIONS:
   --format string, -f string  output format (default: "stylish")
   --help, -h                  show help
```

### Examples

**Stylish output (default):**

```bash
./bin/gendiff file1.json file2.json
```

```
{
  - follow: false
    host: hexlet.io
  - proxy: 123.234.53.22
  - timeout: 50
  + timeout: 20
  + verbose: true
}
```

**Plain format:**

```bash
./bin/gendiff --format plain file1.json file2.json
```

```
Property 'follow' was removed
Property 'proxy' was removed
Property 'timeout' was updated. From 50 to 20
Property 'verbose' was added with value: true
```

**JSON format:**

```bash
./bin/gendiff --format json file1.json file2.json
```

```json
[
  {
    "type": "removed",
    "key": "follow",
    "value": false
  },
  ...
]
```

## Library usage

```go
package main

import (
    "fmt"

    "code"
    "code/internal/parser"
)

func main() {
    // Simple usage
    result, err := code.GenDiff("file1.json", "file2.yml", "stylish")
    if err != nil {
        panic(err)
    }
    fmt.Println(result)

    // Custom configuration (Options pattern)
    customParsers := parser.NewFileParser()
    customParsers.Add(&parser.JSONParser{}, ".json")
    customParsers.Add(&parser.YAMLParser{}, ".yaml", ".yml")

    differ := code.NewDiffer(code.WithFileParser(customParsers))
    result, err = differ.GetDiff("file1.json", "file2.json", "plain")
    if err != nil {
        panic(err)
    }
    fmt.Println(result)
}
```

## Development

```bash
# Run tests
make test

# Run linter
make lint

# Build the binary
make build
```

## Demo

[![asciicast](https://asciinema.org/a/joqZt3OryVvb1Mk08t5NLLKJo.svg)](https://asciinema.org/a/joqZt3OryVvb1Mk08t5NLLKJo)