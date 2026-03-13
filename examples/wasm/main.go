//go:build js && wasm
// +build js,wasm

package main

import (
	"encoding/json"
	"syscall/js"

	"github.com/strugglerx/sjson"
)

func main() {
	api := js.Global().Get("SJSON")
	if api.IsUndefined() || api.IsNull() {
		api = js.Global().Get("Object").New()
		js.Global().Set("SJSON", api)
	}
	api.Set("expand", js.FuncOf(expandJSON))
	api.Set("ready", true)
	js.Global().Set("SJSON", api)
	select {}
}

func expandJSON(_ js.Value, args []js.Value) interface{} {
	if len(args) == 0 {
		return map[string]interface{}{
			"error": "missing input",
		}
	}

	input := args[0].String()
	result, err := sjson.StringWithJsonScanToStringE(json.RawMessage(input))
	if err != nil {
		return map[string]interface{}{
			"error": err.Error(),
		}
	}

	return map[string]interface{}{
		"result": result,
		"parsed": parseExpanded(result),
	}
}

func parseExpanded(result string) interface{} {
	var parsed interface{}
	if err := json.Unmarshal([]byte(result), &parsed); err != nil {
		return map[string]interface{}{
			"error": err.Error(),
		}
	}
	return parsed
}
