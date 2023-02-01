generate/openapi: openapi/generated.go

openapi/generated.go: openapi.yml
	oapi-codegen -package openapi -generate types,chi-server,strict-server openapi.yml > $@
