BIN_DIR=bin

build-runner:
	go build -o $(BIN_DIR)/runner.exe ./runner

build-ci:
	go build -o $(BIN_DIR)/ci.exe ./ci

build-all:
	make build-runner
	make build-ci