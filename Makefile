include .env
LOCAL_BIN=$(CURDIR)/bin

up:
	docker compose up -d

down:
	docker compose down

install-deps:
	GOBIN=$(LOCAL_BIN) go install -tags "postgres" github.com/golang-migrate/migrate/v4/cmd/migrate@latest

migrate-up:
	$(LOCAL_BIN)/migrate -path ./migrations -database "postgresql://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@localhost:5432/$(POSTGRES_DB)?sslmode=disable" up

migrate-down:
	$(LOCAL_BIN)/migrate -path ./migrations -database "postgresql://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@localhost:5432/$(POSTGRES_DB)?sslmode=disable" down

migrate-drop:
	$(LOCAL_BIN)/migrate -path ./migrations -database "postgresql://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@localhost:5432/$(POSTGRES_DB)?sslmode=disable" drop table

.PHONY: up down install-deps migrate-up migrate-down migrate-drop