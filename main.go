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

type Server struct {
	db *sql.DB
}

func (server Server) GetAccountsIdBalance(ctx context.Context, req openapi.GetAccountsIdBalanceRequestObject) (openapi.GetAccountsIdBalanceResponseObject, error) {
	queries := sqlc.New(server.db)

	balanceDecimal, err := queries.GetBalance(ctx, int64(req.Id))
	if err == sql.ErrNoRows {
		var res openapi.GetAccountsIdBalance404Response
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

func (server Server) GetAccountsIdTransactions(ctx context.Context, req openapi.GetAccountsIdTransactionsRequestObject) (openapi.GetAccountsIdTransactionsResponseObject, error) {
	queries := sqlc.New(server.db)

	_, err := queries.GetAccount(ctx, int64(req.Id))
	if err == sql.ErrNoRows {
		var res openapi.GetAccountsIdTransactions404Response
		return res, nil
	}

	if err != nil {
		return nil, fmt.Errorf("query GetAccount: %w", err)
	}

	transactions, err := queries.GetTransactions(ctx, int64(req.Id))
	if err != nil {
		return nil, fmt.Errorf("query GetTransactions: %w", err)
	}

	var res []interface{}
	for i := range transactions {
		tx := transactions[i]
		if tx.MintID.Valid {
			amount, err := strconv.Atoi(tx.MintAmount.String)
			if err != nil {
				return nil, fmt.Errorf("parse amount as decimal: %w", err)
			}

			_type := openapi.MintTypeMint
			mint := openapi.Mint{Amount: &amount, InsertedAt: &tx.InsertedAt, Type: &_type}
			res = append(res, mint)
			continue
		}

		if tx.SpendID.Valid {
			amount, err := strconv.Atoi(tx.SpendAmount.String)
			if err != nil {
				return nil, fmt.Errorf("parse amount as decimal: %w", err)
			}

			_type := openapi.SpendTypeSpend
			spend := openapi.Spend{Amount: &amount, InsertedAt: &tx.InsertedAt, Type: &_type}
			res = append(res, spend)
			continue
		}

		if tx.TransferID.Valid {
			amount, err := strconv.Atoi(tx.TransferAmount.String)
			if err != nil {
				return nil, fmt.Errorf("parse amount as decimal: %w", err)
			}

			_type := openapi.TransferTypeTransfer
			recipient := int(tx.TransferRecipient.Int64)
			transfer := openapi.Transfer{Amount: &amount, InsertedAt: &tx.InsertedAt, Type: &_type, Recipient: &recipient}
			res = append(res, transfer)
			continue
		}
		return nil, fmt.Errorf("failed to determine type of transaction")
	}

	return openapi.GetAccountsIdTransactions200JSONResponse(res), nil
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

	if req.Body.Amount <= 0 {
		res := openapi.PostAccountsIdMint400TextResponse("amount should be positive value")
		return res, nil
	}
	amount := strconv.Itoa(req.Body.Amount)

	accountId := int64(req.Id)
	_, err = queries.GetAccount(ctx, accountId)
	if err == sql.ErrNoRows {
		var res openapi.PostAccountsIdMint404Response
		return res, nil
	}

	mintId, err := queries.InsertMint(ctx, amount)
	if err != nil {
		return nil, fmt.Errorf("Failed to query InsertMint: %w", err)
	}

	txId, err := queries.InsertTransaction(ctx, sqlc.InsertTransactionParams{
		Account: accountId,
		Mint:    sql.NullInt64{Int64: mintId, Valid: true},
	})
	if err != nil {
		return nil, fmt.Errorf("Failed to query InsertTransaction: %w", err)
	}

	err = queries.IncrementBalance(ctx, sqlc.IncrementBalanceParams{
		Account: accountId,
		Amount:  amount,
	})
	if err != nil {
		return nil, fmt.Errorf("Failed to query IncrementBalance: %w", err)
	}

	txIdInt := int(txId)
	res := openapi.PostAccountsIdMint200JSONResponse{
		TransactionId: &txIdInt,
	}

	return res, tx.Commit()
}

func (server Server) PostAccountsIdSpend(ctx context.Context, req openapi.PostAccountsIdSpendRequestObject) (openapi.PostAccountsIdSpendResponseObject, error) {
	tx, err := server.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback()

	queries := sqlc.New(tx)

	if req.Body.Amount <= 0 {
		res := openapi.PostAccountsIdSpend400TextResponse("amount should be positive value")
		return res, nil
	}
	amount := strconv.Itoa(req.Body.Amount)

	accountId := int64(req.Id)
	_, err = queries.GetAccount(ctx, accountId)
	if err == sql.ErrNoRows {
		res := openapi.PostAccountsIdSpend404Response{}
		return res, nil
	}

	spendId, err := queries.InsertSpend(ctx, amount)
	if err != nil {
		return nil, fmt.Errorf("Failed to query InsertSpend: %w", err)
	}

	txId, err := queries.InsertTransaction(ctx, sqlc.InsertTransactionParams{
		Account: accountId,
		Spend:   sql.NullInt64{Int64: spendId, Valid: true},
	})
	if err != nil {
		return nil, fmt.Errorf("Failed to query InsertTransaction: %w", err)
	}

	err = queries.DecrementBalance(ctx, sqlc.DecrementBalanceParams{
		Account: accountId,
		Amount:  amount,
	})
	if err != nil {
		return nil, fmt.Errorf("Failed to query IncrementBalance: %w", err)
	}

	txIdInt := int(txId)
	res := openapi.PostAccountsIdSpend200JSONResponse{
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

	if req.Body.Amount <= 0 {
		res := openapi.PostAccountsIdTransfer400TextResponse("amount should be positive value")
		return res, nil
	}
	amount := strconv.Itoa(req.Body.Amount)

	accountId := int64(req.Id)
	balanceDecimal, err := queries.GetBalance(ctx, accountId)
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

	recipientId := int64(req.Body.To)
	_, err = queries.GetAccount(ctx, recipientId)
	if err == sql.ErrNoRows {
		msg := fmt.Sprintf("recipient was not found by id %d", recipientId)
		res := openapi.PostAccountsIdTransfer400TextResponse(msg)
		return res, nil
	}

	transferId, err := queries.InsertTransfer(ctx, sqlc.InsertTransferParams{
		Recipient: recipientId,
		Amount:    amount,
	})
	if err != nil {
		return nil, fmt.Errorf("query InsertTransfer: %w", err)
	}

	txId, err := queries.InsertTransaction(ctx, sqlc.InsertTransactionParams{
		Account:  accountId,
		Transfer: sql.NullInt64{Int64: transferId, Valid: true},
	})
	if err != nil {
		return nil, fmt.Errorf("query InsertTransaction: %w", err)
	}

	err = queries.DecrementBalance(ctx, sqlc.DecrementBalanceParams{
		Account: accountId,
		Amount:  amount,
	})
	if err != nil {
		return nil, fmt.Errorf("query DecrementBalance: %w", err)
	}

	err = queries.IncrementBalance(ctx, sqlc.IncrementBalanceParams{
		Account: recipientId,
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
