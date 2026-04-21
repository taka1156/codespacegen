BINARY := codespacegen
E2E_TEST_DIR := ./e2e
CMD := ./cmd/codespacegen
BIN_DIR := ./bin

.PHONY: run build test clean e2e bin

run:
	go run $(CMD) $(RUN_ARGS)

build:
	mkdir -p $(BIN_DIR)
	go build -o $(BIN_DIR)/$(BINARY) $(CMD)

test:
	go test ./...

e2e:
	# UPD=--update is updating snapshots mode.
	rm -r $(E2E_TEST_DIR)/${BINARY} || true
	go build -o $(E2E_TEST_DIR)/${BINARY} $(CMD)
	bash $(E2E_TEST_DIR)/e2e.sh $(UPD)

bin:
	mkdir -p $(BIN_DIR)
	go build -o $(BIN_DIR)/$(BINARY) $(CMD)
	$(BIN_DIR)/codespacegen -output $(BIN_DIR)/.devcontainer

clean:
	rm -rf $(BIN_DIR)
