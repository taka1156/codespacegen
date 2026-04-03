BINARY := codespacegen
CMD := ./cmd/codespacegen
BIN_DIR := ./bin

.PHONY: run build test clean

run:
	go run $(CMD) $(RUN_ARGS)

build:
	mkdir -p $(BIN_DIR)
	go build -o $(BIN_DIR)/$(BINARY) $(CMD)

test:
	go test ./...

clean:
	rm -rf $(BIN_DIR)
