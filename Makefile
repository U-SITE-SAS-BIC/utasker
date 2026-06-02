BIN      := task
GIT_TAG  := $(shell git describe --tags --abbrev=0 2>/dev/null || echo "dev")
LDFLAGS  := -s -w -X main.version=$(GIT_TAG)

.PHONY: build test clean install lint release

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
	cp $(BIN) $(GOPATH)/bin/$(BIN) || cp $(BIN) /usr/local/bin/$(BIN)

release:
	goreleaser release --snapshot --clean
