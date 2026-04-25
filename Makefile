BINARY := codespacegen
E2E_TEST_DIR := ./e2e
CMD := ./cmd/codespacegen
BIN_DIR := ./bin
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
LDFLAGS := -X codespacegen/internal/app.Version=$(VERSION)

DIST_TARGETS := \
	linux/amd64/tar.gz \
	linux/arm64/tar.gz \
	darwin/amd64/tar.gz \
	darwin/arm64/tar.gz \
	windows/amd64/exe

.PHONY: fmt run build test test-cover clean e2e bin exec dist

fmt:
	go fmt ./...

run:
	go run $(CMD) $(RUN_ARGS)

build:
	mkdir -p $(BIN_DIR)
	rm -r $(BIN_DIR)/${BINARY} || true
	go build -ldflags="$(LDFLAGS)" -o $(BIN_DIR)/$(BINARY) $(CMD)

test:
	go test ./...

test-cover:
	go test -cover ./... -coverprofile=cover.out
	go tool cover -html=cover.out -o cover.html

e2e:
	# UPD=--update is updating snapshots mode.
	go build -ldflags="$(LDFLAGS)" -o $(BIN_DIR)/$(BINARY) $(CMD)
	mkdir -p $(E2E_TEST_DIR)/devcontainer_config
	mkdir -p $(E2E_TEST_DIR)/codespacegen_config
	cp $(BIN_DIR)/$(BINARY) $(E2E_TEST_DIR)/devcontainer_config/$(BINARY)
	cp $(BIN_DIR)/$(BINARY) $(E2E_TEST_DIR)/codespacegen_config/$(BINARY)
	bash $(E2E_TEST_DIR)/devcontainer_config/devcontainer_config.test.sh $(UPD)
	bash $(E2E_TEST_DIR)/codespacegen_config/codespacegen_config.test.sh $(UPD)
	rm -f $(E2E_TEST_DIR)/devcontainer_config/$(BINARY)
	rm -f $(E2E_TEST_DIR)/codespacegen_config/$(BINARY)

bin:
	mkdir -p $(BIN_DIR)
	go build -ldflags="$(LDFLAGS)" -o $(BIN_DIR)/$(BINARY) $(CMD)

exec:
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
	rm -f cover.out cover.html
