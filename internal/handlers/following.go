package handlers

import (
	"context"
	"fmt"

	"github.com/lucashthiele/gator/internal/database"
	"github.com/lucashthiele/gator/internal/model"
	"github.com/lucashthiele/gator/internal/shared"
)

func printFeeds(f []database.GetFeedFollowsByUserRow) {
	for _, feed := range f {
		fmt.Printf(" - %s\n", feed.FeedName)
	}
}

func HandlerFollowing(s *model.State, cmd model.Command) error {
	user, err := shared.GetCurrentUser(s)
	if err != nil {
		return err
	}

	feeds, err := s.Db.GetFeedFollowsByUser(context.Background(), user.Name)
	if err != nil {
		return err
	}

	printFeeds(feeds)

	return nil
}
