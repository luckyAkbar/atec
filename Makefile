SWAG_VERSION := v1.16.4
SWAG := $(shell which swag)
CURRENT_SWAG_VERSION := $(shell swag --version 2>/dev/null | grep -o 'v[0-9]\+\.[0-9]\+\.[0-9]\+' || echo "")

install-swag:
ifeq ($(SWAG),)
	@echo "Swag is not installed. Installing swag version $(SWAG_VERSION)..."
	go install github.com/swaggo/swag/cmd/swag@$(SWAG_VERSION)
else ifneq ($(CURRENT_SWAG_VERSION),$(SWAG_VERSION))
	@echo "Swag is installed, but version $(CURRENT_SWAG_VERSION) is not $(SWAG_VERSION). Installing swag version $(SWAG_VERSION)..."
	go install github.com/swaggo/swag/cmd/swag@$(SWAG_VERSION)
else
	@echo "Swag version $(SWAG_VERSION) is already installed."
endif

swaggo:
	swag fmt && swag init -g internal/delivery/rest/rest.go

run-dev: go run main.go

watch:
	@modd -f ./.modd/server.modd.conf

install-linter:
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.62.2

lint:
	golangci-lint run --print-issued-lines=false --exclude-use-default=false --enable=goimports  --enable=unconvert --enable=unparam --concurrency=4

# installing modd. ensure you have Go 1.17+ installed
# original docs: https://github.com/cortesi/modd
install-modd:
	go install github.com/cortesi/modd/cmd/modd@latest

# source of truth which mockery version will be used
install-mockery:
	go install github.com/vektra/mockery/v2@v2.52.2

migrate:
	go run main.go migrate

.PHONY: mocks
mocks:
	mockery

EXCLUDED_UNIT_TEST_PATHS := ./mocks ./docs
EXCLUDED_PATTERN := $(shell echo $(EXCLUDED_UNIT_TEST_PATHS) | sed 's/ /|/g')
UNIT_TEST_PACKAGES := $(shell go list ./... | grep -Ev '$(EXCLUDED_PATTERN)')

unit-tests:
	@echo "Running unit tests on: $(UNIT_TEST_PACKAGES)"
	go test -race -cover -coverprofile=coverage.out -count=1 $(UNIT_TEST_PACKAGES)
