package main

import (
	"blog-aggregator/internal/config"
	"fmt"
	"os"
)

type state struct {
	Cfg *config.Config
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

	fmt.Printf("Executing %s successfully.\n", cmd.name)
	return nil
}

func (c *commands) register(name string, f func(*state, command) error) error {
	c.cmdlist[name] = f

	fmt.Printf("new function '%s' is added to command list\n", name)
	return nil
}

func handlerLogins(s *state, cmd command) error {
	if cmd.arguments == nil {
		return fmt.Errorf("not enough arguments were provided")
	} else if len(cmd.arguments) < 1 {
		return fmt.Errorf("username argument is required")
	}

	s.Cfg.Username = cmd.arguments[0]
	if err := config.Write(s.Cfg); err != nil {
		return err
	}

	fmt.Println("User has been set to", cmd.arguments[0])

	return nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("not enough arguments were provided")
		os.Exit(1)
	} else if len(os.Args) == 2 && os.Args[1] == "login" {
		fmt.Println("username is required")
		os.Exit(1)
	}

	cfg, err := config.Read()
	if err != nil {
		panic(err)
	}
	states := state{Cfg: cfg}
	cmds := commands{
		cmdlist: make(map[string]func(*state, command) error),
	}
	cmds.register("login", handlerLogins)

	toRun := command{name: os.Args[1], arguments: os.Args[2:]}
	cmds.run(&states, toRun)
}
