package accounts

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	// "flag"
	"database/sql"

	_ "github.com/lib/pq"
	"github.com/rail44/g/sqlc/generated"
)

func NewController(db *sql.DB) http.Handler {
  model := Model { db: db }
  controller := Controller { db: db, model: &model }
	return Handler(NewStrictHandler(controller, nil))
}

type Controller struct {
	db *sql.DB
	model *Model
}

func (controller Controller) Balance(ctx context.Context, req BalanceRequestObject) (BalanceResponseObject, error) {
  balance, err := controller.model.GetBalance(ctx, req.Id)

  if err == sql.ErrNoRows {
		var res Balance404Response
    return res, nil
  }

	if err != nil {
		return nil, fmt.Errorf("GetBalance: %w", err)
	}

	res := Balance200JSONResponse{
		Balance: &balance,
	}
	return res, nil
}

func (controller Controller) Transactions(ctx context.Context, req TransactionsRequestObject) (TransactionsResponseObject, error) {
	queries := sqlc.New(controller.db)

	_, err := queries.GetAccount(ctx, int64(req.Id))
	if err == sql.ErrNoRows {
		var res Transactions404Response
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

			_type := MintTypeMint
			mint := Mint{Amount: &amount, InsertedAt: &tx.InsertedAt, Type: &_type}
			res = append(res, mint)
			continue
		}

		if tx.SpendID.Valid {
			amount, err := strconv.Atoi(tx.SpendAmount.String)
			if err != nil {
				return nil, fmt.Errorf("parse amount as decimal: %w", err)
			}

			_type := SpendTypeSpend
			spend := Spend{Amount: &amount, InsertedAt: &tx.InsertedAt, Type: &_type}
			res = append(res, spend)
			continue
		}

		if tx.TransferID.Valid {
			amount, err := strconv.Atoi(tx.TransferAmount.String)
			if err != nil {
				return nil, fmt.Errorf("parse amount as decimal: %w", err)
			}

			_type := TransferTypeTransfer
			recipient := int(tx.TransferRecipient.Int64)
			transfer := Transfer{Amount: &amount, InsertedAt: &tx.InsertedAt, Type: &_type, Recipient: &recipient}
			res = append(res, transfer)
			continue
		}
		return nil, fmt.Errorf("failed to determine type of transaction")
	}

	return Transactions200JSONResponse(res), nil
}

func (controller Controller) Register(ctx context.Context, req RegisterRequestObject) (RegisterResponseObject, error) {
	name := req.Body.Name
	if len(name) == 0 {
		res := Register400TextResponse("name is not presented")
		return res, nil
	}

	id, err := controller.model.Register(ctx, name)
	if err != nil {
		return nil, fmt.Errorf("Register: %w", err)
	}
	res := Register200JSONResponse{
		AccountId: &id,
	}

	return res, nil
}

func (controller Controller) Mint(ctx context.Context, req MintRequestObject) (MintResponseObject, error) {
	if req.Body.Amount <= 0 {
		res := Mint400TextResponse("amount should be positive value")
		return res, nil
	}

  txId, err := controller.model.Mint(ctx, req.Id, req.Body.Amount)
  if err == sql.ErrNoRows {
		var res Mint404Response
    return res, nil
  }

	txIdInt := int(txId)
	res := Mint200JSONResponse{
		TransactionId: &txIdInt,
	}

	return res, nil
}

func (controller Controller) Spend(ctx context.Context, req SpendRequestObject) (SpendResponseObject, error) {
	tx, err := controller.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback()

	queries := sqlc.New(tx)

	if req.Body.Amount <= 0 {
		res := Spend400TextResponse("amount should be positive value")
		return res, nil
	}
	amount := strconv.Itoa(req.Body.Amount)

	balance, err := controller.model.GetBalance(ctx, req.Id)
	if err == sql.ErrNoRows {
    var res Spend404Response;
		return res, nil
	}

	if err != nil {
		return nil, fmt.Errorf("getBalance: %w", err)
	}

	if req.Body.Amount > balance {
		msg := fmt.Sprintf("tried to spend %d, but balance was only %d", req.Body.Amount, balance)
		res := Spend400TextResponse(msg)
		return res, nil
	}

	spendId, err := queries.InsertSpend(ctx, amount)
	if err != nil {
		return nil, fmt.Errorf("query InsertSpend: %w", err)
	}

	accountId := int64(req.Id)
	txId, err := queries.InsertTransaction(ctx, sqlc.InsertTransactionParams{
		Account: accountId,
		Spend:   sql.NullInt64{Int64: spendId, Valid: true},
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

	txIdInt := int(txId)
	res := Spend200JSONResponse{
		TransactionId: &txIdInt,
	}

	return res, tx.Commit()
}

func (controller Controller) Transfer(ctx context.Context, req TransferRequestObject) (TransferResponseObject, error) {
	tx, err := controller.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback()

	queries := sqlc.New(tx)

	if req.Body.Amount <= 0 {
		res := Transfer400TextResponse("amount should be positive value")
		return res, nil
	}
	amount := strconv.Itoa(req.Body.Amount)

	accountId := int64(req.Id)
	balance, err := controller.model.GetBalance(ctx, req.Id)
	if err == sql.ErrNoRows {
    var res Transfer404Response
		return res, nil
	}

	if err != nil {
		return nil, fmt.Errorf("getBalance: %w", err)
	}

	if req.Body.Amount > balance {
		msg := fmt.Sprintf("tried to transfer %d, but balance was only %d", req.Body.Amount, balance)
		res := Transfer400TextResponse(msg)
		return res, nil
	}

	recipientId := int64(req.Body.To)
	_, err = queries.GetAccount(ctx, recipientId)
	if err == sql.ErrNoRows {
		msg := fmt.Sprintf("recipient was not found by id %d", recipientId)
		res := Transfer400TextResponse(msg)
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
	res := Transfer200JSONResponse{
		TransactionId: &txIdInt,
	}

	return res, tx.Commit()
}
