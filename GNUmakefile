PREFIX = /usr/local
BINDIR = $(PREFIX)/bin
MANDIR = $(PREFIX)/share/man

MKDIR = mkdir -p
GO = go
PANDOC = pandoc
INSTALL = install
RM = rm -f

DATETIME = "Sep 8, 2025"
VERSION = 2.1.1
LDFLAGS = -ldflags "-X 'main.cliVersion=$(VERSION)'"

BIN = bin/privatebin

.PHONY: all build man install uninstall clean test test-fuzz vet

all: build man

vet:
	$(GO) vet ./...

build:
	$(GO) build $(LDFLAGS) -o $(BIN) cmd/privatebin/main.go cmd/privatebin/cfg.go

man:
	@$(MKDIR) man
	$(PANDOC) --standalone --to man -M footer=$(VERSION) -M date=$(DATETIME) doc/privatebin.1.md -o man/privatebin.1
	$(PANDOC) --standalone --to man -M footer=$(VERSION) -M date=$(DATETIME) doc/privatebin-create.1.md -o man/privatebin-create.1
	$(PANDOC) --standalone --to man -M footer=$(VERSION) -M date=$(DATETIME) doc/privatebin-show.1.md -o man/privatebin-show.1
	$(PANDOC) --standalone --to man -M footer=$(VERSION) -M date=$(DATETIME) doc/privatebin.conf.5.md -o man/privatebin.conf.5

install: build man
	$(INSTALL) -m 755 $(BIN) $(BINDIR)/privatebin
	$(INSTALL) -m 644 man/privatebin.1 $(MANDIR)/man1/privatebin.1
	$(INSTALL) -m 644 man/privatebin-create.1 $(MANDIR)/man1/privatebin-create.1
	$(INSTALL) -m 644 man/privatebin-show.1 $(MANDIR)/man1/privatebin-show.1
	$(INSTALL) -m 644 man/privatebin.conf.5 $(MANDIR)/man5/privatebin.conf.5

uninstall:
	$(RM) $(BINDIR)/privatebin
	$(RM) $(MANDIR)/man1/privatebin.1
	$(RM) $(MANDIR)/man5/privatebin.conf.5

clean:
	$(RM) -r bin man

test:
	$(GO) test -v -race ./...

test-fuzz:
	@for fuzz in $$($(GO) test -list=Fuzz); do \
		echo "Running $$fuzz"; \
		$(GO) test -fuzz=$$fuzz -fuzztime=30s || exit 1; \
	done
