# js-wrapper 示例

[English README](./README.md)

这个示例展示了如何在 WASM 运行时之上再包一层前端友好的 API。

## 构建

```bash
make -C ./examples/js-wrapper build
```

## 本地运行

```bash
make -C ./examples/js-wrapper serve
```

然后打开：

`http://127.0.0.1:4180/examples/js-wrapper/`

## API 形态

`wrapper.js` 暴露：

- `createSJSONClient()`
- `client.expand(input)` -> `{ result, parsed }`
- `client.expandParsed(input)` -> 展开后的对象
- `client.expandString(input)` -> 展开后的 JSON 字符串
