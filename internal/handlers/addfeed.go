package handlers

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/lucashthiele/gator/internal/database"
	"github.com/lucashthiele/gator/internal/model"
)

func createFeed(s *model.State, feedName, feedUrl string) (database.Feed, error) {
	user, err := s.Db.GetUser(context.Background(), s.Config.CurrentUserName)
	if err != nil {
		return database.Feed{}, err
	}

	feed := &database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      feedName,
		Url:       feedUrl,
		UserID:    user.ID,
	}

	createdFeed, err := s.Db.CreateFeed(context.Background(), database.CreateFeedParams(*feed))
	if err != nil {
		return database.Feed{}, err
	}
	return createdFeed, nil
}

func createFollow(s *model.State, createdFeed database.Feed, user *database.User) error {
	feedFollow := &database.FeedFollow{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    createdFeed.ID,
	}

	_, err := s.Db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams(*feedFollow))
	if err != nil {
		return err
	}
	return nil
}

func HandlerAddFeed(s *model.State, cmd model.Command, user *database.User) error {
	expectedArguments := 2
	if len(cmd.Arguments) != expectedArguments {
		return fmt.Errorf(
			"expected %d arguments but got %d arguments\nyou need to provide the name of the feed and it's url",
			expectedArguments,
			len(cmd.Arguments))
	}

	feedName := cmd.Arguments[0]
	feedUrl := cmd.Arguments[1]

	createdFeed, err := createFeed(s, feedName, feedUrl)
	if err != nil {
		return err
	}

	err = createFollow(s, createdFeed, user)
	if err != nil {
		return err
	}

	fmt.Printf("Created Feed: \n%q", createdFeed.ID)
	return nil
}
