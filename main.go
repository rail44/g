package main

import (
	"context"
	"fmt"
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
		res := openapi.GetAccountsIdBalance404Response{}
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
	accountId, err := server.queries.InsertAccount(ctx, sql.NullString{String: req.Body.Name, Valid: true})
	if err != nil {
		return nil, fmt.Errorf("Failed to query InsertAccount: %w", err)
	}

	err = server.queries.InsertBalance(ctx, accountId)
	if err != nil {
		return nil, fmt.Errorf("Failed to query InsertBalance: %w", err)
	}

	accountIdInt := int(accountId)
	res := openapi.PostAccountsRegister200JSONResponse{
		AccountId: &accountIdInt,
	}
	return res, nil
}

func (server Server) PostAccountsIdMint(ctx context.Context, req openapi.PostAccountsIdMintRequestObject) (openapi.PostAccountsIdMintResponseObject, error) {
	mintId, err := server.queries.InsertMint(ctx, sqlc.InsertMintParams{
                Account: int64(req.Id),
                Amount: req.Body.Amount,
        })
	if err != nil {
		return nil, fmt.Errorf("Failed to query InsertMint: %w", err)
	}

	txId, err := server.queries.InsertTransaction(ctx, sqlc.InsertTransactionParams{
                Mint:     sql.NullInt64 { Int64: mintId, Valid: true },
		Transfer: sql.NullInt64 {},
	})
	if err != nil {
		return nil, fmt.Errorf("Failed to query InsertMint: %w", err)
	}

	err = server.queries.IncrementBalance(ctx, sqlc.IncrementBalanceParams{
                Account: int64(req.Id),
                Amount: req.Body.Amount,
        })
	if err != nil {
		return nil, fmt.Errorf("Failed to query InsertMint: %w", err)
	}

	txIdInt := int(txId)
	res := openapi.PostAccountsIdMint200JSONResponse{
		TransactionId: &txIdInt,
	}
	return res, nil
}

func main() {
	db, err := sql.Open("postgres", "user=postgres dbname=g password=password host=localhost sslmode=disable")
	if err != nil {
		panic(err)
	}
	queries := sqlc.New(db)

	server := Server{queries: queries}
	http.ListenAndServe(":3000", openapi.Handler(openapi.NewStrictHandler(server, nil)))
}
