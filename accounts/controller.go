package accounts

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"database/sql"
)

var ValidationError = errors.New("Validation Error")

// 各RouteがErrorをreturnした場合に、その属性によってHTTPレスポンスを分岐させるためのカスタムハンドラ
var ServerOptions = StrictHTTPServerOptions{
	RequestErrorHandlerFunc: func(w http.ResponseWriter, r *http.Request, err error) {
		http.Error(w, err.Error(), http.StatusBadRequest)
	},
	ResponseErrorHandlerFunc: func(w http.ResponseWriter, r *http.Request, err error) {
		if errors.Is(err, NotFoundError) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		if errors.Is(err, ValidationError) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if errors.Is(err, DomainError) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		http.Error(w, err.Error(), http.StatusInternalServerError)
	},
}

func NewController(model *Model) http.Handler {
	controller := Controller{model: model}
	return Handler(NewStrictHandlerWithOptions(controller, nil, ServerOptions))
}

type Controller struct {
	db    *sql.DB
	model *Model
}

// GET /balance
func (controller Controller) Balance(ctx context.Context, req BalanceRequestObject) (BalanceResponseObject, error) {
	balance, err := controller.model.GetBalance(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	res := Balance200JSONResponse{
		Balance: balance,
	}
	return res, nil
}

// GET /transactions
func (controller Controller) Transactions(ctx context.Context, req TransactionsRequestObject) (TransactionsResponseObject, error) {
	transactions, err := controller.model.GetTransactions(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	res := Transactions200JSONResponse(
		transactions,
	)
	return res, nil
}

// POST /
func (controller Controller) Register(ctx context.Context, req RegisterRequestObject) (RegisterResponseObject, error) {
	if len(req.Body.Name) == 0 {
		return nil, fmt.Errorf("name is not presented: %w", ValidationError)
	}

	id, err := controller.model.Register(ctx, req.Body.Name)
	if err != nil {
		return nil, err
	}

	res := Register200JSONResponse{
		AccountId: &id,
	}
	return res, nil
}

// POST /{id}/mint
func (controller Controller) Mint(ctx context.Context, req MintRequestObject) (MintResponseObject, error) {
	if req.Body.Amount <= 0 {
		return nil, fmt.Errorf("amount should be positive value %d: %w", req.Body.Amount, ValidationError)
	}

	txId, err := controller.model.Mint(ctx, req.Id, req.Body.Amount)
	if err != nil {
		return nil, err
	}

	res := Mint200JSONResponse{
		TransactionId: int(txId),
	}
	return res, nil
}

// POST /{id}/spend
func (controller Controller) Spend(ctx context.Context, req SpendRequestObject) (SpendResponseObject, error) {
	if req.Body.Amount <= 0 {
		return nil, fmt.Errorf("amount should be positive value %d: %w", req.Body.Amount, ValidationError)
	}

	txId, err := controller.model.Spend(ctx, req.Id, req.Body.Amount)
	if err != nil {
		return nil, err
	}

	res := Spend200JSONResponse{
		TransactionId: txId,
	}
	return res, nil
}

// POST /{id}/transfer
func (controller Controller) Transfer(ctx context.Context, req TransferRequestObject) (TransferResponseObject, error) {
	if req.Body.Amount <= 0 {
		return nil, fmt.Errorf("amount should be positive value %d: %w", req.Body.Amount, ValidationError)
	}

	if req.Body.Recipient <= 0 {
		return nil, fmt.Errorf("recipient id should be positive value %d: %w", req.Body.Amount, ValidationError)
	}

	txId, err := controller.model.Transfer(ctx, req.Id, req.Body.Recipient, req.Body.Amount)
	if err != nil {
		return nil, err
	}

	res := Transfer200JSONResponse{
		TransactionId: txId,
	}
	return res, nil
}
