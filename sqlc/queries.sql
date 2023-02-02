-- name: GetAccount :one
SELECT * FROM accounts WHERE id=$1 LIMIT 1;

-- name: GetBalance :one
SELECT balance FROM balances WHERE account=$1 LIMIT 1;

-- name: InsertAccount :one
INSERT INTO accounts (
  name
) VALUES (
  $1
) RETURNING id;

-- name: InsertBalance :exec
INSERT INTO balances (
  account, balance
) VALUES (
  $1, 0
);

-- name: InsertMint :one
INSERT INTO mints (
  amount
) VALUES (
  $1
) RETURNING id;

-- name: InsertSpend :one
INSERT INTO spends (
  amount
) VALUES (
  $1
) RETURNING id;

-- name: InsertTransfer :one
INSERT INTO transfers (
  receipient, amount
) VALUES (
  $1, $2
) RETURNING id;

-- name: InsertTransaction :one
INSERT INTO transactions (
  account, mint, spend, transfer
) VALUES (
  $1, $2, $3, $4
) RETURNING id;

-- name: IncrementBalance :exec
UPDATE balances SET balance = balance + sqlc.arg(amount) WHERE account=$1;

-- name: DecrementBalance :exec
UPDATE balances SET balance = balance - sqlc.arg(amount) WHERE account=$1;

