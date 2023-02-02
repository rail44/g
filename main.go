package main

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

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

	balanceDecimal, err := queries.GetBalance(ctx, int64(req.Id))
	if err == sql.ErrNoRows {
		res := openapi.GetAccountsIdBalance404Response{}
		return res, nil
	}

	if err != nil {
		return nil, fmt.Errorf("query GetBalance: %w", err)
	}

	balance, err := strconv.Atoi(balanceDecimal)
	if err != nil {
		return nil, fmt.Errorf("parse balance as decimal: %w", err)
	}
	res := openapi.GetAccountsIdBalance200JSONResponse{
		Balance: &balance,
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

	name := req.Body.Name
	if len(name) == 0 {
		res := openapi.PostAccountsRegister400TextResponse("name is not presented")
		return res, nil
	}

	accountId, err := queries.InsertAccount(ctx, sql.NullString{String: name, Valid: true})
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

	amount := strconv.Itoa(req.Body.Amount)

	accountId := int64(req.Id)
	_, err = queries.GetAccount(ctx, accountId)
	if err == sql.ErrNoRows {
		res := openapi.PostAccountsIdMint404Response{}
		return res, nil
	}

	mintId, err := queries.InsertMint(ctx, sqlc.InsertMintParams{
		Account: accountId,
		Amount:  amount,
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
		Account: accountId,
		Amount:  amount,
	})
	if err != nil {
		return nil, fmt.Errorf("Failed to query InsertMint: %w", err)
	}

	txIdInt := int(txId)
	res := openapi.PostAccountsIdMint200JSONResponse{
		TransactionId: &txIdInt,
	}

	return res, tx.Commit()
}

func (server Server) PostAccountsIdTransfer(ctx context.Context, req openapi.PostAccountsIdTransferRequestObject) (openapi.PostAccountsIdTransferResponseObject, error) {
	tx, err := server.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback()

	queries := sqlc.New(tx)

	amount := strconv.Itoa(req.Body.Amount)

	fromAccountId := int64(req.Id)
	balanceDecimal, err := queries.GetBalance(ctx, fromAccountId)
	if err == sql.ErrNoRows {
		res := openapi.PostAccountsIdTransfer404Response{}
		return res, nil
	}

	balance, err := strconv.Atoi(balanceDecimal)
	if err != nil {
		return nil, fmt.Errorf("parse balance as decimal: %w", err)
	}

	if req.Body.Amount > balance {
		msg := fmt.Sprintf("tried to transfer %d, but balance was only %d", req.Body.Amount, balance)
		res := openapi.PostAccountsIdTransfer400TextResponse(msg)
		return res, nil
	}

	toAccountId := int64(req.Body.To)
	_, err = queries.GetAccount(ctx, toAccountId)
	if err == sql.ErrNoRows {
		msg := fmt.Sprintf("receipient was not found by id %d", toAccountId)
		res := openapi.PostAccountsIdTransfer400TextResponse(msg)
		return res, nil
	}

	transferId, err := queries.InsertTransfer(ctx, sqlc.InsertTransferParams{
		FromAccount: fromAccountId,
		ToAccount:   toAccountId,
		Amount:      amount,
	})
	if err != nil {
		return nil, fmt.Errorf("query InsertTransfer: %w", err)
	}

	txId, err := queries.InsertTransaction(ctx, sqlc.InsertTransactionParams{
		Transfer: sql.NullInt64{Int64: transferId, Valid: true},
		Mint:     sql.NullInt64{},
	})
	if err != nil {
		return nil, fmt.Errorf("query InsertTransaction: %w", err)
	}

	err = queries.DecrementBalance(ctx, sqlc.DecrementBalanceParams{
		Account: fromAccountId,
		Amount:  amount,
	})
	if err != nil {
		return nil, fmt.Errorf("query DecrementBalance: %w", err)
	}

	err = queries.IncrementBalance(ctx, sqlc.IncrementBalanceParams{
		Account: toAccountId,
		Amount:  amount,
	})
	if err != nil {
		return nil, fmt.Errorf("query IncrementBalance: %w", err)
	}

	txIdInt := int(txId)
	res := openapi.PostAccountsIdTransfer200JSONResponse{
		TransactionId: &txIdInt,
	}

	return res, tx.Commit()
}

func main() {
	db, err := sql.Open("postgres", "user=postgres dbname=g password=password host=localhost sslmode=disable")
	if err != nil {
		panic(err)
	}

	server := Server{db: db}
	http.ListenAndServe(":3000", openapi.Handler(openapi.NewStrictHandler(server, nil)))
}
