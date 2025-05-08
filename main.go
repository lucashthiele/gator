package main

import (
	"fmt"
	"os"

	"github.com/lucashthiele/gator/internal/config"
)

type state struct {
	Config *config.Config
}

type command struct {
	Name      string
	Arguments []string
}

type commands struct {
	Handlers map[string]func(*state, command) error
}

func (c *commands) run(s *state, cmd command) error {
	handler, found := c.Handlers[cmd.Name]
	if !found {
		return fmt.Errorf("command not supported")
	}
	return handler(s, cmd)
}

func (c *commands) register(name string, handler func(*state, command) error) {
	c.Handlers[name] = handler
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.Arguments) != 1 {
		return fmt.Errorf(
			"expected %d arguments but got %d arguments\nyou need to provide only the username for login",
			1,
			len(cmd.Arguments))
	}

	username := cmd.Arguments[0]

	err := s.Config.SetUser(username)
	if err != nil {
		return err
	}

	fmt.Printf("user has been set!")

	return nil
}

func registerHandlers(cmd *commands) {
	cmd.register("login", handlerLogin)
}

func getCommand(args []string) (command, error) {
	if len(args) < 2 {
		return command{}, fmt.Errorf("no command provided")
	}

	cmdName := args[1]
	var cmdArgs []string
	if len(args) > 2 {
		cmdArgs = args[2:]
	}

	return command{
		Name:      cmdName,
		Arguments: cmdArgs,
	}, nil
}

func main() {
	c, err := config.Read()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	s := &state{
		Config: &c,
	}

	cmds := commands{
		Handlers: make(map[string]func(*state, command) error),
	}
	registerHandlers(&cmds)

	args := os.Args
	cmd, err := getCommand(args)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err = cmds.run(s, cmd)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
