.POSIX:
.SUFFIXES:

GO = go
RM = rm
GOFLAGS =
PREFIX = /usr/local
BINDIR = $(PREFIX)/bin

all: kimchi

kimchi:
	$(GO) build $(GOFLAGS) .

clean:
	$(RM) -f kimchi

install: all
	mkdir -p $(DESTDIR)$(BINDIR)
	cp -f kimchi $(DESTDIR)$(BINDIR)
