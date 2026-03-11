package main

import (
	"blog-aggregator/internal/config"
	"blog-aggregator/internal/database"
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

type state struct {
	Cfg *config.Config
	db  *database.Queries
}
type command struct {
	name      string
	arguments []string
}

type commands struct {
	cmdlist map[string]func(*state, command) error
}

func (c *commands) run(s *state, cmd command) error {
	function, ok := c.cmdlist[cmd.name]
	if !ok {
		return fmt.Errorf("no such as command exist")
	}

	if err := function(s, cmd); err != nil {
		return err
	}

	fmt.Printf("\nExecuting '%s' successfully.\n", cmd.name)
	return nil
}

func (c *commands) register(name string, f func(*state, command) error) error {
	c.cmdlist[name] = f

	fmt.Printf("  new function '%s' is added to command list\n", name)
	return nil
}

func initCmds(cmds commands) {

	cmds.register("login", handlerLogins)
	cmds.register("register", handleRegister)
	cmds.register("reset", handlerReset)
	cmds.register("users", handlerUsers)
	cmds.register("agg", handlerAggregator)
	cmds.register("addfeed", middlewareLoggedIn(handlerAddFeed))
	cmds.register("feeds", handlerFeeds)
	cmds.register("follow", middlewareLoggedIn(handlerFollow))
	cmds.register("following", middlewareLoggedIn(handlerFollowing))
	cmds.register("unfollow", middlewareLoggedIn(handlerUnfollow))
	cmds.register("browse", middlewareLoggedIn(handlerBrowse))

	fmt.Println("")
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("not enough arguments were provided")
		os.Exit(1)
	}

	cfg, err := config.Read()
	if err != nil {
		panic(err)
	}
	db, err := sql.Open("postgres", cfg.DbUrl)
	if err != nil {
		log.Fatal(err)
	}

	dbQueries := database.New(db)

	states := state{Cfg: cfg, db: dbQueries}
	cmds := commands{
		cmdlist: make(map[string]func(*state, command) error),
	}

	initCmds(cmds)

	toRun := command{name: os.Args[1], arguments: os.Args[2:]}
	status := cmds.run(&states, toRun)
	if status != nil {
		fmt.Println("error occured: ", status)
		os.Exit(1)
	} else {
		fmt.Println("program exiting with 0 code")
		os.Exit(0)
	}
}
