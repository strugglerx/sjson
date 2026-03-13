# js-wrapper Example

[中文说明](./README_ZH.md)

This example shows a frontend-friendly wrapper API on top of the WASM runtime.

## Build

```bash
make -C ./examples/js-wrapper build
```

## Run Locally

```bash
make -C ./examples/js-wrapper serve
```

Then open:

`http://127.0.0.1:4180/examples/js-wrapper/`

## API Shape

`wrapper.js` exposes:

- `createSJSONClient()`
- `client.expand(input)` -> `{ result, parsed }`
- `client.expandParsed(input)` -> parsed object
- `client.expandString(input)` -> expanded JSON string
