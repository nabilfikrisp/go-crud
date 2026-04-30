ifneq ($(wildcard .env),)
include .env
export
else
$(warning WARNING: .env file not found! Using .env.example)
include .env.example
export
endif

BASE_STACK = docker compose -f docker-compose.yml

# HELP =================================================================================================================
# This will output the help for each task
# thanks to https://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
.PHONY: help

help: ## Display this help screen
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

test-integration: ### Run integration tests
	go test -v ./integration_test/*
.PHONY: test-integration

compose-up-db: ### Run docker compose db container in background
	$(BASE_STACK) up -d db
	@echo "DB running on localhost:5433"
.PHONY: compose-up-db

compose-stop-db: ### Stop the db container (keeps it for fast restart)
	$(BASE_STACK) stop db
.PHONY: compose-stop-db

migrate-create:  ### create new migration
	migrate create -ext sql -dir migrations '$(word 2,$(MAKECMDGOALS))'
.PHONY: migrate-create

migrate-up: ### migration up
	migrate -path migrations -database '$(PG_URL)' up
.PHONY: migrate-up

migrate-down: ### rollback the last migration
	migrate -path migrations -database '$(PG_URL)' down 1
.PHONY: migrate-down