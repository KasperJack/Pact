SHELL := powershell.exe
.SHELLFLAGS := -NoProfile -Command





BIN_DIR=../bin
.DEFAULT_GOAL := build



help: ## List available targets
	@Select-String -Path '$(MAKEFILE_LIST)' -Pattern '^[a-zA-Z0-9_-]+:.*?## .*$$' | ForEach-Object { $$_.Line -match '^([a-zA-Z0-9_-]+):.*?## (.*)$$' | Out-Null; '{0,-15} {1}' -f $$matches[1], $$matches[2] }


build: ## Build binary for repo checker (DEFAULT)
	$$env:GOOS='windows'; $$env:GOARCH='amd64'; go build -o $(BIN_DIR)/chrp.exe .




.PHONY: build help