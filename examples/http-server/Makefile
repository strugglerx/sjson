SHELL := /bin/zsh

GO ?= go
BIN_DIR := ../../bin
BIN := $(BIN_DIR)/sjson-http-server

.PHONY: build run clean

build:
	mkdir -p $(BIN_DIR)
	$(GO) build -o $(BIN) .

run:
	$(GO) run .

clean:
	rm -f $(BIN)
