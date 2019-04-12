### Tools
GOTOOLS_CHECK = dep gin golangci-lint

all: check_tools fmt ensure-deps test linter

### Tools & dependencies

check_tools:
	@# https://stackoverflow.com/a/25668869
	@echo "Found tools: $(foreach tool,$(GOTOOLS_CHECK),\
        $(if $(shell which $(tool)),$(tool),$(error "No $(tool) in PATH")))"

### Testing

test:
	go test ./... -covermode=atomic -coverpkg=./... -count=1 -race

test-cover:
	go test ./... -covermode=atomic -coverprofile=/tmp/coverage.out -coverpkg=./... -count=1
	go tool cover -html=/tmp/coverage.out

### Formatting, linting, and deps

fmt:
	go fmt ./...

linter:
	@echo "==> Running linter"
	golangci-lint run ./...

ensure-deps:
	@echo "==> Running dep ensure"
	dep ensure

run:
	gin --port 8080 run main.go

# To avoid unintended conflicts with file names, always add to .PHONY
# unless there is a reason not to.
# https://www.gnu.org/software/make/manual/html_node/Phony-Targets.html
.PHONY: check_tools test test-cover fmt linter ensure-deps run 
