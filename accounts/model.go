package accounts

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/rail44/g/sqlc/generated"
	"strconv"
)

type Model struct {
	db *sql.DB
}

func (model *Model) Exists(ctx context.Context, id int) (bool, error) {
	queries := sqlc.New(model.db)

	_, err := queries.GetAccount(ctx, int64(id))

	if err == sql.ErrNoRows {
		return false, nil
	}

	if err != nil {
		return false, fmt.Errorf("Exists: %w", err)
	}

	return true, nil
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

func (model *Model) Mint(ctx context.Context, accountId int, amount int) (int, error) {
	tx, err := model.db.Begin()
	if err != nil {
		return 0, fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback()
	queries := sqlc.New(tx)

	amountDecimal := strconv.Itoa(amount)
	mintId, err := queries.InsertMint(ctx, amountDecimal)
	if err != nil {
		return 0, fmt.Errorf("query InsertMint: %w", err)
	}

	accountIdInt64 := int64(accountId)
	txId, err := queries.InsertTransaction(ctx, sqlc.InsertTransactionParams{
		Account: accountIdInt64,
		Mint:    sql.NullInt64{Int64: mintId, Valid: true},
	})
	if err != nil {
		return 0, fmt.Errorf("query InsertTransaction: %w", err)
	}

	err = queries.IncrementBalance(ctx, sqlc.IncrementBalanceParams{
		Account: accountIdInt64,
		Amount:  strconv.Itoa(amount),
	})
	if err != nil {
		return 0, fmt.Errorf("query IncrementBalance: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return 0, fmt.Errorf("commit: %w", err)
	}

	return int(txId), nil
}

func (model *Model) GetTransactions(ctx context.Context, accountId int) ([]interface{}, error) {
	queries := sqlc.New(model.db)
	transactions, err := queries.GetTransactions(ctx, int64(accountId))
	if err != nil {
		return nil, fmt.Errorf("query GetTransactions: %w", err)
	}

	var result []interface{}
	for _, v := range transactions {
		row, err := mapToSubtype(v)
		if err != nil {
			return nil, fmt.Errorf("mapToSubtype: %w", err)
		}

		result = append(result, row)
	}

	return result, nil
}

func (model *Model) Spend(ctx context.Context, accountId int, amount int) (int, error) {
	tx, err := model.db.Begin()
	if err != nil {
		return 0, fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback()
	queries := sqlc.New(tx)

	amountDecimal := strconv.Itoa(amount)
	mintId, err := queries.InsertMint(ctx, amountDecimal)
	if err != nil {
		return 0, fmt.Errorf("query InsertMint: %w", err)
	}

	accountIdInt64 := int64(accountId)
	txId, err := queries.InsertTransaction(ctx, sqlc.InsertTransactionParams{
		Account: accountIdInt64,
		Mint:    sql.NullInt64{Int64: mintId, Valid: true},
	})
	if err != nil {
		return 0, fmt.Errorf("query InsertTransaction: %w", err)
	}

	err = queries.DecrementBalance(ctx, sqlc.DecrementBalanceParams{
		Account: accountIdInt64,
		Amount:  strconv.Itoa(amount),
	})
	if err != nil {
		return 0, fmt.Errorf("query DecrementBalance: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return 0, fmt.Errorf("commit: %w", err)
	}

	return int(txId), nil
}
