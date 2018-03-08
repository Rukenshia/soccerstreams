BINARY_NAME := web
SOURCES = $(wildcard *.go) $(wildcard **/*.go) $(wildcard **/*.html) $(wildcard assets/**/*) $(wildcard conf/*) $(wildcard ../../pkg/**/*.go)

bin: $(BINARY_NAME).tar.gz
.PHONY: bin

$(BINARY_NAME).tar.gz: $(SOURCES)
	CGO_ENABLED=0 GOOS=linux revel package github.com/Rukenshia/soccerstreams/cmd/web

image: $(BINARY_NAME).tar.gz Dockerfile version
	docker build -t gcr.io/soccerstreams-web/$(BINARY_NAME):`cat version` .
	gcloud docker -- push gcr.io/soccerstreams-web/$(BINARY_NAME):`cat version`
	@touch $@

test:
	go test ./...
.PHONY: test

lint:
	gometalinter --fast ./...
.PHONY: lint
