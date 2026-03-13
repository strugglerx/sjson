# sjson

[![Go Doc](https://godoc.org/github.com/strugglerx/sjson?status.svg)](https://godoc.org/github.com/strugglerx/sjson)
[![Production Ready](https://img.shields.io/badge/production-ready-blue.svg)](https://github.com/strugglerx/sjson)
[![License](https://img.shields.io/github/license/strugglerx/sjson.svg?style=flat)](https://github.com/strugglerx/sjson)

> 将结构体中**以字符串形式存储的内嵌 JSON** 字段展开为合法 JSON。常用于 MySQL 等数据库中 JSON 字段以字符串存储的场景。

## 功能示例

将此结构（字段值为 JSON 字符串）：

```json
{
    "key1": "{\"key\":[\"中文\", \"english\"]}",
    "key2": "[{\"url\":\"https://example.com\", \"desc\": \"换行\n换行\"}]",
    "nick": "[天呐]"
}
```

转换为合法 JSON（内嵌 JSON 自动展开，普通字符串保持原样）：

```json
{
    "key1": {"key":["中文", "english"]},
    "key2": [{"url":"https://example.com", "desc":"换行\n换行"}],
    "nick": "[天呐]"
}
```

## 安装

```bash
go get -u -v github.com/strugglerx/sjson
```

## API 说明

| 函数 | 说明 | 推荐度 |
|:---|:---|:---|
| `ToJsonByte(v)` | 标准 JSON 序列化，返回 `[]byte` | ✅ 通用 |
| `ToJsonString(v)` | 标准 JSON 序列化，返回 `string` | ✅ 通用 |
| `StringWithJsonScanToString(v)` | 极致优化的单次扫描展开内嵌 JSON | ⭐ **推荐** |
| `StringWithJsonSafetyRegexToString(v)` | 正则展开内嵌 JSON（旧版本，对比保留） | ⚠️ 较慢 |

## 性能表现

测试环境：Apple M1 Pro，处理约 **1MB** 大数据集对比：

| 指标 | `SafetyRegex` (正则版) | **`ScanToString` (扫描版)** | 提升幅度 |
|:---|:---|:---|:---|
| **执行耗时** | ~115 ms | **~3.6 ms** | **🚀 32x 加速** |
| **内存分配** | ~133 MB | **~1.06 MB** | **🔥 减少 99%** |
| **分配次数** | 23,654次 | **13,762次** | **降低 40%** |

### 为什么这么快？

1. **零正则引擎开销**：完全摒弃了重量级的正则匹配，改为单次 O(n) 字节扫描。
2. **全链路字节流**：减少了不必要的字符串与字节数组之间的头信息拷贝。
3. **Scratch Buffer 池化**：用于内嵌 JSON 校验的临时空间使用 `sync.Pool` 复用，将内存分配降至理论极小值。
4. **分段批量拷贝**：在处理脱义字符时，寻找下一个特殊字符并批量拷贝中间段，而非逐字节处理。

## 内存安全性

代码中使用了大量的 `sync.Pool` 来优化 GC 压力，为了防止内存泄露或内存爆炸，我们实施了以下策略：

- **垃圾回收兼容**：`sync.Pool` 中的对象会在 GC 时自动释放，不会导致常驻内存（RSS）无限增长。
- **大容量熔断**：在还回 `scratchPool` 时，若缓冲区容量超过 4KB，则会放弃复用直接销毁。这防止了单次超长数据的处理导致缓冲池被撑大后永远占用内存。
- **并发安全**：所有池化操作均为并发安全，适用于高并发 Web 服务器。

## 推荐用法

```golang
package main

import (
    "fmt"
    "github.com/strugglerx/sjson"
)

func main() {
    data := map[string]interface{}{
        "nick":  "[天呐]",
        "extra": "{\"score\":100, \"tags\":[\"go\"]}",
    }

    // ⭐ 推荐：性能最佳，且输出经 json.Valid 严格验证
    result := sjson.StringWithJsonScanToString(data)
    fmt.Println(result)
}
```

## License

`sjson` is licensed under the [MIT License](LICENSE), 100% free and open-source, forever.
