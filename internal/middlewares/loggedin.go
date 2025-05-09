package middlewares

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/lucashthiele/gator/internal/database"
	"github.com/lucashthiele/gator/internal/model"
)

func LoggedIn(handler func(s *model.State, cmd model.Command, user *database.User) error) func(*model.State, model.Command) error {
	return func(s *model.State, cmd model.Command) error {
		user, err := s.Db.GetUser(context.Background(), s.Config.CurrentUserName)
		if err == sql.ErrNoRows {
			return fmt.Errorf("user does not exists")
		} else if err != nil {
			return err
		}

		return handler(s, cmd, &user)
	}
}
