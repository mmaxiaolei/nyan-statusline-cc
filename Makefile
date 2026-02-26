BINARY_NAME := nyan-statusline
VERSION := 0.1.0
INSTALL_DIR := $(HOME)/.claude

.PHONY: build build-all install clean test lint

build:
	go build -ldflags="-s -w" -o $(BINARY_NAME) .

build-all:
	GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w" -o $(BINARY_NAME)-darwin-arm64 .
	GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o $(BINARY_NAME)-darwin-amd64 .

install: build
	mkdir -p $(INSTALL_DIR)
	cp $(BINARY_NAME) $(INSTALL_DIR)/$(BINARY_NAME)
	chmod +x $(INSTALL_DIR)/$(BINARY_NAME)
	@echo "Installed to $(INSTALL_DIR)/$(BINARY_NAME)"
	@echo "Add to ~/.claude/settings.json:"
	@echo '  "statusLine": { "type": "command", "command": "~/.claude/nyan-statusline", "padding": 0 }'

clean:
	rm -f $(BINARY_NAME) $(BINARY_NAME)-darwin-*

test:
	go test ./...

lint:
	go vet ./...
