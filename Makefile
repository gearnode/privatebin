MKDIR=	mkdir -p
GO=	go
PANDOC=	pandoc
CP=	cp
RM=	rm
INSTALL=	install

DATETIME=	"Sep 28, 2022"
VERSION=	1.2.0
LDFLAGS=	-ldflags "-X 'gearno.de/privatebin/internal/version.Version=$(VERSION)'"

BIN=bin/privatebin

SRC = cmd/privatebin/main.go \
      privatebin.go \
      go.sum \
      go.mod

all: build man

build: $(BIN)

man:
	@$(MKDIR) man
	$(PANDOC) \
		--standalone \
		--to man \
		-M footer=$(VERSION) \
		-M date=$(DATETIME) \
		doc/privatebin.1.md \
		-o man/privatebin.1
	$(PANDOC) \
		--standalone \
		--to man \
		-M footer=$(VERSION) \
		-M date=$(DATETIME) \
		doc/privatebin.conf.5.md \
		-o man/privatebin.conf.5

install: build man
	$(INSTALL) -m 555 bin/privatebin /usr/local/bin/privatebin
	$(CP) man/privatebin.1 /usr/local/man/man1/privatebin.1
	$(CP) man/privatebin.conf.5 /usr/local/man/man5/privatebin.conf.5

deinstall:
	$(RM) -f /usr/local/bin/privatebin
	$(RM) -f /usr/local/man/man1/privatebin.1
	$(RM) -f /usr/local/man/man5/privatebin.conf.5

clean:
	$(RM) -rf bin
	$(RM) -rf man

$(BIN): $(SRC)
	$(GO) build $(LDFLAGS) -o $(BIN) cmd/privatebin/main.go

.PHONY: all man build clean install deinstall
