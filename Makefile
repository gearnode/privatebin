all: build man

build: bin/privatebin

man: man/privatebin.1

clean:
	rm -rf bin
	rm -rf man

bin/privatebin: main.go go.sum go.mod privatebin/privatebin.go
	go build -o bin/privatebin

man/privatebin.1: doc/privatebin.1.md
	mkdir -p man
	pandoc --standalone --to man $> -o $@


.PHONY: all man build clean
