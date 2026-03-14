package main

import (
	"fmt"
)

type command struct {
	name string
	args []string
}

type commands struct {
	allCmds map[string]func(*state, command) error
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
