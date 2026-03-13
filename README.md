# sjson

[![Go Doc](https://godoc.org/github.com/strugglerx/sjson?status.svg)](https://godoc.org/github.com/strugglerx/sjson)
[![Production Ready](https://img.shields.io/badge/production-ready-blue.svg)](https://github.com/strugglerx/sjson)
[![License](https://img.shields.io/github/license/strugglerx/sjson.svg?style=flat)](https://github.com/strugglerx/sjson)

[中文文档](./README_ZH.md)

> Expand nested JSON strings inside structs, maps, or generic Go values into valid JSON output.

`sjson` is designed for data pipelines where some fields are stored as JSON strings instead of real JSON objects or arrays. It marshals your Go value first, then scans the encoded output and expands string values that are actually valid JSON containers.

## Recommended Default

If you just want the default recommended API, use `StringWithJsonScanToString(v)`.

- It is the main high-performance path in this package.
- It returns a final JSON `string`, which is the most common output form.
- If you prefer explicit error handling over panic-style behavior, use `StringWithJsonScanToStringE(v)`.

## Typical Use Cases

- Database fields stored as stringified JSON, such as MySQL or PostgreSQL text columns.
- Third-party APIs that return JSON-in-JSON payloads.
- Audit logs, event payloads, or ETL records that embed object or array content as strings.
- Admin tools that need readable JSON output without manual post-processing.

For packaged integration examples, see [examples/README.md](./examples/README.md).
Available scenarios now include `cli`, `http-server`, `wasm`, `js-wrapper`, and `js-wrapper-vite`.

## What It Does

Input:

```json
{
  "profile": "{\"name\":\"alice\",\"tags\":[\"go\",\"json\"]}",
  "events": "[{\"type\":\"login\",\"ok\":true}]",
  "note": "[not real json]"
}
```

Output:

```json
{
  "profile": {"name":"alice","tags":["go","json"]},
  "events": [{"type":"login","ok":true}],
  "note": "[not real json]"
}
```

Valid nested JSON strings are expanded. Normal strings stay as strings.

## Installation

```bash
go get -u -v github.com/strugglerx/sjson
```

For local development helpers, this repo also includes a `Makefile`:

```bash
make test
make examples
```

## API

| Function | Description | Notes |
|:---|:---|:---|
| `ToJsonByte(v)` | Marshal to `[]byte` | Must-style API, panics on encode error |
| `ToJsonByteE(v)` | Marshal to `[]byte` | Returns `error` |
| `ToJsonString(v)` | Marshal to `string` | Must-style API, panics on encode error |
| `ToJsonStringE(v)` | Marshal to `string` | Returns `error` |
| `StringWithJsonScanToString(v)` | Expand nested JSON strings and return `string` | Recommended default |
| `StringWithJsonScanToStringE(v)` | Same as above with `error` return | Stable API for callers that avoid panic |
| `StringWithJsonScanToBytes(v)` | Expand nested JSON strings and return `[]byte` | Useful when caller already works with bytes |
| `StringWithJsonScanToBytesE(v)` | Same as above with `error` return | Stable bytes API |
| `StringWithJsonSafetyRegexToString(v)` | Legacy regex implementation | Slower, kept mainly for comparison |

## Usage

```go
package main

import (
	"fmt"

	"github.com/strugglerx/sjson"
)

func main() {
	data := map[string]interface{}{
		"user":    "{\"id\":1,\"name\":\"alice\"}",
		"roles":   "[\"admin\",\"dev\"]",
		"comment": "[not real json]",
	}

	out := sjson.StringWithJsonScanToString(data)
	fmt.Println(out)
}
```

If you do not want panic-style behavior:

```go
out, err := sjson.StringWithJsonScanToStringE(data)
if err != nil {
	panic(err)
}
fmt.Println(out)
```

If your downstream code already needs bytes:

```go
b, err := sjson.StringWithJsonScanToBytesE(data)
if err != nil {
	panic(err)
}
_ = b
```

## Behavior Notes

- Only string values that decode into a valid JSON object or array are expanded.
- Strings like `"[hello]"`, `"{not-json}"`, or other invalid JSON-looking content stay unchanged.
- Expansion is based on the JSON-encoded form of the value, not the raw Go string bytes.
- The scanner expands one level in the final encoded output. It does not recursively keep decoding forever.
- Encoding still follows Go's standard `encoding/json` behavior, including UTF-8 normalization rules.

## Performance

Current local benchmark environment:

- Machine: Apple M1 Pro
- Go: module target is `go 1.16`
- Dataset: roughly `1MB`

Latest benchmark snapshot:

| Benchmark | Time | Memory | Allocs |
|:---|:---|:---|:---|
| `SafetyRegex_1MB` | `~115.7 ms/op` | `~133.7 MB/op` | `23656 allocs/op` |
| `Scanner_1MB` | `~3.62 ms/op` | `~1.57 MB/op` | `13766 allocs/op` |

That is roughly:

- `~32x` faster than the legacy regex path
- about `98%+` less memory than the regex path

## Why It Is Fast

- Single-pass byte scanning instead of regex-heavy rewriting
- `sync.Pool` reuse for temporary buffers
- Segment-based copying instead of byte-by-byte output writes on hot paths
- Validation only when a candidate string actually looks like a JSON container

## Safety And Testing

The package is covered by:

- Standard unit tests
- Complex hand-written edge cases
- Benchmarks
- Go fuzz testing for random string inputs

Recent fuzz runs completed successfully against the scan path without producing invalid JSON or panics in the tested window.

## Concurrency And Memory

- Pool usage is concurrency-safe
- Large pooled objects are capped before being returned, which helps avoid buffer growth sticking around forever
- `sync.Pool` remains GC-friendly for server workloads

## When To Use This Package

Use `sjson` when:

- you need a JSON string output
- some fields may contain JSON strings that should become real nested JSON
- you want a lightweight solution without first decoding into custom structs

It may be the wrong tool when:

- you need full recursive semantic transformation of arbitrary JSON documents
- you want schema-aware decoding into strongly typed Go structs
- your input is already a clean, correctly typed Go value graph

## License

`sjson` is licensed under the [MIT License](LICENSE).
