BINARY_NAME := web.tar.gz
SOURCES = $(wildcard *.go) $(wildcard **/*.go) $(wildcard **/*.html) $(wildcard assets/**/*) $(wildcard conf/*)

bin: $(BINARY_NAME)
.PHONY: bin

$(BINARY_NAME): $(SOURCES)
	CGO_ENABLED=0 GOOS=linux revel package github.com/Rukenshia/soccerstreams/cmd/web

test:
	go test ./...
.PHONY: test

lint:
	gometalinter --fast ./...
.PHONY: lint