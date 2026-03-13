# 示例目录

这个目录提供了几种常见接入方式的最小示例。

## 包含的场景

- `cli`：打包成命令行工具
- `http-server`：通过 HTTP API 暴露能力
- `wasm`：通过 WebAssembly 在浏览器里调用
- `js-wrapper`：面向浏览器调用的 WASM 包装层
- `js-wrapper-vite`：基于 Vite 的前端包装层示例

## 什么时候用哪种

- `cli`：适合本地工具、CI、桌面宿主进程调用
- `http-server`：适合前端页面调用后端服务
- `wasm`：适合必须在浏览器本地处理数据的场景
- `js-wrapper`：适合前端项目想要更稳定、更清晰的调用 API（不直接碰全局导出细节）
- `js-wrapper-vite`：适合项目本身就是 Vite 技术栈，希望快速接入

## 快速开始

用根目录 `Makefile` 一次性构建：

```bash
make examples
```

也可以按场景单独构建。

CLI：

```bash
make -C ./examples/cli build
./bin/sjson-cli < input.json
```

HTTP 服务：

```bash
make -C ./examples/http-server build
./bin/sjson-http-server
```

示例请求：

```json
{
  "data": {
    "profile": "{\"name\":\"alice\",\"tags\":[\"go\",\"json\"]}",
    "items": "[{\"id\":1,\"ok\":true},{\"id\":2,\"ok\":false}]",
    "note": "[not real json]"
  }
}
```

示例返回：

```json
{
  "result": "{\"profile\":{\"name\":\"alice\",\"tags\":[\"go\",\"json\"]},\"items\":[{\"id\":1,\"ok\":true},{\"id\":2,\"ok\":false}],\"note\":\"[not real json]\"}",
  "parsed": {
    "profile": {"name":"alice","tags":["go","json"]},
    "items": [{"id":1,"ok":true},{"id":2,"ok":false}],
    "note": "[not real json]"
  }
}
```

这里可以很直观地看到：

- `profile` 被展开成对象
- `items` 被展开成数组
- `note` 仍然保持普通字符串

WASM：

```bash
make -C ./examples/wasm build
```

这个命令会同时生成：

- `sjson.wasm`
- `wasm_exec.js`

WASM 示例现在按更接近生产可维护的方式拆分：

- `main.go`：Go 导出层
- `loader.js`：运行时加载和就绪处理
- `index.html`：纯 UI 展示层

浏览器示例页会同时展示：

- WASM 的原始返回值
- 展开后的对象结果

JS 包装层：

```bash
make -C ./examples/js-wrapper build
```

包装层对前端暴露的 API：

```js
import { createSJSONClient } from "./wrapper.js";

const client = await createSJSONClient();
const { parsed, result } = client.expand(data);
```

目录说明：

- `index.html`：包装层示例页面
- `wrapper.js`：`createSJSONClient` 包装 API
- `loader.js`：构建时复制的运行时加载器
- 详细说明见：[js-wrapper/README.md](./js-wrapper/README.md)

Vite 包装层示例：

```bash
cd ./examples/js-wrapper-vite
npm install
npm run dev
```

详细文档：

- [js-wrapper-vite/README.md](./js-wrapper-vite/README.md)
- [js-wrapper-vite/README_ZH.md](./js-wrapper-vite/README_ZH.md)
