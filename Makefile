.PHONY: build lint test ci all

build:
	./build.sh host

lint:
	./scripts/lint.sh

test:
	go test ./...

ci: build lint

all: build lint test


