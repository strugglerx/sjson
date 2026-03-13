(() => {
const DEFAULT_TIMEOUT_MS = 6000;

function ensureNamespace() {
  window.SJSON = window.SJSON || {};
  if (typeof window.SJSON.ready !== "boolean") {
    window.SJSON.ready = false;
  }
  return window.SJSON;
}

function loadScript(src) {
  return new Promise((resolve, reject) => {
    const existing = document.querySelector(`script[data-src="${src}"]`);
    if (existing) {
      if (existing.dataset.loaded === "true") {
        resolve();
        return;
      }
      existing.addEventListener("load", () => resolve(), { once: true });
      existing.addEventListener("error", () => reject(new Error(`failed to load ${src}`)), { once: true });
      return;
    }

    const script = document.createElement("script");
    script.src = src;
    script.async = true;
    script.dataset.src = src;
    script.addEventListener("load", () => {
      script.dataset.loaded = "true";
      resolve();
    }, { once: true });
    script.addEventListener("error", () => reject(new Error(`failed to load ${src}`)), { once: true });
    document.head.appendChild(script);
  });
}

async function instantiateWasm(go, wasmPath) {
  const response = await fetch(wasmPath);
  if (!response.ok) {
    throw new Error(`failed to fetch ${wasmPath}: HTTP ${response.status}`);
  }

  try {
    return await WebAssembly.instantiateStreaming(response, go.importObject);
  } catch (_) {
    const fallbackResponse = await fetch(wasmPath);
    if (!fallbackResponse.ok) {
      throw new Error(`fallback fetch failed for ${wasmPath}: HTTP ${fallbackResponse.status}`);
    }
    const bytes = await fallbackResponse.arrayBuffer();
    return await WebAssembly.instantiate(bytes, go.importObject);
  }
}

function waitForReady(timeoutMs = DEFAULT_TIMEOUT_MS) {
  return new Promise((resolve, reject) => {
    const started = Date.now();

    const tick = () => {
      if (window.SJSON && window.SJSON.ready === true && typeof window.SJSON.expand === "function") {
        resolve(window.SJSON);
        return;
      }

      if (Date.now() - started > timeoutMs) {
        reject(new Error("window.SJSON.expand was not registered in time."));
        return;
      }

      window.setTimeout(tick, 30);
    };

    tick();
  });
}

async function initSJSONWasm(options = {}) {
  const wasmExecPath = options.wasmExecPath || "./wasm_exec.js";
  const wasmPath = options.wasmPath || "./sjson.wasm";
  const timeoutMs = options.timeoutMs || DEFAULT_TIMEOUT_MS;

  ensureNamespace();
  await loadScript(wasmExecPath);

  if (typeof Go !== "function") {
    throw new Error("wasm_exec.js loaded, but Go runtime is unavailable.");
  }

  const go = new Go();
  const readyPromise = waitForReady(timeoutMs);
  const { instance } = await instantiateWasm(go, wasmPath);

  const runPromise = go.run(instance);
  if (runPromise && typeof runPromise.catch === "function") {
    runPromise.catch((err) => {
      console.error("go.run failed", err);
    });
  }

  return readyPromise;
}

window.initSJSONWasm = initSJSONWasm;
})();
