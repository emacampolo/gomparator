### Tools
GOTOOLS_CHECK = golangci-lint

all: check_tools test fmt linter

### Tools & dependencies

check_tools:
	@# https://stackoverflow.com/a/25668869
	@echo "Found tools: $(foreach tool,$(GOTOOLS_CHECK),\
        $(if $(shell which $(tool)),$(tool),$(error "No $(tool) in PATH")))"

### Testing

test:
	go test ./...

test-cover:
	go test ./... -covermode=atomic -coverprofile=/tmp/coverage.out
	go tool cover -html=/tmp/coverage.out

### Formatting, linting, and vetting

fmt:
	go fmt ./...

metalinter:
	@echo "==> Running linter"
	golangci-lint run ./...

# To avoid unintended conflicts with file names, always add to .PHONY
# unless there is a reason not to.
# https://www.gnu.org/software/make/manual/html_node/Phony-Targets.html
.PHONY: check_tools test test-cover fmt linter
