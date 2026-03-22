BINARY := upbit
BUILD_DIR := bin
INSTALL_DIR := $(HOME)/.local/bin
MODULE := github.com/kyungw00k/upbit
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
LDFLAGS := -ldflags "-X $(MODULE)/internal/cli.Version=$(VERSION)"

.PHONY: build install test lint clean

build:
	go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY) ./cmd/upbit

install: build
	mkdir -p $(INSTALL_DIR)
	cp $(BUILD_DIR)/$(BINARY) $(INSTALL_DIR)/$(BINARY)
	@echo "Installed to $(INSTALL_DIR)/$(BINARY)"
	@sh scripts/check-path.sh $(INSTALL_DIR)

test:
	go test ./... -v

lint:
	go vet ./...

clean:
	rm -rf $(BUILD_DIR)
