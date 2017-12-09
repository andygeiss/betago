# Binary settings
USER=andygeiss
APPNAME=$(shell cat APPNAME)
TS=$(shell date -u '+%Y/%m/%d %H:%M:%S') 
LDFLAGS="-s"

all: clean test build

build/$(APPNAME):
	@echo $(TS) Building binaries ...
	@go build -ldflags $(LDFLAGS) -o build/$(APPNAME) platform/$(APPNAME)/main.go
	@echo $(TS) Done.

build: build/$(APPNAME)

clean:
	@echo $(TS) Cleaning up previous build ...
	@rm -f build/*
	@echo $(TS) Done.

install:
	@echo $(TS) Installing $(APPNAME) ...
	@sudo cp build/$(APPNAME) /usr/local/bin/
	@echo $(TS) Done.

test:
	@echo $(TS) Testing ...
	@go test -v github.com/$(USER)/$(APPNAME)/...
	@echo $(TS) Done.
