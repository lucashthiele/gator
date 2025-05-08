package main

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
	"github.com/lucashthiele/gator/internal/config"
	"github.com/lucashthiele/gator/internal/database"
	"github.com/lucashthiele/gator/internal/handlers"
	"github.com/lucashthiele/gator/internal/model"
)

func registerHandlers(cmd *model.Commands) {
	cmd.Register("login", handlers.HandlerLogin)
	cmd.Register("register", handlers.HandlerRegister)
	cmd.Register("reset", handlers.HandlerReset)
}

func getCommand(args []string) (model.Command, error) {
	if len(args) < 2 {
		return model.Command{}, fmt.Errorf("no command provided")
	}

	cmdName := args[1]
	var cmdArgs []string
	if len(args) > 2 {
		cmdArgs = args[2:]
	}

	return model.Command{
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

	dbURL := c.DbURL
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	dbQueries := database.New(db)

	s := &model.State{
		Db:     dbQueries,
		Config: c,
	}

	cmds := model.Commands{
		Handlers: make(map[string]func(*model.State, model.Command) error),
	}
	registerHandlers(&cmds)

	args := os.Args
	cmd, err := getCommand(args)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err = cmds.Run(s, cmd)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
