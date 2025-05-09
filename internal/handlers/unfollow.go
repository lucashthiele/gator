package handlers

import (
	"context"
	"fmt"

	"github.com/lucashthiele/gator/internal/database"
	"github.com/lucashthiele/gator/internal/model"
	"github.com/lucashthiele/gator/internal/repository"
)

func HandlerUnfollow(s *model.State, cmd model.Command, user *database.User) error {
	expectedArguments := 1

	if len(cmd.Arguments) != expectedArguments {
		return fmt.Errorf(
			"expected %d arguments but got %d arguments\nyou need to provide only the url to unfollow",
			expectedArguments,
			len(cmd.Arguments))
	}

	feedUrl := cmd.Arguments[0]

	feed, err := repository.GetFeedFromURL(s, feedUrl)
	if err != nil {
		return err
	}

	deleteFeedParams := database.DeleteFeedFollowParams{
		UserID: user.ID,
		FeedID: feed.ID,
	}

	err = s.Db.DeleteFeedFollow(context.Background(), deleteFeedParams)
	if err != nil {
		return err
	}

	return nil
}
