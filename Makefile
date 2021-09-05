BIN_DIR = ${PWD}/bin

all: build

build: FORCE
	GOBIN=$(BIN_DIR) go install ./...

vet:
	go vet ./...

test:
	go test -cover -race ./...

clean:
	rm -rf $(BIN_DIR)

FORCE:

.PHONY: all vet test clean
