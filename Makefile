help: ## show this message
	@echo "All commands can be run on local machine as well as inside dev container."
	@echo ""
	@sed -nE 's/^ *([^[:blank:]]+)[[:blank:]]*:[^#]*##[[:blank:]]*(.+)/\1\n\2/p' $(MAKEFILE_LIST) | tr '\n' '\0' | xargs -0 -n 2 printf '%-25s%s\n'
.PHONY: help

.DEFAULT_GOAL := help

ifndef INSIDE_DEV_CONTAINER
  NOT_INSIDE_DEV_CONTAINER = 1
endif

test: ## run all tests
	@echo "+ $@"
	go test -race -count 1 -p 8 -parallel 8 -timeout 1m ./...
.PHONY: test

test-cover: ## run all tests with code coverage
	@echo "+ $@"
	go test -race -count 1 -p 8 -parallel 8 -timeout 1m -coverpkg ./... -coverprofile coverage.out ./...
.PHONY: test-cover

lint: compose-build ## run linter
	@echo "+ $@"
	$(call RUN_IN_DEV_CONTAINER, golangci-lint run)
.PHONY: lint

bash: compose-build ## run bash inside container for development
 ifndef INSIDE_DEV_CONTAINER
	@echo "+ $@"
	docker-compose run --rm service-a bash
 endif
.PHONY: bash

start: compose-build ## start service in docker-compose
 ifdef NOT_INSIDE_DEV_CONTAINER
	@echo "+ $@"
	docker-compose up
 endif
.PHONY: start

stop: ## stop docker-compose
 ifdef NOT_INSIDE_DEV_CONTAINER
	@echo "+ $@"
	docker-compose down
 endif
.PHONY: stop

compose-build: ## build docker-compose
 ifdef NOT_INSIDE_DEV_CONTAINER
	@echo "+ $@"
	COMPOSE_DOCKER_CLI_BUILD=1 DOCKER_BUILDKIT=1 docker-compose build
 endif
.PHONY: compose-build

build-service-a: ## build service-a binary
	@echo "+ $@"
	go build -v -o ./.bin/service-a ./service-a
.PHONY: build-service-a

build-service-b: ## build service-b binary
	@echo "+ $@"
	go build -v -o ./.bin/service-b ./service-b
.PHONY: build-service-b

# $(1) - command to run inside container
RUN_IN_DEV_CONTAINER = $(if $(NOT_INSIDE_DEV_CONTAINER), docker-compose run --rm --no-deps service-a) $(1)
