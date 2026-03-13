import "./style.css";
import { createSJSONClient } from "./sjson-client.js";

const app = document.querySelector("#app");
app.innerHTML = `
  <main class="shell">
    <section class="hero">
      <p class="badge">Vite + WASM Wrapper</p>
      <h1>sjson frontend wrapper demo</h1>
      <p class="lead">Use a stable JS API without touching global bridge details.</p>
    </section>
    <section class="grid">
      <article class="panel">
        <h2>Input</h2>
        <textarea id="input">{
  "profile": "{\\"name\\":\\"alice\\",\\"tags\\":[\\"go\\",\\"json\\"]}",
  "items": "[{\\"id\\":1,\\"ok\\":true},{\\"id\\":2,\\"ok\\":false}]",
  "note": "[not real json]"
}</textarea>
        <div class="actions">
          <button id="run" disabled>Expand via Wrapper</button>
          <span id="status">Loading runtime...</span>
        </div>
      </article>
      <article class="panel">
        <h2>Parsed Output</h2>
        <pre id="parsed"></pre>
      </article>
      <article class="panel">
        <h2>Result String</h2>
        <pre id="result"></pre>
      </article>
      <article class="panel">
        <h2>Raw Payload</h2>
        <pre id="raw"></pre>
      </article>
    </section>
  </main>
`;

const runButton = document.querySelector("#run");
const statusEl = document.querySelector("#status");
const parsedEl = document.querySelector("#parsed");
const resultEl = document.querySelector("#result");
const rawEl = document.querySelector("#raw");
const inputEl = document.querySelector("#input");

let client = null;

function setStatus(message, isError = false) {
  statusEl.textContent = message;
  statusEl.style.color = isError ? "#bf1b2c" : "";
}

function pretty(value) {
  return JSON.stringify(value, null, 2);
}

async function bootstrap() {
  try {
    client = await createSJSONClient();
    runButton.disabled = false;
    setStatus("Runtime ready.");
    runButton.click();
  } catch (err) {
    setStatus(`Runtime init failed: ${String(err)}`, true);
    rawEl.textContent = String(err);
  }
}

runButton.addEventListener("click", () => {
  if (!client) {
    setStatus("Client is not ready yet.", true);
    return;
  }

  try {
    const output = client.expand(inputEl.value);
    parsedEl.textContent = pretty(output.parsed);
    resultEl.textContent = pretty(JSON.parse(output.result));
    rawEl.textContent = pretty(output);
    setStatus("Expanded successfully.");
  } catch (err) {
    parsedEl.textContent = "No parsed output.";
    resultEl.textContent = "No result string.";
    rawEl.textContent = String(err);
    setStatus(`Expansion failed: ${String(err)}`, true);
  }
});

bootstrap();
