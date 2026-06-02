BIN      := utasker
GIT_TAG  := $(shell git describe --tags --abbrev=0 2>/dev/null || echo "dev")
LDFLAGS  := -s -w -X main.version=$(GIT_TAG) -X main.commit=$(shell git rev-parse --short HEAD 2>/dev/null) -X main.date=$(shell date +%Y-%m-%d)

.PHONY: build test clean install lint release snapshot

build:
	go build -ldflags "$(LDFLAGS)" -o $(BIN) .

test:
	go test ./...

vet:
	go vet ./...

lint: vet

clean:
	rm -f $(BIN)

install: build
	@mkdir -p $(shell go env GOPATH)/bin
	cp $(BIN) $(shell go env GOPATH)/bin/$(BIN)
	@echo "✓ Installed to $$(go env GOPATH)/bin/$(BIN)"

sudo-install: build
	sudo cp $(BIN) /usr/local/bin/$(BIN)
	@echo "✓ Installed to /usr/local/bin/$(BIN)"

release:
	goreleaser release --clean

snapshot:
	goreleaser release --snapshot --clean
