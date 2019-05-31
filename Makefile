ROOT:=$(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))
GAEA_OUT:=$(ROOT)/bin/gaea
GAEA_CC_OUT:=$(ROOT)/bin/gaea-cc
PKG:=$(shell go list -m)

.PHONY: all build gaea gaea-cc parser clean test build_with_coverage
all: build test

build: parser gaea gaea-cc

gaea:
	go build -o $(GAEA_OUT) $(shell bash gen_ldflags.sh $(GAEA_OUT) $(PKG)/core $(PKG)/cmd/gaea)

gaea-cc:
	go build -o $(GAEA_CC_OUT) $(shell bash gen_ldflags.sh $(GAEA_CC_OUT) $(PKG)/core $(PKG)/cmd/gaea-cc)

parser:
	cd parser && make && cd ..

clean:
	@rm -rf bin
	@rm -f .coverage.out .coverage.html

test:
	go test -coverprofile=.coverage.out ./...
	go tool cover -func=.coverage.out -o .coverage.func
	tail -1 .coverage.func
	go tool cover -html=.coverage.out -o .coverage.html

build_with_coverage:
	go test -c cmd/gaea/main.go cmd/gaea/main_test.go -coverpkg ./... -covermode=count -o bin/gaea
