

APPNAME=$(shell cat APPNAME)
BUILD=$(shell date -u +%Y%m%d%H%M%S)
VERSION=$(shell cat VERSION)
LDFLAGS="-X main.APPNAME=$(APPNAME) -X main.BUILD=$(BUILD) -X main.VERSION=$(VERSION)"

# needed for building static binaries which uses networking
export CGO_ENABLED=0

all: clean test build

build/$(APPNAME):
	@echo Building binaries ...
	@go build -ldflags $(LDFLAGS) -o build/$(APPNAME) platform/$(APPNAME)/main.go
	@echo Done.

build: build/$(APPNAME)

clean:
	@echo Cleaning up previous build ...
	@rm -f build/*
	@echo Done.

init:
	@echo Creating initial commit ...
	@rm -rf .git
	@git init
	@git add .
	@git commit -m "Initial commit"
	@git remote add origin git@github.com:andygeiss/$(APPNAME).git
	@git push -u --force origin master

install:
	@echo Installing $(APPNAME) ...
	@sudo cp build/$(APPNAME) /usr/local/bin/
	@echo Done.

test:
	@echo Testing ...
	@go test -v github.com/andygeiss/$(APPNAME)/...
	@echo Done.
