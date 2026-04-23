BINARY := codespacegen
E2E_TEST_DIR := ./e2e
CMD := ./cmd/codespacegen
BIN_DIR := ./bin
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
LDFLAGS := -X main.version=$(VERSION)

DIST_TARGETS := \
	linux/amd64/tar.gz \
	linux/arm64/tar.gz \
	darwin/amd64/tar.gz \
	darwin/arm64/tar.gz \
	windows/amd64/exe

.PHONY: fmt run build test clean e2e bin dist

fmt:
	go fmt ./...

run:
	go run $(CMD) $(RUN_ARGS)

build:
	mkdir -p $(BIN_DIR)
	go build -ldflags="$(LDFLAGS)" -o $(BIN_DIR)/$(BINARY) $(CMD)

test:
	go test ./...

e2e:
	# UPD=--update is updating snapshots mode.
	rm -r $(E2E_TEST_DIR)/${BINARY} || true
	go build -ldflags="$(LDFLAGS)" -o $(E2E_TEST_DIR)/${BINARY} $(CMD)
	bash $(E2E_TEST_DIR)/e2e.sh $(UPD)

bin:
	mkdir -p $(BIN_DIR)
	go build -ldflags="$(LDFLAGS)" -o $(BIN_DIR)/$(BINARY) $(CMD)
	$(BIN_DIR)/codespacegen -output $(BIN_DIR)/.devcontainer

dist:
	@for target in $(DIST_TARGETS); do \
		GOOS=$$(echo $$target | cut -d/ -f1) \
		GOARCH=$$(echo $$target | cut -d/ -f2) \
		ARCHIVE=$$(echo $$target | cut -d/ -f3) \
		bash scripts/build.sh; \
	done

clean:
	rm -rf $(BIN_DIR) dist tmp
