package main

import (
	"fmt"
	"context"
	"net/http"

	// "flag"
	"database/sql"

	_ "github.com/lib/pq"
	"github.com/rail44/g/openapi"
	"github.com/rail44/g/sqlc/generated"
)

type Secret struct {
	SupabaseKey string `toml:"supabaseKey"`
}

type GraphqlTranstport struct {
	supabaseKey *string
}

func (trastport *GraphqlTranstport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Add("apikey", *trastport.supabaseKey)
	return http.DefaultTransport.RoundTrip(req)
}

type Server struct {
        queries *sqlc.Queries
}

func (server Server) GetAccountsIdBalance(ctx context.Context, req openapi.GetAccountsIdBalanceRequestObject) (openapi.GetAccountsIdBalanceResponseObject, error) {
        model, err := server.queries.GetBalance(ctx, int64(req.Id))

        if err == sql.ErrNoRows {
                res := openapi.GetAccountsIdBalance404Response {}
                return res, nil
        }

        if err != nil {
                return nil, fmt.Errorf("Failed to query GetBalance: %w", err)
        }

	res := openapi.GetAccountsIdBalance200JSONResponse{
		Balance: &model.Balance,
	}
	return res, nil
}

func (server Server) PostAccountsRegister(ctx context.Context, req openapi.PostAccountsRegisterRequestObject) (openapi.PostAccountsRegisterResponseObject, error) {
        int64Id, err := server.queries.RegisterAccount(ctx, sql.NullString{ String: req.Body.Name, Valid: true })
        if err != nil {
                return nil, fmt.Errorf("Failed to query RegisterAccount: %w", err)
        }

        id := int(int64Id)
	res := openapi.PostAccountsRegister200JSONResponse{
		AccountId: &id,
	}
	return res, nil
}

func main() {
        db, err := sql.Open("postgres", "user=postgres dbname=g password=password host=localhost sslmode=disable")
        if err != nil {
                panic(err)
        }
        queries := sqlc.New(db)

        server := Server{ queries: queries }
	http.ListenAndServe(":3000", openapi.Handler(openapi.NewStrictHandler(server, nil)))
}
