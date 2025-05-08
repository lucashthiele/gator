-- name: CreateFeed :one
INSERT INTO FEEDS (
  ID,
  CREATED_AT,
  UPDATED_AT,
  NAME,
  URL,
  USER_ID
) VALUES (
  $1,
  $2,
  $3,
  $4,
  $5,
  $6
)
RETURNING *;

-- name: GetFeeds :many
SELECT F.NAME,
       F.URL,
       U.NAME USERNAME
  FROM FEEDS F
  JOIN USERS U ON F.USER_ID = U.ID;