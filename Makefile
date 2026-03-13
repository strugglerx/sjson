SHELL := /bin/zsh

GO ?= go

.PHONY: help test fuzz bench fmt examples clean-examples

help:
	@echo "Available targets:"
	@echo "  make test            - run unit tests"
	@echo "  make fuzz            - run a short fuzz session"
	@echo "  make bench           - run benchmarks with memory stats"
	@echo "  make fmt             - format Go files"
	@echo "  make examples        - build all example scenarios"
	@echo "  make clean-examples  - clean all example build outputs"

test:
	$(GO) test ./...

fuzz:
	$(GO) test -run=^$$ -fuzz=FuzzStringWithJsonScanToString -fuzztime=5s

bench:
	$(GO) test -bench . -benchmem

fmt:
	$(GO) fmt ./...

examples:
	$(MAKE) -C ./examples/cli build
	$(MAKE) -C ./examples/http-server build
	$(MAKE) -C ./examples/wasm build
	$(MAKE) -C ./examples/js-wrapper build

clean-examples:
	$(MAKE) -C ./examples/cli clean
	$(MAKE) -C ./examples/http-server clean
	$(MAKE) -C ./examples/wasm clean
	$(MAKE) -C ./examples/js-wrapper clean
	rm -rf ./bin
