BIN_DIR=../bin

.DEFAULT_GOAL := build

build: build-x64 build-x86

build-x64:
	powershell -Command "$$env:GOOS='windows'; $$env:GOARCH='amd64'; go build -o $(BIN_DIR)/runner-x64.exe ."

build-x86:
	powershell -Command "$$env:GOOS='windows'; $$env:GOARCH='386'; go build -o $(BIN_DIR)/runner-x86.exe ."

release: release-x64 release-x86

release-x64:
	powershell -Command "$$env:GOOS='windows'; $$env:GOARCH='amd64'; go build -ldflags='-s -w' -o $(BIN_DIR)/runner-x64.exe ."

release-x86:
	powershell -Command "$$env:GOOS='windows'; $$env:GOARCH='386'; go build -ldflags='-s -w' -o $(BIN_DIR)/runner-x86.exe ."

.PHONY: build build-x64 build-x86 release release-x64 release-x86