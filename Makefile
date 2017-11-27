
APPNAME=$(shell cat APPNAME)
VERSION=$(shell cat VERSION)

# needed for building static binaries which uses networking
export CGO_ENABLED=0

all: clean test build

build/$(APPNAME):
	@echo Building binaries ...
	@go build -o $@ platform/$(APPNAME)/main.go
	@echo Done.

build: build/$(APPNAME)

clean:
	@echo Cleaning up previous build ...
	@rm -f build/$(APPNAME)
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
	@go test github.com/andygeiss/$(APPNAME)/...
	@echo Done.