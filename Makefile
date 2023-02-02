POSTGRES_PASSWORD := password
POSTGRES_PORT := 5432
POSTGRES_DB := g

db/launch:
	docker run --rm -p $(POSTGRES_PORT):$(POSTGRES_PORT) -e POSTGRES_PASSWORD=$(POSTGRES_PASSWORD) -e POSTGRES_DB=$(POSTGRES_DB) postgres

db/schema:
	psql -f sqlc/schema.sql postgresql://postgres:$(POSTGRES_PASSWORD)@localhost:$(POSTGRES_PORT)/$(POSTGRES_DB)

openapi/generate: accounts/openapi.gen.go

accounts/openapi.gen.go: accounts/openapi.yml
	oapi-codegen -package accounts -generate types,chi-server,strict-server $< > $@
