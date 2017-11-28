
APPNAME=$(shell cat APPNAME)
VERSION=$(shell cat VERSION)

# needed for building static binaries which uses networking
export CGO_ENABLED=0

all: clean test build

build/$(APPNAME):
	@echo Building binaries ...
	@go build -o build/alphago platform/alphago/main.go
	@go build -o build/betago platform/betago/main.go
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

test:
	@echo Testing ...
	@go test -v github.com/andygeiss/$(APPNAME)/...
	@echo Done.
