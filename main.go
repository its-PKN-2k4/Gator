package main

import (
	"fmt"
	"log"
	"os"

	"github.com/its-PKN-2k4/Gator/internal/config"
)

type state struct {
	cfgPtr *config.Config
}

type command struct {
	name string
	args []string
}

type commands struct {
	allCmds map[string]func(*state, command) error
}

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}
	fmt.Printf("Read config: %+v\n", cfg)

	currState := state{
		cfgPtr: &cfg,
	}

	cmds := commands{
		allCmds: make(map[string]func(*state, command) error),
	}

	if _, exist := cmds.allCmds["login"]; !exist {
		cmds.register("login", handlerLogin)
	}

	args := os.Args
	if len(args) < 2 {
		log.Fatalf("Fewer than 2 arguments are provided. Needs at least 2 arguments")
		os.Exit(1)
	}

	cmd := command{
		name: args[1],
		args: args[2:],
	}

	err1 := cmds.run(&currState, cmd)
	if err1 != nil {
		fmt.Print(err1.Error())
		os.Exit(1)
	}
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("This command needs 1 argument: username\n")
	}

	err := s.cfgPtr.SetUser(cmd.args[0])
	if err != nil {
		return fmt.Errorf("Couldn't set current user for login: %v\n", err)
	}
	fmt.Printf("User has been set to: %v\n", s.cfgPtr.CurrentUserName)
	return nil
}

func (c *commands) run(s *state, cmd command) error {
	function, exists := c.allCmds[cmd.name]
	if !exists {
		return fmt.Errorf("The requested command does not exist\n")
	}

	err := function(s, cmd)
	if err != nil {
		return fmt.Errorf("Error encountered while executing requested command: %v\n", err)
	}
	return nil
}

func (c *commands) register(name string, f func(*state, command) error) {
	_, exists := c.allCmds[name]
	if !exists {
		c.allCmds[name] = f
		fmt.Printf("Successfully bind new command '%v' to its handler\n", name)
	} else {
		fmt.Printf("Cannot register command '%v' as new command since it already exists\n", name)
	}
}
