package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/lucashthiele/gator/internal/database"
	"github.com/lucashthiele/gator/internal/model"
)

func GetFeedFromURL(s *model.State, feedUrl string) (*database.Feed, error) {
	feed, err := s.Db.GetFeedByURL(context.Background(), feedUrl)
	if err == sql.ErrNoRows {
		return &database.Feed{}, fmt.Errorf("feed with provided url does not exist")
	} else if err != nil {
		return &database.Feed{}, err
	}
	return &feed, nil
}
