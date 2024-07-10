include .env
export

base_migrations_command=goose -dir="./migrations" clickhouse "tcp://$(CLICKHOUSE_HOST):9000?username=$(CLICKHOUSE_ADMIN_USER)&password=$(CLICKHOUSE_ADMIN_PASSWORD)"

migrations-up:
	$(base_migrations_command) up
migrations-down:
	$(base_migrations_command) down

run:
	go run cmd/api/main.go
run-debug:
	go run cmd/api/main.go -debug

generate:
	go generate ./...
test: generate
	go test ./... -v
