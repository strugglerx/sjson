# sjson

[![Go Doc](https://godoc.org/github.com/strugglerx/sjson?status.svg)](https://godoc.org/github.com/strugglerx/sjson)
[![Production Ready](https://img.shields.io/badge/production-ready-blue.svg)](https://github.com/strugglerx/sjson)
[![License](https://img.shields.io/github/license/strugglerx/sjson.svg?style=flat)](https://github.com/strugglerx/sjson)

[中文文档](./ReadMe_zh.md)

> Expand nested JSON strings within structs/maps into valid JSON. Commonly used for MySQL JSON fields stored as strings.

## Example

Convert this (where fields are JSON strings):

```json
{
    "key1": "{\"key\":[\"CN\", \"EN\"]}",
    "key2": "[{\"url\":\"https://example.com\", \"desc\": \"Line\nBreak\"}]",
    "nick": "[Wow]"
}
```

Into valid JSON (nested JSON automatically expanded, normal strings untouched):

```json
{
    "key1": {"key":["CN", "EN"]},
    "key1": [{"url":"https://example.com", "desc":"Line\nBreak"}],
    "nick": "[Wow]"
}
```

## Installation

```bash
go get -u -v github.com/strugglerx/sjson
```

## API Reference

| Function | Description | Recommended |
|:---|:---|:---|
| `ToJsonByte(v)` | Standard JSON Marshal -> `[]byte` | ✅ |
| `ToJsonString(v)` | Standard JSON Marshal -> `string` | ✅ |
| `StringWithJsonScanToString(v)` | Extreme optimized single-pass scanner | ⭐ **Recommended** |
| `StringWithJsonSafetyRegexToString(v)` | Regex-based expansion (Legacy/Comparison) | ⚠️ Slower |

## Performance

Environment: Apple M1 Pro, Benchmark with **~1MB** dataset:

| Metric | `SafetyRegex` (Legacy) | **`ScanToString` (New)** | Boost |
|:---|:---|:---|:---|
| **Latency** | ~115 ms | **~3.6 ms** | **🚀 32x Faster** |
| **Memory Alloc** | ~133 MB | **~1.06 MB** | **🔥 99% Reduced** |
| **Alloc Count** | 23,654 | **13,762** | **40% Lower** |

### Why is it so fast?

1. **Zero Regex Engine Overhead**: Replaced heavy regex with a single O(n) byte scanner.
2. **Full Bytestream Pipeline**: Reduced unnecessary string-to-byte conversions.
3. **Scratch Buffer Pooling**: Used `sync.Pool` for temporary validation buffers.
4. **Segmented Copy**: Optimized the unescape loop with bulk segment copies.

## Memory Safety

- **GC Friendly**: Uses `sync.Pool` which is automatically reclaimed by GC.
- **Capacity Guard**: Buffers exceeding 4KB are discarded instead of pooled to prevent memory bloat.
- **Thread Safe**: All pooled operations are safe for concurrent use.

## Usage

```golang
package main

import (
    "fmt"
    "github.com/strugglerx/sjson"
)

func main() {
    data := map[string]interface{}{
        "nick":  "[Wow]",
        "extra": "{\"score\":100, \"tags\":[\"go\"]}",
    }

    // ⭐ Recommended: Best performance, validated by json.Valid
    result := sjson.StringWithJsonScanToString(data)
    fmt.Println(result)
}
```

## License

`sjson` is licensed under the [MIT License](LICENSE), 100% free and open-source, forever.