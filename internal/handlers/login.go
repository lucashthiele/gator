package handlers

import (
	"fmt"

	"github.com/lucashthiele/gator/internal/model"
	"github.com/lucashthiele/gator/internal/shared"
)

func HandlerLogin(s *model.State, cmd model.Command) error {
	expectedArguments := 1

	if len(cmd.Arguments) != expectedArguments {
		return fmt.Errorf(
			"expected %d arguments but got %d arguments\nyou need to provide only the username for login",
			expectedArguments,
			len(cmd.Arguments))
	}

	username := cmd.Arguments[0]

	_, err := shared.GetCurrentUser(s)
	if err != nil {
		return err
	}

	err = s.Config.SetUser(username)
	if err != nil {
		return err
	}

	fmt.Printf("user has been set!")

	return nil
}
