# sjson

[![Go Doc](https://godoc.org/github.com/strugglerx/sjson?status.svg)](https://godoc.org/github.com/strugglerx/sjson)
[![Production Ready](https://img.shields.io/badge/production-ready-blue.svg)](https://github.com/strugglerx/sjson)
[![License](https://img.shields.io/github/license/strugglerx/sjson.svg?style=flat)](https://github.com/strugglerx/sjson)

[中文文档](./README_ZH.md)

> Expand nested JSON strings inside structs, maps, or generic Go values into valid JSON output.

`sjson` is designed for data pipelines where some fields are stored as JSON strings instead of real JSON objects or arrays. It marshals your Go value first, then scans the encoded output and expands string values that are actually valid JSON containers.

## Depth-Aware Decoding

`sjson` now also supports depth-aware recursive decoding for payloads that contain JSON strings inside JSON strings.

If a field looks like `"data": "{\"inner\": \"{\\\"deep\\\": 1}\"}"`, you can choose how many extra levels should be expanded:

- `depth = 0`: only expand the first JSON string layer
- `depth = 1`: expand one more nested stringified JSON layer
- `depth = 2`: expand two more nested layers
- `depth = -1`: keep expanding until no more valid nested JSON containers remain

Important: the depth APIs are a separate feature path. The default scanner APIs stay on the original fast path and are still the recommended choice unless you explicitly need recursive expansion.

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
	"profile": { "name": "alice", "tags": ["go", "json"] },
	"events": [{ "type": "login", "ok": true }],
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

| Function                                     | Description                                    | Notes                                       |
| :------------------------------------------- | :--------------------------------------------- | :------------------------------------------ |
| `ToJsonByte(v)`                              | Marshal to `[]byte`                            | Must-style API, panics on encode error      |
| `ToJsonByteE(v)`                             | Marshal to `[]byte`                            | Returns `error`                             |
| `ToJsonString(v)`                            | Marshal to `string`                            | Must-style API, panics on encode error      |
| `ToJsonStringE(v)`                           | Marshal to `string`                            | Returns `error`                             |
| `StringWithJsonScanToString(v)`              | Expand nested JSON strings and return `string` | Recommended default                         |
| `StringWithJsonScanToStringE(v)`             | Same as above with `error` return              | Stable API for callers that avoid panic     |
| `StringWithJsonScanToBytes(v)`               | Expand nested JSON strings and return `[]byte` | Useful when caller already works with bytes |
| `StringWithJsonScanToBytesE(v)`              | Same as above with `error` return              | Stable bytes API                            |
| `StringWithJsonScanDepthToString(v, depth)`  | Decode strings recursively by `depth`          | `-1` for infinite recursion                 |
| `StringWithJsonScanDepthToStringE(v, depth)` | Same as above with `error` return              | `-1` for infinite recursion                 |
| `StringWithJsonScanDepthToBytes(v, depth)`   | Decode recursively by `depth`                  | returns `[]byte`                            |
| `StringWithJsonScanDepthToBytesE(v, depth)`  | Same as above with `error` return              | returns `[]byte`                            |
| `StringWithJsonSafetyRegexToString(v)`       | Legacy regex implementation                    | Slower, kept mainly for comparison          |

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

	nestedData := map[string]interface{}{
		"data": "{\"inner\": \"{\\\"deep\\\": 1}\"}",
	}

	// Will recursively unescape the nested values
	// Output: {"data":{"inner":{"deep": 1}}}
	outDeep := sjson.StringWithJsonScanDepthToString(nestedData, -1)
	fmt.Println(outDeep)
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
- `StringWithJsonScanToString` and `StringWithJsonScanToBytes` expand one level in the final encoded output.
- If you need recursive expansion, use the `StringWithJsonScanDepth...` APIs.
- Encoding still follows Go's standard `encoding/json` behavior, including UTF-8 normalization rules.

## Performance

Current local benchmark environment:

- Machine: Apple M1 Pro
- Go: module target is `go 1.16`
- Benchmarks below include `1MB`, `10MB`, and `100MB` payloads

Latest benchmark snapshot:

### 1. Decoding Standard Payloads (vs Legacy Regex)

| Benchmark          | Time           | Memory         | Allocs              |
| :----------------- | :------------- | :------------- | :------------------ |
| `SafetyRegex_1MB`  | `~115.7 ms/op` | `~133.7 MB/op` | `23656 allocs/op`   |
| `Scanner_1MB`      | `~3.93 ms/op`  | `~1.53 MB/op`  | `13765 allocs/op`   |
| `SafetyRegex_10MB` | `~8.91 s/op`   | `~12.25 GB/op` | `237446 allocs/op`  |
| `Scanner_10MB`     | `~39.6 ms/op`  | `~20.8 MB/op`  | `137546 allocs/op`  |
| `Scanner_100MB`    | `~376 ms/op`   | `~252.7 MB/op` | `1375061 allocs/op` |

That is roughly:

- `~32x` faster than the legacy regex path
- about `98%+` less memory than the regex path
- still practical on very large payloads when using the scanner path

### 2. Depth APIs

Depth-aware decoding is available for payloads that truly need recursive expansion, but it is intentionally kept on a separate API path:

```json
{
	"data": "{\"inner1\": \"{\\\"inner2\\\": \\\"{\\\\\\\"inner3\\\\\\\": \\\\\\\"{\\\\\\\\\\\\\\\"deep\\\\\\\\\\\\\\\": 1}\\\\\\\"}\\\"}\"}"
}
```

- `StringWithJsonScanDepthToString(v, 0)` keeps the same first-layer behavior as the normal scanner path
- `StringWithJsonScanDepthToString(v, 1)` expands one more nested layer
- `StringWithJsonScanDepthToString(v, -1)` keeps expanding until the nested JSON string chain ends

Use the depth APIs only when you actually need recursive decoding. For normal API responses, DB rows, and event payload cleanup, the default non-depth scanner path is still the best choice.

## Why It Is Fast

- Single-pass byte scanning instead of regex-heavy rewriting
- `sync.Pool` reuse for temporary buffers
- Segment-based copying instead of byte-by-byte output writes on hot paths
- Validation only when a candidate string actually looks like a JSON container

## Large Payload Conclusion

- For normal backend responses and internal tools, the scanner path is already fast enough to use with confidence.
- On `10MB` payloads, the scanner path stayed in the `~35ms` range, while the legacy regex path moved into multi-second territory.
- On `100MB` payloads, the scanner path still completed in sub-second time on the test machine.
- After adding the depth feature, the default scanner path still keeps its own fast implementation and did not show a meaningful regression in local benchmark comparison.
- The legacy regex path became impractical at `100MB` scale, so the scanner path should be treated as the real production path.

In short: if your workload includes large or messy JSON payloads, `StringWithJsonScanToString` is the path you should use by default.

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
