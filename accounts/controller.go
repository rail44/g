package accounts

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"database/sql"
)

var ValidationError = errors.New("Validation Error")
var ServerOptions = StrictHTTPServerOptions{
	RequestErrorHandlerFunc: func(w http.ResponseWriter, r *http.Request, err error) {
		http.Error(w, err.Error(), http.StatusBadRequest)
	},
	ResponseErrorHandlerFunc: func(w http.ResponseWriter, r *http.Request, err error) {
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

	return HandlerWithOptions(NewStrictHandlerWithOptions(controller, nil, ServerOptions), ChiServerOptions{})
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
		Balance: balance,
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

func (controller Controller) Mint(ctx context.Context, req MintRequestObject) (MintResponseObject, error) {
	txId, err := controller.model.Mint(ctx, req.Id, req.Body.Amount)
	if err != nil {
		return nil, err
	}

	res := Mint200JSONResponse{
		TransactionId: int(txId),
	}
	return res, nil
}

func (controller Controller) Spend(ctx context.Context, req SpendRequestObject) (SpendResponseObject, error) {
	txId, err := controller.model.Spend(ctx, req.Id, req.Body.Amount)
	if err != nil {
		return nil, err
	}

	res := Spend200JSONResponse{
		TransactionId: txId,
	}
	return res, nil
}

func (controller Controller) Transfer(ctx context.Context, req TransferRequestObject) (TransferResponseObject, error) {
	txId, err := controller.model.Transfer(ctx, req.Id, req.Body.To, req.Body.Amount)
	if err != nil {
		return nil, err
	}

	res := Transfer200JSONResponse{
		TransactionId: txId,
	}
	return res, nil
}
