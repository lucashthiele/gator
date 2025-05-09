package handlers

import (
	"context"
	"fmt"

	"github.com/lucashthiele/gator/internal/database"
	"github.com/lucashthiele/gator/internal/model"
)

func printFeeds(f []database.GetFeedFollowsByUserRow) {
	for _, feed := range f {
		fmt.Printf(" - %s\n", feed.FeedName)
	}
}

func HandlerFollowing(s *model.State, cmd model.Command, user *database.User) error {
	feeds, err := s.Db.GetFeedFollowsByUser(context.Background(), user.Name)
	if err != nil {
		return err
	}

	printFeeds(feeds)

	return nil
}
