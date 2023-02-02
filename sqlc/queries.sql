-- name: GetBalance :one
SELECT * FROM balances WHERE id=$1 LIMIT 1;

-- name: RegisterAccount :one
WITH ids AS (
  INSERT INTO accounts (
    name
  ) VALUES (
    $1
  ) RETURNING id
)
INSERT INTO balances (
  id, balance
) SELECT
  id, 0
FROM ids
RETURNING id;
