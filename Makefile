help: ## show this message
	@echo "All commands can be run on local machine as well as inside dev container."
	@echo ""
	@sed -nE 's/^ *([^[:blank:]]+)[[:blank:]]*:[^#]*##[[:blank:]]*(.+)/\1\n\2/p' $(MAKEFILE_LIST) | tr '\n' '\0' | xargs -0 -n 2 printf '%-25s%s\n'
.PHONY: help

.DEFAULT_GOAL := help

ifndef INSIDE_DEV_CONTAINER
  NOT_INSIDE_DEV_CONTAINER = 1
endif

gen: compose-build ## run code generation
	@echo "+ $@"
	$(call RUN_IN_DEV_CONTAINER, make _gen)
.PHONY: lint

_gen:
	sqlc generate --file service-a/postgres/sqlc.yml
	sqlc generate --file service-b/postgres/sqlc.yml
	go generate ./...
.PHONY: _gen

lint: compose-build ## run linter
	@echo "+ $@"
	$(call RUN_IN_DEV_CONTAINER, golangci-lint run)
.PHONY: lint

test: ## run all tests
	@echo "+ $@"
	go test -race -count 1 -p 8 -parallel 8 -timeout 1m ./...
.PHONY: test

test-cover: ## run all tests with code coverage
	@echo "+ $@"
	go test -race -count 1 -p 8 -parallel 8 -timeout 1m -coverpkg ./... -coverprofile coverage.out ./...
.PHONY: test-cover

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

migrate-service-a: compose-build ## apply all migrations to service-a
	@echo "+ $@"
	$(call RUN_IN_DEV_CONTAINER, ./scripts/migrate.sh service-a)
.PHONY: migrate-service-a

migrate-service-b: compose-build ## apply all migrations to service-b
	@echo "+ $@"
	$(call RUN_IN_DEV_CONTAINER, ./scripts/migrate.sh service-b)
.PHONY: migrate-service-b

check-tidy: ## ensure go.mod is tidy
	@echo "+ $@"
	cp go.mod go.check.mod
	cp go.sum go.check.sum
	go mod tidy -modfile=go.check.mod
	diff -u go.mod go.check.mod
	diff -u go.sum go.check.sum
	rm go.check.mod go.check.sum
.PHONY: check-tidy

check-vendor: ## ensure vendor is up-to-date
	@echo "+ $@"
	@$(call ENSURE_NO_CHANGES, check-vendor is unavailable: there are uncommitted changes in git)
	@go mod vendor
	@$(call ENSURE_NO_CHANGES, vendor is outdated: please run 'go mod vendor' before commit)
	@echo OK
.PHONY: check-vendor

check-gen: ## ensure generated code is up-to-date
	@echo "+ $@"
	@$(call ENSURE_NO_CHANGES, check-gen is unavailable: there are uncommitted changes in git)
	@$(MAKE) gen
	@$(call ENSURE_NO_CHANGES, generated code is outdated: please run 'make gen' before commit)
	@echo OK
.PHONY: check-gen

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

# $(1) - error message
ENSURE_NO_CHANGES = test -z "`git status --porcelain`" || (git diff; echo "\nERROR:$(1)\n"; exit 1)

# $(1) - command to run inside container
RUN_IN_DEV_CONTAINER = $(if $(NOT_INSIDE_DEV_CONTAINER), docker-compose run --rm --no-deps service-a) $(1)
