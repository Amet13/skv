.PHONY: build lint test ci all

build:
	bash ./build.sh host

lint:
	bash ./scripts/lint.sh

test:
	go test ./...

ci: build lint

all: build lint test


