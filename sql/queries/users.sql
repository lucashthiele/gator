-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, name)
VALUES (
    $1,
    $2,
    $3,
    $4
)
RETURNING *;

-- name: GetUser :one
SELECT ID,
       CREATED_AT,
       UPDATED_AT,
       NAME
  FROM USERS
 WHERE NAME = $1;

-- name: ResetUsersTable :exec
DELETE FROM USERS;

-- name: GetUsers :many
SELECT ID,
       CREATED_AT,
       UPDATED_AT,
       NAME
  FROM USERS;