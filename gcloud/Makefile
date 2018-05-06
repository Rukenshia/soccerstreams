DEPLOYMENT_EXISTS := $(shell gcloud deployment-manager deployments list --filter='name = buckets' --simple-list | wc -l | awk '{print $$1}')
update: deploy
	gcloud deployment-manager deployments update buckets --config buckets.yml

deploy:
ifeq ($(DEPLOYMENT_EXISTS), 1)
	@echo "Deployment already exists"
	@exit 0
else
	gcloud deployment-manager deployments create buckets --config buckets.yml
endif

clean:
	gcloud deployment-manager deployments delete buckets

.PHONY: deploy clean update