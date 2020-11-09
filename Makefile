.POSIX:
.SUFFIXES:

GO = go
RM = rm
GOFLAGS =
PREFIX = /usr/local
BINDIR = $(PREFIX)/bin
SYSCONFDIR = /etc

goflags = $(GOFLAGS) \
	-ldflags="-X 'main.configPath=$(SYSCONFDIR)/kimchi/config'"

all: kimchi

kimchi:
	$(GO) build $(goflags) .

clean:
	$(RM) -f kimchi

install: all
	mkdir -p $(DESTDIR)$(BINDIR)
	cp -f kimchi $(DESTDIR)$(BINDIR)
