# Print help
@help:
    just --list

# Build and run
@run *args:
    go run cmd/main.go {{args}}

# Build
@build:
    version="0.0.0-$(git rev-parse --short HEAD)"; \
    CGO_ENABLED=0 go build -o bin/cuemix \
      -ldflags="-s -w -X 'github.com/amir-ahmad/cuemix/cmd/cuemix.version=$version'" \
      cmd/main.go

# Run tests
@test:
    GOFLAGS="-count=1" gotestsum -f testname ./...

# Test a path
@test-path *args:
    go test -count=1 {{args}}

# Format code with go fmt
@fmt:
    go fmt ./...
