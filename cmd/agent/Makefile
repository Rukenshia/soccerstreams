BINARY_NAME := agent
SOURCES = $(wildcard *.go) $(wildcard **/*.go)

bin: $(BINARY_NAME)
.PHONY: bin

image: $(BINARY_NAME) Dockerfile
	docker build -t soccerstreams-$(BINARY_NAME) .
	@touch $@

$(BINARY_NAME): $(SOURCES)
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo

test:
	go test ./...
.PHONY: test

lint:
	gometalinter --fast ./...
.PHONY: lint