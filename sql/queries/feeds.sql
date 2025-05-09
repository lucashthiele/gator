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

-- name: GetFeedByURL :one
SELECT ID,
       CREATED_AT,
       UPDATED_AT,
       NAME,
       URL,
       USER_ID
  FROM FEEDS
 WHERE URL = $1;

-- name: MarkFeedFetched :exec
UPDATE FEEDS
   SET LAST_FETCHED_AT = $1,
       UPDATED_AT = $2
 WHERE ID = $3;

-- name: GetNextFeedToFetch :many
SELECT ID,
       NAME,
       URL
  FROM FEEDS
 ORDER BY LAST_FETCHED_AT ASC NULLS FIRST;