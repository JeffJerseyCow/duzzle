GOCMD=go
GOFMT=$(GOCMD)fmt
GOMOD=$(GOCMD) mod
GOBUILD=$(GOCMD) build
GOINSTALL=$(GOCMD) install

VERSION=$(shell git describe)

INFILES=duzzle.go
OUTFILES=duzzle

INSTALLFILES=/usr/local/bin/duzzle

all: fmt \
     build

fmt:
	$(GOFMT) -w .

mod:
	$(GOMOD) init

build: $(INFILES)
	$(GOBUILD) -ldflags="-X main.version=${VERSION}"

install: $(OUTFILES)
	cp duzzle /usr/local/bin

uninstall: $(INSTALLFILES)
	${RM} $(INSTALLFILES)

clean:
	$(RM) $(OUTFILES)
