package model

import (
	"fmt"

	"github.com/lucashthiele/gator/internal/config"
	"github.com/lucashthiele/gator/internal/database"
)

type State struct {
	Db     *database.Queries
	Config *config.Config
}

type Command struct {
	Name      string
	Arguments []string
}

type Commands struct {
	Handlers map[string]func(*State, Command) error
}

func (c *Commands) Run(s *State, cmd Command) error {
	handler, found := c.Handlers[cmd.Name]
	if !found {
		return fmt.Errorf("Command not supported")
	}
	return handler(s, cmd)
}

func (c *Commands) Register(name string, handler func(*State, Command) error) {
	c.Handlers[name] = handler
}
