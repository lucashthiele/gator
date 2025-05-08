package handlers

import (
	"context"
	"fmt"

	"github.com/lucashthiele/gator/internal/model"
)

func HandlerFeeds(s *model.State, cmd model.Command) error {
	feeds, err := s.Db.GetFeeds(context.Background())
	if err != nil {
		return err
	}

	fmt.Printf("Feeds: /n%q", feeds)

	return nil

}
