PORT := 3000
PG_HOST := localhost
PG_PORT := 5432
PG_USER := postgres
PG_PASS := password
PG_DB := g

g/run:
	@go run . -port=$(PORT) -dbuser=$(PG_USER) -dbpass=$(PG_PASS) -dbhost=$(PG_HOST) -dbport=$(PG_PORT) -dbname=$(PG_DB)

db/up:
	docker run -ti --rm -p $(PG_PORT):$(PG_PORT) -e POSTGRES_USER=$(PG_USER) -e POSTGRES_PASSWORD=$(PG_PASS) -e POSTGRES_DB=$(PG_DB) postgres

db/schema db/reset:
	psql -f sqlc/schema.sql postgresql://$(PG_USER):$(PG_PASS)@$(PG_HOST):$(PG_PORT)/$(PG_DB)

openapi/generate: accounts/openapi.gen.go
accounts/openapi.gen.go: accounts/openapi.yml
	oapi-codegen -package accounts -generate types,chi-server,strict-server $< > $@

sqlc/generate:
	sqlc generate

fmt:
	go fmt ./...
