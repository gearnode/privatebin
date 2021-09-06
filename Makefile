BIN=bin/privatebin

SRC = main.go \
      privatebin/privatebin.go \
      go.sum \
      go.mod

all: build man

build: $(BIN)

man:
	mkdir -p man
	pandoc --standalone --to man doc/privatebin.1.md -o man/privatebin.1
	pandoc --standalone --to man doc/privatebin.conf.5.md -o man/privatebin.conf.5

install: build man
	install -m 555 bin/privatebin /usr/local/bin/privatebin
	cp man/privatebin.1 /usr/local/man/man1/privatebin.1
	cp man/privatebin.conf.5 /usr/local/man/man5/privatebin.conf.5

clean:
	rm -rf bin
	rm -rf man

$(BIN): $(SRC)
	go build -o bin/privatebin

.PHONY: all man build clean
