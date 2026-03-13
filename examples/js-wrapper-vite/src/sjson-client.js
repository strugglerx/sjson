function normalizeInput(input) {
  if (typeof input === "string") {
    return input;
  }
  return JSON.stringify(input);
}

export async function createSJSONClient(options = {}) {
  if (typeof window.initSJSONWasm !== "function") {
    throw new Error("window.initSJSONWasm is unavailable. Ensure /wasm/loader.js is loaded.");
  }

  await window.initSJSONWasm({
    wasmExecPath: options.wasmExecPath || "/wasm/wasm_exec.js",
    wasmPath: options.wasmPath || "/wasm/sjson.wasm",
    timeoutMs: options.timeoutMs || 6000,
  });

  return {
    expand(input) {
      const payload = normalizeInput(input);
      const result = window.SJSON.expand(payload);
      if (result && result.error) {
        throw new Error(result.error);
      }
      return result;
    },
    expandParsed(input) {
      return this.expand(input).parsed;
    },
    expandString(input) {
      return this.expand(input).result;
    },
  };
}
