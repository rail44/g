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
	db *sql.DB
}

func (server Server) GetAccountsIdBalance(ctx context.Context, req openapi.GetAccountsIdBalanceRequestObject) (openapi.GetAccountsIdBalanceResponseObject, error) {
	queries := sqlc.New(server.db)
	model, err := queries.GetBalance(ctx, int64(req.Id))

	if err == sql.ErrNoRows {
		res := openapi.GetAccountsIdBalance404Response{}
		return res, nil
	}

	if err != nil {
		return nil, fmt.Errorf("query GetBalance: %w", err)
	}

	res := openapi.GetAccountsIdBalance200JSONResponse{
		Balance: &model.Balance,
	}
	return res, nil
}

func (server Server) PostAccountsRegister(ctx context.Context, req openapi.PostAccountsRegisterRequestObject) (openapi.PostAccountsRegisterResponseObject, error) {
	tx, err := server.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback()
	queries := sqlc.New(tx)

	accountId, err := queries.InsertAccount(ctx, sql.NullString{String: req.Body.Name, Valid: true})
	if err != nil {
		return nil, fmt.Errorf("query InsertAccount: %w", err)
	}

	err = queries.InsertBalance(ctx, accountId)
	if err != nil {
		return nil, fmt.Errorf("query InsertBalance: %w", err)
	}

	accountIdInt := int(accountId)
	res := openapi.PostAccountsRegister200JSONResponse{
		AccountId: &accountIdInt,
	}

	return res, tx.Commit()
}

func (server Server) PostAccountsIdMint(ctx context.Context, req openapi.PostAccountsIdMintRequestObject) (openapi.PostAccountsIdMintResponseObject, error) {
	tx, err := server.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback()

	queries := sqlc.New(tx)
	mintId, err := queries.InsertMint(ctx, sqlc.InsertMintParams{
		Account: int64(req.Id),
		Amount:  req.Body.Amount,
	})
	if err != nil {
		return nil, fmt.Errorf("Failed to query InsertMint: %w", err)
	}

	txId, err := queries.InsertTransaction(ctx, sqlc.InsertTransactionParams{
		Mint:     sql.NullInt64{Int64: mintId, Valid: true},
		Transfer: sql.NullInt64{},
	})
	if err != nil {
		return nil, fmt.Errorf("Failed to query InsertMint: %w", err)
	}

	err = queries.IncrementBalance(ctx, sqlc.IncrementBalanceParams{
		Account: int64(req.Id),
		Amount:  req.Body.Amount,
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

	server := Server{db: db}
	http.ListenAndServe(":3000", openapi.Handler(openapi.NewStrictHandler(server, nil)))
}
