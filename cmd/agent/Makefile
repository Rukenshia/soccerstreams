BINARY_NAME := agent
SOURCES = $(wildcard *.go) $(wildcard **/*.go) $(wildcard ../../pkg/**/*.go)

bin: $(BINARY_NAME)
.PHONY: bin

image: $(BINARY_NAME) Dockerfile version
	docker build -t gcr.io/soccerstreams-web/$(BINARY_NAME):`cat version` .
	gcloud docker -- push gcr.io/soccerstreams-web/$(BINARY_NAME):`cat version`
	@touch $@

$(BINARY_NAME): $(SOURCES)
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags '-w -s' -a -installsuffix cgo

test:
	go test ./...
.PHONY: test

lint:
	gometalinter --fast ./...
.PHONY: lint