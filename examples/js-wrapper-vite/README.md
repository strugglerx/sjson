# js-wrapper Vite Example

This example shows how to use `sjson` in a Vite frontend via a small wrapper module.

## Quick Start

```bash
cd examples/js-wrapper-vite
npm install
npm run dev
```

Open:

`http://127.0.0.1:5173`

## Scripts

- `npm run sync:wasm`: build and copy WASM runtime files into `public/wasm`
- `npm run dev`: sync runtime files and start Vite dev server
- `npm run build`: sync runtime files and build production assets
- `npm run preview`: preview production build

## Key Files

- `src/sjson-client.js`: wrapper API (`createSJSONClient`)
- `src/main.js`: demo page logic
- `index.html`: includes `/wasm/loader.js`
