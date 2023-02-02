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
  account, amount
) VALUES (
  $1, $2
) RETURNING id;

-- name: InsertTransfer :one
INSERT INTO transfers (
  from_account, to_account, amount
) VALUES (
  $1, $2, $3
) RETURNING id;

-- name: InsertTransaction :one
INSERT INTO transactions (
  mint, transfer
) VALUES (
  $1, $2
) RETURNING id;

-- name: IncrementBalance :exec
UPDATE balances SET balance = balance + sqlc.arg(amount) WHERE account=$1;

-- name: DecrementBalance :exec
UPDATE balances SET balance = balance - sqlc.arg(amount) WHERE account=$1;

