binaries = agent sweeper web

$(binaries): %: cmd/%
	cd cmd/$@ && make image

deps:
	go get -t ./... github.com/revel/cmd/revel github.com/alecthomas/gometalinter

test:
	go test ./...
.PHONY: test

lint:
	gometalinter --fast ./...
.PHONY: lint

secrets:
	kubectl create secret generic sentry-agent --from-file ./cmd/agent/sentry
	kubectl create secret generic sentry-sweeper --from-file ./cmd/sweeper/sentry
	kubectl create secret generic sentry-web --from-file ./cmd/web/sentry
	kubectl create secret generic graw --from-file ./graw
	kubectl create secret generic gcloud-service-account --from-file ./gcloud-service-account.json
	kubectl create secret generic dns-service-account --from-file ./dns-service-account.json

deploy:
	kubectl apply -f k8s/agent.yml
	kubectl apply -f k8s/sweeper.yml
	kubectl apply -f k8s/web.yml
	kubectl apply -f k8s/web-service.yml
.PHONY: secrets