package handlers

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/lucashthiele/gator/internal/database"
	"github.com/lucashthiele/gator/internal/model"
	"github.com/lucashthiele/gator/internal/repository"
)

func HandlerFollow(s *model.State, cmd model.Command, user *database.User) error {
	expectedArguments := 1

	if len(cmd.Arguments) != expectedArguments {
		return fmt.Errorf(
			"expected %d arguments but got %d arguments\nyou need to provide only the url to follow",
			expectedArguments,
			len(cmd.Arguments))
	}

	feedUrl := cmd.Arguments[0]

	feed, err := repository.GetFeedFromURL(s, feedUrl)
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
