.PHONY: test bench lint build clean

BINARY_DIR = bin
GO = go

test:
	$(GO) test ./...

bench:
	$(GO) test -bench=. ./...

lint:
	golangci-lint run ./...

build:
	mkdir -p $(BINARY_DIR)
	$(GO) build -o $(BINARY_DIR)/wc ./cmd/wc
	$(GO) build -o $(BINARY_DIR)/downloader ./cmd/downloader

clean:
	rm -rf $(BINARY_DIR)