package accounts

import (
	"fmt"
	"strconv"
	"context"
	"database/sql"
	"github.com/rail44/g/sqlc/generated"
)

type Model struct {
	db *sql.DB
}

func (model *Model) getBalance(ctx context.Context, id int) (int, error) {
	queries := sqlc.New(model.db)

	balanceDecimal, err := queries.GetBalance(ctx, int64(id))
	if err != nil {
		return 0, fmt.Errorf("query GetBalance: %w", err)
	}

	balance, err := strconv.Atoi(balanceDecimal)
	if err != nil {
		return 0, fmt.Errorf("parse balance as decimal: %w", err)
	}
	return balance, nil
}
