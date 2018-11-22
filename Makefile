GOCMD=go
GOMOD=$(GOCMD) mod
GOBUILD=$(GOCMD) build
GOINSTALL=$(GOCMD) install

INFILES=duzzle.go
OUTFILES=go.mod \
	 duzzle \
	 go.sum 
INSTALLFILES=/usr/local/bin/duzzle

all: build

mod: 
	$(GOMOD) init github.com/jeffjerseycow/duzzle

build: $(INFILES) \
       mod
	$(GOBUILD)

install: $(OUTFILES)
	cp duzzle /usr/local/bin

uninstall: $(INSTALLFILES)
	${RM} $(INSTALLFILES)

clean:
	$(RM) $(OUTFILES)
