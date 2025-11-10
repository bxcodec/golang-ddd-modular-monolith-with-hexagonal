# Database
POSTGRES_USER ?= user
POSTGRES_PASSWORD ?= password
POSTGRES_HOST ?= 127.0.0.1
POSTGRES_PORT ?= 5432
POSTGRES_DB ?= payment

# Exporting bin folder to the path for makefile
export PATH   := $(PWD)/bin:$(PATH)
# Default Shell
export SHELL  := bash
# Type of OS: Linux or Darwin.
export OSTYPE := $(shell uname -s | tr A-Z a-z)
export ARCH := $(shell uname -m)



# --- Tooling & Variables ----------------------------------------------------------------
include ./.misc/make/tools.Makefile
include ./.misc/make/help.Makefile

# ~~~ Development Environment ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

init: install-deps go-generate ## Initialize the project
up: dev-env dev-air             ## Startup / Spinup Docker Compose and air
down: docker-stop               ## Stop Docker
destroy: docker-teardown clean  ## Teardown (removes volumes, tmp files, etc...)

install-deps: migrate air gotestsum tparse mockery ## Install Development Dependencies (localy).
deps: $(MIGRATE) $(AIR) $(GOTESTSUM) $(TPARSE) $(MOCKERY) $(GOLANGCI) ## Checks for Global Development Dependencies.
deps:
	@echo "Required Tools Are Available"

dev-env: ## Bootstrap Environment (with a Docker compose help).
	@ docker compose up -d --build postgres

dev-env-test: dev-env ## Run application (within a Docker compose help)
	@ $(MAKE) image-build
	docker compose up web

dev-air: $(AIR) ## Starts AIR ( Continuous Development app).
	POSTGRES_HOST=$(POSTGRES_HOST) \
	POSTGRES_PORT=$(POSTGRES_PORT) \
	POSTGRES_USER=$(POSTGRES_USER) \
	POSTGRES_PASSWORD=$(POSTGRES_PASSWORD) \
	POSTGRES_DB=$(POSTGRES_DB) \
	air
 
docker-stop:
	@ docker compose down

docker-teardown:
	@ docker compose down --remove-orphans -v

# ~~~ Code Actions ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

lint: $(GOLANGCI) ## Runs golangci-lint with predefined configuration
	@echo "Applying linter"
	golangci-lint version
	golangci-lint run -c .golangci.yaml ./...

# -trimpath - will remove the filepathes from the reports, good to same money on network trafic,
#             focus on bug reports, and find issues fast.
# - race    - adds a racedetector, in case of racecondition, you can catch report with sentry.
#             https://golang.org/doc/articles/race_detector.html
#
# todo(butuzov): add additional flags to compiler to have an `version` flag.
build: ## Builds binary
	@ printf "Building aplication... "
	@ go build \
		-trimpath  \
		-o engine \
		./application/
	@ echo "done"


build-race: ## Builds binary (with -race flag)
	@ printf "Building aplication with race flag... "
	@ go build \
		-trimpath  \
		-race      \
		-o engine \
		./application/
	@ echo "done"


go-generate: $(MOCKERY) ## Runs mockery to generate mocks
	mockery


TESTS_ARGS := --format testname --jsonfile gotestsum.json.out
TESTS_ARGS += --max-fails 2
TESTS_ARGS += -- ./...
TESTS_ARGS += -test.parallel 2
TESTS_ARGS += -test.count    1
TESTS_ARGS += -test.failfast
TESTS_ARGS += -test.coverprofile   coverage.out
TESTS_ARGS += -test.timeout        5s
TESTS_ARGS += -race

tests: test-unit ## Run unit tests by default

test-unit: $(GOTESTSUM) ## Run only unit tests (excludes e2e)
	@ gotestsum $(TESTS_ARGS) -short

test-e2e: $(GOTESTSUM) ## Run only E2E tests
	@ gotestsum --format testname --jsonfile gotestsum-e2e.json.out -- -tags=e2e -test.parallel 1 -test.count 1 -test.timeout 5m -test.coverprofile coverage-e2e.out ./...

test-all: $(GOTESTSUM) ## Run all tests (unit + e2e)
	@ gotestsum --format testname --jsonfile gotestsum-all.json.out -- -tags=e2e -test.parallel 1 -test.count 1 -test.timeout 5m -test.coverprofile coverage-all.out -race ./...

tests-complete: test-all $(TPARSE) ## Run Tests & parse details
	@cat gotestsum-all.json.out | $(TPARSE) -all -notests

# ~~~ Docker Build ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

IMAGE_NAME ?= payment-app
IMAGE_TAG ?= latest

.ONESHELL:
image-build: ## Build Docker image
	@ echo "Building Docker image: $(IMAGE_NAME):$(IMAGE_TAG)"
	@ docker build \
		--file Dockerfile \
		--tag $(IMAGE_NAME):$(IMAGE_TAG) \
		--build-arg BUILDKIT_INLINE_CACHE=1 \
			.
	@ echo "Image built successfully: $(IMAGE_NAME):$(IMAGE_TAG)"

image-build-no-cache: ## Build Docker image without cache
	@ echo "Building Docker image without cache: $(IMAGE_NAME):$(IMAGE_TAG)"
	@ docker build \
		--no-cache \
		--file Dockerfile \
		--tag $(IMAGE_NAME):$(IMAGE_TAG) \
			.
	@ echo "Image built successfully: $(IMAGE_NAME):$(IMAGE_TAG)"

# ~~~ Database Migrations ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

POSTGRES_DSN := "postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(POSTGRES_HOST):$(POSTGRES_PORT)/$(POSTGRES_DB)?sslmode=disable"

.PHONY: migrate-up
migrate-up: $(MIGRATE) ## Apply all (or N up) migrations.
	@ read -p "How many migration you wants to perform (default value: [all]): " N; \
	migrate -database $(POSTGRES_DSN) -path=./migrations up $${N}

.PHONY: migrate-down
migrate-down: $(MIGRATE) ## Apply all (or N down) migrations.
	@ read -p "How many migration you wants to perform (default value: [all]): " N; \
	migrate -database $(POSTGRES_DSN) -path=./migrations down $${N}

.PHONY: migrate-drop
migrate-drop: $(MIGRATE) ## Drop everything inside the database.
	@ migrate -database $(POSTGRES_DSN) -path=./migrations drop

.PHONY: migrate-create
migrate-create: $(MIGRATE) ## Create a set of up/down migrations with a specified name.
	@ read -p "Please provide name for the migration: " Name; \
	migrate create -ext sql -dir ./migrations $${Name}

# ~~~ Cleans ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

clean: clean-artifacts clean-docker

clean-artifacts: ## Removes Artifacts (*.out)
	@printf "Cleanning artifacts... "
	@rm -f *.out
	@echo "done."


clean-docker: ## Removes dangling docker images
	@ docker image prune -f
