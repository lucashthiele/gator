package handlers

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/lucashthiele/gator/internal/database"
	"github.com/lucashthiele/gator/internal/model"
	"github.com/lucashthiele/gator/internal/shared"
)

func getFeedFromURL(s *model.State, feedUrl string) (*database.Feed, error) {
	feed, err := s.Db.GetFeedByURL(context.Background(), feedUrl)
	if err == sql.ErrNoRows {
		return &database.Feed{}, fmt.Errorf("feed with provided url does not exist")
	} else if err != nil {
		return &database.Feed{}, err
	}
	return &feed, nil
}

func HandlerFollow(s *model.State, cmd model.Command) error {
	expectedArguments := 1

	if len(cmd.Arguments) != expectedArguments {
		return fmt.Errorf(
			"expected %d arguments but got %d arguments\nyou need to provide only the url to follow",
			expectedArguments,
			len(cmd.Arguments))
	}

	user, err := shared.GetCurrentUser(s)
	if err != nil {
		return err
	}

	feedUrl := cmd.Arguments[0]

	feed, err := getFeedFromURL(s, feedUrl)
	if err != nil {
		return err
	}

	feedFollow := &database.FeedFollow{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	}

	createdFeedFollow, err := s.Db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams(*feedFollow))
	if err != nil {
		return err
	}

	fmt.Printf("Created Feed Follow: \n%q", createdFeedFollow)

	return nil
}
