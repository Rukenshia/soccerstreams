BINARY_NAME := agent-py
SOURCES = $(wildcard *.py) $(wildcard **/*.py)

bin: $(BINARY_NAME)
.PHONY: bin

image: $(BINARY_NAME) Dockerfile version
	docker build -t gcr.io/soccerstreams-web/$(BINARY_NAME):`cat version` .
	gcloud docker -- push gcr.io/soccerstreams-web/$(BINARY_NAME):`cat version`
	@touch $@

$(BINARY_NAME): $(SOURCES)
	@touch $@

test:
	exit 0
.PHONY: test

lint:
	exit 0
.PHONY: lint