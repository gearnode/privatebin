PREFIX = /usr/local
BINDIR = $(PREFIX)/bin
MANDIR = $(PREFIX)/share/man

MKDIR = mkdir -p
GO = go
PANDOC = pandoc
INSTALL = install
RM = rm -f

DATETIME = "Jan 08, 2023"
VERSION = 1.4.0
LDFLAGS = -ldflags "-X 'gearno.de/privatebin/internal/version.Version=$(VERSION)'"

BIN = bin/privatebin
SRC = cmd/privatebin/main.go privatebin.go go.sum go.mod

.PHONY: all build man install uninstall clean

all: build man

build: $(BIN)

$(BIN): $(SRC)
	$(GO) build $(LDFLAGS) -o $@ cmd/privatebin/main.go

man:
	@$(MKDIR) man
	$(PANDOC) --standalone --to man -M footer=$(VERSION) -M date=$(DATETIME) doc/privatebin.1.md -o man/privatebin.1
	$(PANDOC) --standalone --to man -M footer=$(VERSION) -M date=$(DATETIME) doc/privatebin.conf.5.md -o man/privatebin.conf.5

install: build man
	$(INSTALL) -m 755 $(BIN) $(BINDIR)/privatebin
	$(INSTALL) -m 644 man/privatebin.1 $(MANDIR)/man1/privatebin.1
	$(INSTALL) -m 644 man/privatebin.conf.5 $(MANDIR)/man5/privatebin.conf.5

uninstall:
	$(RM) $(BINDIR)/privatebin
	$(RM) $(MANDIR)/man1/privatebin.1
	$(RM) $(MANDIR)/man5/privatebin.conf.5

clean:
	$(RM) -r bin man
