GOCMD=go
GOFMT=$(GOCMD) fmt
GOBUILD=$(GOCMD) build
GOINSTALL=$(GOCMD) install

INFILES=duzzle.go 
OUTFILES=duzzle

all: build

fmt: $(INFILES)
	$(GOFMT) $(INFILES)

build: fmt $(INFILES)
	$(GOBUILD) $(INFILES)

install: $(INFILES)
	$(GOINSTALL) github.com/jeffjerseycow/duzzle

clean:
	$(RM) $(OUTFILES)
