package handlers

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/lucashthiele/gator/internal/model"
)

func HandlerLogin(s *model.State, cmd model.Command) error {
	if len(cmd.Arguments) != 1 {
		return fmt.Errorf(
			"expected %d arguments but got %d arguments\nyou need to provide only the username for login",
			1,
			len(cmd.Arguments))
	}

	username := cmd.Arguments[0]

	_, err := s.Db.GetUser(context.Background(), username)
	if err == sql.ErrNoRows {
		return fmt.Errorf("user does not exists")
	} else if err != nil {
		return err
	}

	err = s.Config.SetUser(username)
	if err != nil {
		return err
	}

	fmt.Printf("user has been set!")

	return nil
}
