SHELL := /bin/zsh

GO ?= go
WASM_OUT := ./sjson.wasm
WASM_EXEC_SRC := $(shell $(GO) env GOROOT)/lib/wasm/wasm_exec.js
WASM_EXEC_OUT := ./wasm_exec.js

.PHONY: build clean

build:
	GOOS=js GOARCH=wasm $(GO) build -o $(WASM_OUT) .
	cp $(WASM_EXEC_SRC) $(WASM_EXEC_OUT)

clean:
	rm -f $(WASM_OUT)
	rm -f $(WASM_EXEC_OUT)
