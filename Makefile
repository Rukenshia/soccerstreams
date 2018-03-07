binaries = agent sweeper web

$(binaries): %: cmd/%
	cd cmd/$@ && make image

test:
	go test ./...
.PHONY: test

lint:
	gometalinter --fast ./...
.PHONY: lint