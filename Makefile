ifeq ($(POSTGRES_SETUP_TEST),)
	POSTGRES_SETUP_TEST := user=postgres password=postgres dbname=postgres host=localhost port=5432 sslmode=disable
endif

INTERNAL_PKG_PATH=$(CURDIR)/internal/pkg
MIGRATION_FOLDER=$(INTERNAL_PKG_PATH)/db/migrations

.PHONY: .up docker-compose
.up:
	docker-compose -f build/docker-compose.yml up -d

.PHONY: test-migration-up
test-migration-up:
	goose -dir "$(MIGRATION_FOLDER)" postgres "$(POSTGRES_SETUP_TEST)" up

.PHONY: test-migration-down
test-migration-down:
	goose -dir "$(MIGRATION_FOLDER)" postgres "$(POSTGRES_SETUP_TEST)" down

.PHONY: .test-unit
.test-unit:
	$(info Running tests...)
	go test ./internal/pkg/server -run unit -v

.PHONY: .test-integration
.test-integration:
	$(info Running tests...)
	go test ./internal/pkg/server -run integration -v
	go test ./internal/pkg/repository/postgresql -run integration -v