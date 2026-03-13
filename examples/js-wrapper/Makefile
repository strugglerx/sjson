SHELL := /bin/zsh

GO ?= go
WASM_OUT := ./sjson.wasm
WASM_EXEC_SRC := $(shell $(GO) env GOROOT)/lib/wasm/wasm_exec.js
WASM_EXEC_OUT := ./wasm_exec.js
LOADER_SRC := ../wasm/loader.js
LOADER_OUT := ./loader.js

.PHONY: build serve clean

build:
	GOOS=js GOARCH=wasm $(GO) build -o $(WASM_OUT) ../wasm
	cp $(WASM_EXEC_SRC) $(WASM_EXEC_OUT)
	cp $(LOADER_SRC) $(LOADER_OUT)

serve:
	python3 -m http.server 4180

clean:
	rm -f $(WASM_OUT)
	rm -f $(WASM_EXEC_OUT)
	rm -f $(LOADER_OUT)
