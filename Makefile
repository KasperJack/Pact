BIN_DIR=bin

build-pactci:
	go build -o $(BIN_DIR)/pact-ci.exe ./cmd/ci

build-pcatbuild:
	go build -o $(BIN_DIR)/pcat-check.exe ./cmd/check

build-all:
	make build-pactci
	make build-pcatbuild