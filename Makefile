BIN=bin/privatebin

SRC = cmd/privatebin/main.go \
      privatebin.go \
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

deinstall:
	rm -f /usr/local/bin/privatebin
	rm -f /usr/local/man/man1/privatebin.1
	rm -f /usr/local/man/man5/privatebin.conf.5

clean:
	rm -rf bin
	rm -rf man

$(BIN): $(SRC)
	go build -o bin/privatebin cmd/privatebin/main.go

.PHONY: all man build clean
