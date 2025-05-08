package handlers

import (
	"context"

	"github.com/lucashthiele/gator/internal/model"
)

func HandlerReset(s *model.State, cmd model.Command) error {
	err := s.Db.ResetUsersTable(context.Background())
	if err != nil {
		return err
	}

	return nil
}
