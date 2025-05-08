package shared

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/lucashthiele/gator/internal/database"
	"github.com/lucashthiele/gator/internal/model"
)

func GetCurrentUser(s *model.State) (*database.User, error) {
	user, err := s.Db.GetUser(context.Background(), s.Config.CurrentUserName)
	if err == sql.ErrNoRows {
		return &database.User{}, fmt.Errorf("user does not exists")
	} else if err != nil {
		return &database.User{}, err
	}

	return &user, nil
}
