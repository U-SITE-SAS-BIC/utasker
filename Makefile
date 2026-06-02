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
	cp $(BIN) $(GOPATH)/bin/$(BIN) || sudo cp $(BIN) /usr/local/bin/$(BIN)

release:
	goreleaser release --clean

snapshot:
	goreleaser release --snapshot --clean
