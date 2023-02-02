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

func (model *Model) GetBalance(ctx context.Context, id int) (int, error) {
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

func (model *Model) Register(ctx context.Context, name string) (int, error) {
	tx, err := model.db.Begin()
	if err != nil {
		return 0, fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback()
	queries := sqlc.New(tx)

	accountId, err := queries.InsertAccount(ctx, sql.NullString{String: name, Valid: true})
	if err != nil {
		return 0, fmt.Errorf("query InsertAccount: %w", err)
	}

	err = queries.InsertBalance(ctx, accountId)
	if err != nil {
		return 0, fmt.Errorf("query InsertBalance: %w", err)
	}

	err = tx.Commit()
  if err != nil {
		return 0, fmt.Errorf("commit: %w", err)
  }

  return int(accountId), nil
}
