BIN_DIR=../bin

.DEFAULT_GOAL := build

build:
	go build -o $(BIN_DIR)/runner.exe .

release:
	go build -ldflags="-s -w" -o $(BIN_DIR)/runner.exe .

.PHONY: build release