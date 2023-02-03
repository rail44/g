package accounts

import (
	"context"

	"database/sql"
	"errors"
	"fmt"
	"github.com/rail44/g/sqlc/generated"
	"strconv"
)

var NotFoundError = errors.New("NotFound")

// データを問い合わせて発覚するロジックのエラー
var DomainError = errors.New("Domain Error")

// 関数をラップして単一トランザクションでクエリを実行するためのユーティリティ関数
// 下の方に使用例があります
func WithTransaction[T interface{}](db *sql.DB, f func(tx *sql.Tx) (T, error)) (T, error) {
	var v T
	tx, err := db.Begin()
	if err != nil {
		return v, fmt.Errorf("begining transaction: %w", err)
	}

	// Commit()まで到達せずにスコープを抜けた場合はRollback
	defer tx.Rollback()

	v, err = f(tx)
	if err != nil {
		return v, err
	}

	err = tx.Commit()
	if err != nil {
		return v, fmt.Errorf("commit: %w", err)
	}

	return v, nil
}

// sqlcで生成したクエリを発行するためのドメインモデル
// QueryとCommandに分けるのもありかもしれないです
// (Commandについてはinstantiate毎にトランザクションを発行してしまうと見通しがよくなるかもしれない
type Model struct {
	db *sql.DB
}

func NewModel(db *sql.DB) *Model {
	return &Model{db: db}
}

// idのaccountsが存在していなければNotFound Error
func (model *Model) Exists(ctx context.Context, id int) error {
	queries := sqlc.New(model.db)

	_, err := queries.GetAccount(ctx, int64(id))

	if err == sql.ErrNoRows {
		// NotFoundErrorをwrapしてearly return
		return fmt.Errorf("Not found account by id %d: %w", id, NotFoundError)
	}

	if err != nil {
		// こっちはUnexpectedErrorなので500になります
		return fmt.Errorf("queryng GetAccount: %w", err)
	}

	return nil
}

// idのaccountsがamount以上の残高をもっていなければDomainError
func (model *Model) HasEnough(ctx context.Context, id int, amount int) error {
	balance, err := model.GetBalance(ctx, id)

	if err != nil {
		return err
	}

	if amount > balance {
		return fmt.Errorf("%d amount was requested, but balance was only %d: %w", amount, balance, DomainError)
	}

	return nil
}

func (model *Model) GetBalance(ctx context.Context, id int) (int, error) {
	queries := sqlc.New(model.db)

	err := model.Exists(ctx, id)
	if err != nil {
		return 0, err
	}

	balanceDecimal, err := queries.GetBalance(ctx, int64(id))
	if err != nil {
		return 0, fmt.Errorf("querying GetBalance: %w", err)
	}

	balance, err := strconv.Atoi(balanceDecimal)
	if err != nil {
		return 0, fmt.Errorf("parsing balance as decimal: %w", err)
	}
	return balance, nil
}

func (model *Model) Register(ctx context.Context, name string) (int, error) {
	return WithTransaction(model.db, func(tx *sql.Tx) (int, error) {
		queries := sqlc.New(tx)
		accountId, err := queries.InsertAccount(ctx, sql.NullString{String: name, Valid: true})
		if err != nil {
			return 0, fmt.Errorf("querying InsertAccount: %w", err)
		}

		err = queries.InsertBalance(ctx, accountId)
		if err != nil {
			return 0, fmt.Errorf("querying InsertBalance: %w", err)
		}

		return int(accountId), nil
	})
}

func (model *Model) Mint(ctx context.Context, accountId int, amount int) (int, error) {
	return WithTransaction(model.db, func(tx *sql.Tx) (int, error) {
		queries := sqlc.New(tx)

		err := model.Exists(ctx, accountId)
		if err != nil {
			return 0, err
		}

		amountDecimal := strconv.Itoa(amount)
		mintId, err := queries.InsertMint(ctx, amountDecimal)
		if err != nil {
			return 0, fmt.Errorf("querying InsertMint: %w", err)
		}

		accountIdInt64 := int64(accountId)
		txId, err := queries.InsertTransaction(ctx, sqlc.InsertTransactionParams{
			Account: accountIdInt64,
			Mint:    sql.NullInt64{Int64: mintId, Valid: true},
		})
		if err != nil {
			return 0, fmt.Errorf("querying InsertTransaction: %w", err)
		}

		err = queries.IncrementBalance(ctx, sqlc.IncrementBalanceParams{
			Account: accountIdInt64,
			Amount:  strconv.Itoa(amount),
		})
		if err != nil {
			return 0, fmt.Errorf("querying IncrementBalance: %w", err)
		}
		return int(txId), nil
	})
}

func (model *Model) GetTransactions(ctx context.Context, accountId int) ([]interface{}, error) {
	queries := sqlc.New(model.db)

	err := model.Exists(ctx, accountId)
	if err != nil {
		return nil, err
	}

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
	return WithTransaction(model.db, func(tx *sql.Tx) (int, error) {
		queries := sqlc.New(tx)
		err := model.Exists(ctx, accountId)
		if err != nil {
			return 0, err
		}

		err = model.HasEnough(ctx, accountId, amount)
		if err != nil {
			return 0, err
		}

		amountDecimal := strconv.Itoa(amount)
		mintId, err := queries.InsertSpend(ctx, amountDecimal)
		if err != nil {
			return 0, fmt.Errorf("query InsertSpend: %w", err)
		}

		accountIdInt64 := int64(accountId)
		txId, err := queries.InsertTransaction(ctx, sqlc.InsertTransactionParams{
			Account: accountIdInt64,
			Spend:   sql.NullInt64{Int64: mintId, Valid: true},
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

		return int(txId), nil
	})
}

func (model *Model) Transfer(ctx context.Context, senderAccountId int, recipientAccountId int, amount int) (int, error) {
	return WithTransaction(model.db, func(tx *sql.Tx) (int, error) {
		queries := sqlc.New(tx)

		err := model.Exists(ctx, senderAccountId)
		if err != nil {
			return 0, err
		}

		err = model.Exists(ctx, recipientAccountId)
		if err != nil {
			return 0, err
		}

		err = model.HasEnough(ctx, senderAccountId, amount)
		if err != nil {
			return 0, err
		}

		amountDecimal := strconv.Itoa(amount)
		transferId, err := queries.InsertTransfer(ctx, sqlc.InsertTransferParams{
			Recipient: int64(recipientAccountId),
			Amount:    amountDecimal,
		})
		if err != nil {
			return 0, fmt.Errorf("query InsertTransfer: %w", err)
		}

		txId, err := queries.InsertTransaction(ctx, sqlc.InsertTransactionParams{
			Account:  int64(senderAccountId),
			Transfer: sql.NullInt64{Int64: transferId, Valid: true},
		})
		if err != nil {
			return 0, fmt.Errorf("query InsertTransaction: %w", err)
		}

		err = queries.DecrementBalance(ctx, sqlc.DecrementBalanceParams{
			Account: int64(senderAccountId),
			Amount:  amountDecimal,
		})
		if err != nil {
			return 0, fmt.Errorf("query DecrementBalance: %w", err)
		}

		err = queries.IncrementBalance(ctx, sqlc.IncrementBalanceParams{
			Account: int64(recipientAccountId),
			Amount:  amountDecimal,
		})
		if err != nil {
			return 0, fmt.Errorf("query IncrementBalance: %w", err)
		}

		return int(txId), nil
	})
}
