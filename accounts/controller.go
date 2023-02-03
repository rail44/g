package accounts

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"database/sql"

	"github.com/rail44/g/sqlc/generated"
)

func NewController(db *sql.DB) http.Handler {
	model := Model{db: db}
  controller := Controller{ db: db, model: &model }

  options := StrictHTTPServerOptions {
		RequestErrorHandlerFunc: func(w http.ResponseWriter, r *http.Request, err error) {
			http.Error(w, err.Error(), http.StatusBadRequest)
		},
		ResponseErrorHandlerFunc: func(w http.ResponseWriter, r *http.Request, err error) {
      if errors.Is(err, ValidationError) {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
      }

			http.Error(w, err.Error(), http.StatusInternalServerError)
		},
	}

	return HandlerWithOptions(NewStrictHandlerWithOptions(controller, nil, options), ChiServerOptions{
  })
}

type Controller struct {
	db    *sql.DB
	model *Model
}

func (controller Controller) Balance(ctx context.Context, req BalanceRequestObject) (BalanceResponseObject, error) {
	err := controller.model.Exists(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	balance, err := controller.model.GetBalance(ctx, req.Id)

	if err != nil {
		return nil, err
	}

	res := Balance200JSONResponse{
		Balance: &balance,
	}
	return res, nil
}

func (controller Controller) Transactions(ctx context.Context, req TransactionsRequestObject) (TransactionsResponseObject, error) {
	err := controller.model.Exists(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	transactions, err := controller.model.GetTransactions(ctx, req.Id)

	if err != nil {
		return nil, err
	}

	res := Transactions200JSONResponse(
		transactions,
	)
	return res, nil
}

func (controller Controller) Register(ctx context.Context, req RegisterRequestObject) (RegisterResponseObject, error) {
	id, err := controller.model.Register(ctx, req.Body.Name)
	if err != nil {
		return nil, err
	}
	res := Register200JSONResponse{
		AccountId: &id,
	}

	return res, nil
}

func (controller Controller) Mint(ctx context.Context, req MintRequestObject) (MintResponseObject, error) {
	txId, err := controller.model.Mint(ctx, req.Id, req.Body.Amount)
	if err != nil {
		return nil, err
	}

	txIdInt := int(txId)
	res := Mint200JSONResponse{
		TransactionId: &txIdInt,
	}

	return res, nil
}

func (controller Controller) Spend(ctx context.Context, req SpendRequestObject) (SpendResponseObject, error) {
  txId, err := controller.model.Spend(ctx, req.Id, req.Body.Amount)
	if err != nil {
		return nil, err
	}

	res := Spend200JSONResponse{
		TransactionId: &txId,
	}

	return res, nil
}

func (controller Controller) Transfer(ctx context.Context, req TransferRequestObject) (TransferResponseObject, error) {
	err := controller.model.Exists(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	err = controller.model.Exists(ctx, req.Body.To)
	if err != nil {
		return nil, err
	}

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

	if err != nil {
		return nil, fmt.Errorf("GetBalance: %w", err)
	}

	if req.Body.Amount > balance {
		msg := fmt.Sprintf("tried to transfer %d, but balance was only %d", req.Body.Amount, balance)
		res := Transfer400TextResponse(msg)
		return res, nil
	}
	recipientId := int64(req.Body.To)
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
