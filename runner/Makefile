SHELL := powershell.exe
.SHELLFLAGS := -NoProfile -Command

BIN_DIR=../bin
.DEFAULT_GOAL := build



help: ## List available targets
	@Select-String -Path '$(MAKEFILE_LIST)' -Pattern '^[a-zA-Z0-9_-]+:.*?## .*$$' | ForEach-Object { $$_.Line -match '^([a-zA-Z0-9_-]+):.*?## (.*)$$' | Out-Null; '{0,-15} {1}' -f $$matches[1], $$matches[2] }



build: build-x64 build-x86 ## Build debug binaries for x64 and x86 (DEFAULT)


build-x64: ## Build debug binary for Windows x64
	$$env:GOOS='windows'; $$env:GOARCH='amd64'; go build -o $(BIN_DIR)/pact64.exe .


build-x86: ## Build debug binary for Windows x86
	$$env:GOOS='windows'; $$env:GOARCH='386'; go build -o $(BIN_DIR)/pact86.exe .





release: release-x64 release-x86 ## Build stripped release binaries for x64 and x86


release-x64: ## Build stripped release binary for Windows x64
	$$env:GOOS='windows'; $$env:GOARCH='amd64'; go build -ldflags='-s -w' -o $(BIN_DIR)/pact64.exe .


release-x86: ## Build stripped release binary for Windows x86
	$$env:GOOS='windows'; $$env:GOARCH='386'; go build -ldflags='-s -w' -o $(BIN_DIR)/pact86.exe .




.PHONY: build build-x64 build-x86 release release-x64 release-x86 help