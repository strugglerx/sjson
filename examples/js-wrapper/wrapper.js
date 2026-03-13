const DEFAULT_VERSION = "20260313";

function asJSONString(input) {
  if (typeof input === "string") {
    return input;
  }
  return JSON.stringify(input);
}

export async function createSJSONClient(options = {}) {
  if (typeof window.initSJSONWasm !== "function") {
    throw new Error("window.initSJSONWasm is unavailable. Did you load loader.js?");
  }

  const version = options.version || DEFAULT_VERSION;
  const wasmExecPath = options.wasmExecPath || `./wasm_exec.js?v=${version}`;
  const wasmPath = options.wasmPath || `./sjson.wasm?v=${version}`;

  await window.initSJSONWasm({
    wasmExecPath,
    wasmPath,
    timeoutMs: options.timeoutMs || 6000,
  });

  return {
    expand(input) {
      const payload = asJSONString(input);
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
