BIN_DIR = $(CURDIR)/bin

all: build

build: FORCE
	GOBIN=$(BIN_DIR) go install ./...

vet:
	go vet ./...

test:
	go test -cover -race ./...

clean:
	$(RM) -r $(BIN_DIR)

FORCE:

.PHONY: all vet test clean
