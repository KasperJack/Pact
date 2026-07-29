#BIN_DIR=bin

#build-runner:
#	go build -o $(BIN_DIR)/runner.exe ./runner

#build-ci:
#	go build -o $(BIN_DIR)/ci.exe ./ci

#build-all:
#	make build-runner
#	make build-ci




# TODO: Keep sub Makefiles overrideable.
# Root Makefile controls shared paths (bin, dist, test output) for full workspace builds.
# root Makefile is the orchestrator


.PHONY: default
default:
	@echo Available targets:
	@echo make vet    - Run go vet on all modules
	@echo # TODO: Add build targets for runner


.PHONY: vet
vet:
	go vet ./core/... ./runner/... ./ci/... ./psbridge/...