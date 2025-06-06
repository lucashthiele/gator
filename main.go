package main

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
	"github.com/lucashthiele/gator/internal/config"
	"github.com/lucashthiele/gator/internal/database"
	"github.com/lucashthiele/gator/internal/handlers"
	"github.com/lucashthiele/gator/internal/middlewares"
	"github.com/lucashthiele/gator/internal/model"
)

func registerHandlers(cmd *model.Commands) {
	cmd.Register("login", handlers.HandlerLogin)
	cmd.Register("register", handlers.HandlerRegister)
	cmd.Register("reset", handlers.HandlerReset)
	cmd.Register("users", handlers.HandlerUsers)
	cmd.Register("agg", handlers.HandlerAgg)
	cmd.Register("addfeed", middlewares.LoggedIn(handlers.HandlerAddFeed))
	cmd.Register("feeds", handlers.HandlerFeeds)
	cmd.Register("follow", middlewares.LoggedIn(handlers.HandlerFollow))
	cmd.Register("following", middlewares.LoggedIn(handlers.HandlerFollowing))
	cmd.Register("unfollow", middlewares.LoggedIn(handlers.HandlerUnfollow))
	cmd.Register("browse", middlewares.LoggedIn(handlers.HandlerBrowse))
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
