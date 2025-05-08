package handlers

import (
	"context"
	"fmt"

	"github.com/lucashthiele/gator/internal/database"
	"github.com/lucashthiele/gator/internal/model"
)

func printUsers(users []database.User, currentUser string) {
	for _, user := range users {
		isCurrentUser := user.Name == currentUser
		if isCurrentUser {
			fmt.Printf("* %s (current)\n", user.Name)
			continue
		}
		fmt.Printf("* %s\n", user.Name)
	}
}

func HandlerUsers(s *model.State, cmd model.Command) error {
	users, err := s.Db.GetUsers(context.Background())
	if err != nil {
		return err
	}

	printUsers(users, s.Config.CurrentUserName)

	return nil
}
