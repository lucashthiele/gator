// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0
// source: feed_follows.sql

package database

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const createFeedFollow = `-- name: CreateFeedFollow :one
WITH INSERTED_FEED_FOLLOW AS (
  INSERT INTO FEED_FOLLOWS (
    ID,
    CREATED_AT,
    UPDATED_AT,
    USER_ID,
    FEED_ID
  ) VALUES (
    $1,
    $2,
    $3,
    $4,
    $5
  )
  RETURNING id, created_at, updated_at, user_id, feed_id
)
SELECT iff.id, iff.created_at, iff.updated_at, iff.user_id, iff.feed_id,
       F.NAME FEED_NAME,
       U.NAME USERNAME
  FROM INSERTED_FEED_FOLLOW IFF
  JOIN FEEDS F ON IFF.FEED_ID = F.ID
  JOIN USERS U ON IFF.USER_ID = U.ID
`

type CreateFeedFollowParams struct {
	ID        uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
	UserID    uuid.UUID
	FeedID    uuid.UUID
}

type CreateFeedFollowRow struct {
	ID        uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
	UserID    uuid.UUID
	FeedID    uuid.UUID
	FeedName  string
	Username  string
}

func (q *Queries) CreateFeedFollow(ctx context.Context, arg CreateFeedFollowParams) (CreateFeedFollowRow, error) {
	row := q.db.QueryRowContext(ctx, createFeedFollow,
		arg.ID,
		arg.CreatedAt,
		arg.UpdatedAt,
		arg.UserID,
		arg.FeedID,
	)
	var i CreateFeedFollowRow
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.UserID,
		&i.FeedID,
		&i.FeedName,
		&i.Username,
	)
	return i, err
}

const deleteFeedFollow = `-- name: DeleteFeedFollow :exec
DELETE FROM FEED_FOLLOWS WHERE USER_ID = $1 AND FEED_ID = $2
`

type DeleteFeedFollowParams struct {
	UserID uuid.UUID
	FeedID uuid.UUID
}

func (q *Queries) DeleteFeedFollow(ctx context.Context, arg DeleteFeedFollowParams) error {
	_, err := q.db.ExecContext(ctx, deleteFeedFollow, arg.UserID, arg.FeedID)
	return err
}

const getFeedFollowsByUser = `-- name: GetFeedFollowsByUser :many
SELECT F.NAME FEED_NAME,
       U.NAME USERNAME
  FROM FEED_FOLLOWS FF
  JOIN FEEDS F ON FF.FEED_ID = F.ID
  JOIN USERS U ON FF.USER_ID = U.ID
 WHERE U.NAME = $1
`

type GetFeedFollowsByUserRow struct {
	FeedName string
	Username string
}

func (q *Queries) GetFeedFollowsByUser(ctx context.Context, name string) ([]GetFeedFollowsByUserRow, error) {
	rows, err := q.db.QueryContext(ctx, getFeedFollowsByUser, name)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetFeedFollowsByUserRow
	for rows.Next() {
		var i GetFeedFollowsByUserRow
		if err := rows.Scan(&i.FeedName, &i.Username); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
