# Examples

[中文说明](./README_ZH.md)

This directory contains minimal integration examples for different delivery styles.

## Included Scenarios

- `cli`: package `sjson` as a command-line tool
- `http-server`: expose `sjson` through a JSON HTTP API
- `wasm`: expose `sjson` to browser JavaScript via WebAssembly
- `js-wrapper`: browser-facing wrapper API around WASM runtime
- `js-wrapper-vite`: Vite-based frontend example using the wrapper API

## When To Use Which

- Use `cli` when you want a local binary, CI tool, or desktop-host subprocess.
- Use `http-server` when frontend code should call a backend service.
- Use `wasm` when data must stay in the browser and you accept the extra integration cost.
- Use `js-wrapper` when frontend code wants a cleaner API than calling `window.SJSON.expand` directly.
- Use `js-wrapper-vite` when your app is already on Vite and you want drop-in integration.

## Quick Start

With `Makefile`:

```bash
make examples
```

Or build one scenario at a time:

CLI:

```bash
make -C ./examples/cli build
./bin/sjson-cli < input.json
```

HTTP server:

```bash
make -C ./examples/http-server build
./bin/sjson-http-server
```

Example request:

```json
{
  "data": {
    "profile": "{\"name\":\"alice\",\"tags\":[\"go\",\"json\"]}",
    "items": "[{\"id\":1,\"ok\":true},{\"id\":2,\"ok\":false}]",
    "note": "[not real json]"
  }
}
```

Example response:

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

WASM build:

```bash
make -C ./examples/wasm build
```

This also copies `wasm_exec.js` from your local Go installation into `examples/wasm/`.

The WASM example is structured like a small production-style demo:

- `main.go`: Go export layer
- `loader.js`: runtime boot and readiness handling
- `index.html`: UI only

The demo page now shows both:

- the raw WASM return payload
- the parsed expanded object for easier visual checking

JS wrapper build:

```bash
make -C ./examples/js-wrapper build
```

The wrapper example exposes a stable frontend API:

```js
import { createSJSONClient } from "./wrapper.js";

const client = await createSJSONClient();
const { parsed, result } = client.expand(data);
```

Wrapper files:

- `index.html`: wrapper demo page
- `wrapper.js`: `createSJSONClient` API
- `loader.js`: copied runtime loader (generated on build)
- detailed notes: [js-wrapper/README.md](./js-wrapper/README.md)
- Chinese notes: [js-wrapper/README_ZH.md](./js-wrapper/README_ZH.md)

Vite wrapper example:

```bash
cd ./examples/js-wrapper-vite
npm install
npm run dev
```

Details:

- [js-wrapper-vite/README.md](./js-wrapper-vite/README.md)
- [js-wrapper-vite/README_ZH.md](./js-wrapper-vite/README_ZH.md)
