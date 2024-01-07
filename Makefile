PROJECT_PATH = $(shell pwd | sed 's/.*\(github.*\)/\1/')
GO := $(or $(GOROOT),/usr/lib/go)/bin/go
GO_BIN := $(or $(GOPATH),/usr/lib/go)/bin

all: welcome optimize_imports format linter testing building running

welcome:
	@echo "Make project: ($(PROJECT_PATH))"

optimize_imports:
	@echo "- Imports optimization"

	@/bin/find . -name '*.go' -exec ./scripts/clean_imports.sh {} \;
	@$(GO) install golang.org/x/tools/cmd/goimports@latest
	@$(GO_BIN)/goimports --local "$(PROJECT_PATH)" -w .

	@echo "	Done"

format:
	@echo "- Files formatting"

	@$(GO) install mvdan.cc/gofumpt@latest
	@$(GO_BIN)/gofumpt -w .

	@echo "	Done"
linter:
	@echo "- Linter check"

	@$(GO) install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@$(GO_BIN)/golangci-lint run

	@echo "	Done"
testing:
	@echo "- Testing"

	@$(GO) test ./...

	@echo "	Done"
building:
	@echo "- Building"

	@docker-compose -f build/docker-compose.yml build

	@echo "	Done"

running:
	@echo "- Running..."

	@docker-compose -f build/docker-compose.yml up
