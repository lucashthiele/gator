package handlers

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/lucashthiele/gator/internal/database"
	"github.com/lucashthiele/gator/internal/model"
)

func HandlerRegister(s *model.State, cmd model.Command) error {
	if len(cmd.Arguments) != 1 {
		return fmt.Errorf(
			"expected %d arguments but got %d arguments\nyou need to provide only the username for register",
			1,
			len(cmd.Arguments))
	}
	username := cmd.Arguments[0]

	user := &database.User{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      username,
	}

	createdUser, err := s.Db.CreateUser(context.Background(), database.CreateUserParams(*user))
	if err != nil {
		return err
	}

	err = s.Config.SetUser(username)
	if err != nil {
		return err
	}

	fmt.Printf("user created: %q", createdUser)
	return nil
}
