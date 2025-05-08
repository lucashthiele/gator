package handlers

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/lucashthiele/gator/internal/database"
	"github.com/lucashthiele/gator/internal/model"
)

func HandlerAddFeed(s *model.State, cmd model.Command) error {
	expectedArguments := 2
	if len(cmd.Arguments) != expectedArguments {
		return fmt.Errorf(
			"expected %d arguments but got %d arguments\nyou need to provide the name of the feed and it's url",
			expectedArguments,
			len(cmd.Arguments))
	}

	feedName := cmd.Arguments[0]
	feedUrl := cmd.Arguments[1]

	user, err := s.Db.GetUser(context.Background(), s.Config.CurrentUserName)
	if err != nil {
		return err
	}

	feed := &database.Feed{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      feedName,
		Url:       feedUrl,
		UserID:    user.ID,
	}

	createdFeed, err := s.Db.CreateFeed(context.Background(), database.CreateFeedParams(*feed))
	if err != nil {
		return err
	}

	fmt.Printf("Created Feed: \n%q", createdFeed)

	return nil
}
