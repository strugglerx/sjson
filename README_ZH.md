# sjson

[![Go Doc](https://godoc.org/github.com/strugglerx/sjson?status.svg)](https://godoc.org/github.com/strugglerx/sjson)
[![Production Ready](https://img.shields.io/badge/production-ready-blue.svg)](https://github.com/strugglerx/sjson)
[![License](https://img.shields.io/github/license/strugglerx/sjson.svg?style=flat)](https://github.com/strugglerx/sjson)

[English README](./README.md)

> 将结构体、`map` 或任意 Go 值中“以字符串形式存储的 JSON 对象 / 数组”展开成真正的 JSON 输出。

`sjson` 适合处理这类数据：字段本质上是 JSON，但上游却把它存成了字符串。这个包会先按 `encoding/json` 对 Go 值做标准编码，再扫描编码后的结果，把其中实际是合法 JSON 对象或数组的字符串值展开。

## 默认推荐用法

如果你不确定该用哪个接口，默认推荐使用 `StringWithJsonScanToString(v)`。

- 它是这个包当前主推的高性能扫描路径。
- 返回值是最终 JSON `string`，最符合大多数调用场景。
- 如果你不想使用 panic 风格接口，可以改用 `StringWithJsonScanToStringE(v)`。

## 典型应用场景

- MySQL、PostgreSQL、ES 等系统里，字段内容实际是 JSON，但存储层给的是字符串。
- 第三方接口返回 JSON 字符串嵌在 JSON 结构里。
- 日志、埋点、审计事件、消息队列载荷里有内嵌 JSON 文本。
- 后台管理、调试工具、导出工具，希望直接得到更可读的最终 JSON。

如果你想看不同打包和接入方式的最小示例，可以直接看 [examples/README.md](./examples/README.md)。
目前示例已经包含 `cli`、`http-server`、`wasm`、`js-wrapper`、`js-wrapper-vite` 五种场景。

一个更直观的判断方式是：

- `profile: "{\"name\":\"alice\"}"` 这种值本身是合法 JSON 对象字符串，会展开
- `items: "[{\"id\":1}]"` 这种值本身是合法 JSON 数组字符串，也会展开
- `note: "[not real json]"` 这种只是长得像 JSON，但并不是合法 JSON，会保持原字符串

## 它解决什么问题

输入：

```json
{
  "profile": "{\"name\":\"alice\",\"tags\":[\"go\",\"json\"]}",
  "events": "[{\"type\":\"login\",\"ok\":true}]",
  "note": "[not real json]"
}
```

输出：

```json
{
  "profile": {"name":"alice","tags":["go","json"]},
  "events": [{"type":"login","ok":true}],
  "note": "[not real json]"
}
```

真正的内嵌 JSON 会被展开，普通字符串保持不变。

## 安装

```bash
go get -u -v github.com/strugglerx/sjson
```

仓库里也附带了一个 `Makefile`，本地开发时可以直接用：

```bash
make test
make examples
```

## API 说明

| 函数 | 说明 | 备注 |
|:---|:---|:---|
| `ToJsonByte(v)` | 标准 JSON 编码，返回 `[]byte` | Must 风格，编码失败会 panic |
| `ToJsonByteE(v)` | 标准 JSON 编码，返回 `[]byte` | 返回 `error` |
| `ToJsonString(v)` | 标准 JSON 编码，返回 `string` | Must 风格，编码失败会 panic |
| `ToJsonStringE(v)` | 标准 JSON 编码，返回 `string` | 返回 `error` |
| `StringWithJsonScanToString(v)` | 展开内嵌 JSON 字符串，返回 `string` | 默认推荐 |
| `StringWithJsonScanToStringE(v)` | 同上，但返回 `error` | 适合不希望 panic 的调用方 |
| `StringWithJsonScanToBytes(v)` | 展开内嵌 JSON 字符串，返回 `[]byte` | 下游本来就处理字节流时更方便 |
| `StringWithJsonScanToBytesE(v)` | 同上，但返回 `error` | 稳定字节版接口 |
| `StringWithJsonSafetyRegexToString(v)` | 旧版正则实现 | 更慢，主要保留用于对比 |

## 使用示例

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

如果你不希望使用 panic 风格接口：

```go
out, err := sjson.StringWithJsonScanToStringE(data)
if err != nil {
	panic(err)
}
fmt.Println(out)
```

如果你的下游逻辑本来就是处理 `[]byte`：

```go
b, err := sjson.StringWithJsonScanToBytesE(data)
if err != nil {
	panic(err)
}
_ = b
```

## 行为边界

- 只有“字符串值本身能解析成合法 JSON 对象或数组”时，才会被展开。
- 类似 `"[hello]"`、`"{not-json}"` 这种看起来像 JSON、但实际不是合法 JSON 的内容，会保持原字符串。
- 判断依据是 JSON 编码后的结果，不是原始 Go 字符串的裸字节。
- 当前行为是对最终编码结果做一层展开，不会无限递归反复解码。
- 编码语义遵循 Go 标准库 `encoding/json`，包括 UTF-8 规范化等行为。

## 性能表现

当前本地 benchmark 环境：

- 机器：Apple M1 Pro
- Go 模块版本：`go 1.16`
- 数据量：约 `1MB`

最新 benchmark 快照：

| Benchmark | 耗时 | 内存 | 分配次数 |
|:---|:---|:---|:---|
| `SafetyRegex_1MB` | `~115.7 ms/op` | `~133.7 MB/op` | `23656 allocs/op` |
| `Scanner_1MB` | `~3.62 ms/op` | `~1.57 MB/op` | `13766 allocs/op` |
| `SafetyRegex_10MB` | `~8.91 s/op` | `~12.25 GB/op` | `237446 allocs/op` |
| `Scanner_10MB` | `~35.8 ms/op` | `~20.2 MB/op` | `137543 allocs/op` |
| `Scanner_100MB` | `~363 ms/op` | `~308.6 MB/op` | `1375078 allocs/op` |

大致可以理解为：

- 相比旧版正则路径，速度提升约 `32x`
- 内存占用下降到原来的 `1%` 到 `2%`
- 到了更大的数据量，扫描版依然有很强的可用性

## 为什么快

- 用单次字节扫描替代正则重写
- 使用 `sync.Pool` 复用临时缓冲
- 热路径采用分段拷贝，减少逐字节写入
- 只有候选值看起来像 JSON 容器时才做进一步校验

## 大数据结论

- 对普通后端接口、内部平台和数据处理脚本来说，扫描版的性能已经足够放心使用。
- 在 `10MB` 量级下，扫描版仍然保持在 `~35ms` 级别，而旧正则版已经进入多秒区间。
- 在 `100MB` 量级下，扫描版在测试机器上依然能在亚秒级完成。
- 旧正则版在 `100MB` 量级已经不再适合作为实际生产路径，真正应该默认使用的是扫描版。

可以把结论理解成一句话：

- 如果你的数据里经常混着大对象或脏 JSON，默认直接用 `StringWithJsonScanToString` 就对了。

## 健壮性与测试

目前仓库里已经包含：

- 常规单元测试
- 更复杂的手工边界测试
- benchmark
- Go fuzz test

最近一轮 fuzz 已在扫描路径上跑通，在测试时间窗口内没有打出 panic，也没有生成非法 JSON。

## 并发与内存行为

- 所有池化复用逻辑都是并发安全的
- 大容量缓冲不会无限制回收到池中，避免偶发超大请求把池子永久撑胖
- `sync.Pool` 对服务端场景比较友好，能减轻 GC 压力

## 适合在什么场景使用

适合：

- 你最终需要的是 JSON 字符串或字节输出
- 你的数据里经常混着“字符串形式的 JSON”
- 你想快速把输出变得可读，而不想先写一套复杂结构体来反序列化

不太适合：

- 你需要对任意 JSON 文档做深层、递归、语义级别的转换
- 你需要强类型、schema-aware 的完整解码流程
- 你的输入本来就是类型干净、结构正确的 Go 对象

## License

`sjson` 使用 [MIT License](LICENSE)。
