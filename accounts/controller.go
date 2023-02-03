package accounts

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"database/sql"
	"github.com/rail44/g/sqlc/generated"
)

func NewController(db *sql.DB) http.Handler {
	model := Model{db: db}
	controller := Controller{db: db, model: &model}
	return Handler(NewStrictHandler(controller, nil))
}

type Controller struct {
	db    *sql.DB
	model *Model
}

func (controller Controller) Balance(ctx context.Context, req BalanceRequestObject) (BalanceResponseObject, error) {
	exists, err := controller.model.Exists(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	if !exists {
		var res Balance404Response
		return res, nil
	}

	balance, err := controller.model.GetBalance(ctx, req.Id)

	if err != nil {
		return nil, fmt.Errorf("GetBalance: %w", err)
	}

	res := Balance200JSONResponse{
		Balance: &balance,
	}
	return res, nil
}

func (controller Controller) Transactions(ctx context.Context, req TransactionsRequestObject) (TransactionsResponseObject, error) {
	exists, err := controller.model.Exists(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	if !exists {
		var res Transactions404Response
		return res, nil
	}

	transactions, err := controller.model.GetTransactions(ctx, req.Id)

	if err != nil {
		return nil, fmt.Errorf("GetBalance: %w", err)
	}

	res := Transactions200JSONResponse(
		transactions,
	)
	return res, nil
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
	exists, err := controller.model.Exists(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	if !exists {
		var res Mint404Response
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
	exists, err := controller.model.Exists(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	if !exists {
		var res Spend404Response
		return res, nil
	}

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
		var res Spend404Response
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
	exists, err := controller.model.Exists(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	if !exists {
		var res Transfer404Response
		return res, nil
	}

	exists, err = controller.model.Exists(ctx, req.Body.To)
	if err != nil {
		return nil, err
	}

	if !exists {
		msg := fmt.Sprintf("recipient was not found by id %d", req.Body.To)
		res := Transfer400TextResponse(msg)
		return res, nil
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
