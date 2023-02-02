-- name: GetBalance :one
SELECT * FROM balances WHERE account=$1 LIMIT 1;

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

-- name: InsertTransaction :one
INSERT INTO transactions (
  mint, transfer
) VALUES (
  $1, $2
) RETURNING id;

-- name: IncrementBalance :exec
UPDATE balances SET balance = balance + sqlc.arg(amount) WHERE account=$1;

