# js-wrapper Vite 示例

[English README](./README.md)

这个示例展示了如何在 Vite 前端项目中，通过包装模块调用 `sjson` 的 WASM 运行时。

## 快速开始

```bash
cd examples/js-wrapper-vite
npm install
npm run dev
```

打开：

`http://127.0.0.1:5173`

## 脚本说明

- `npm run sync:wasm`：构建并复制 WASM 运行时文件到 `public/wasm`
- `npm run dev`：同步运行时文件并启动 Vite 开发服务
- `npm run build`：同步运行时文件并执行生产构建
- `npm run preview`：预览生产构建结果

## 关键文件

- `src/sjson-client.js`：包装 API（`createSJSONClient`）
- `src/main.js`：示例页面逻辑
- `index.html`：加载 `/wasm/loader.js`
