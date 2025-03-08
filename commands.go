package main

import (
    "fmt"
)

type command struct {
    name string
    args []string
}

type commands struct{
    registry map[string]func(*state, command) error
}

/*
registers a new handler function for a given command
*/
func (c *commands) register(name string, f func(*state, command) error) {
    c.registry[name] = f
}

/*
runs a given command if it exists in the registry
*/
func (c *commands) run(s *state, cmd command) error {
    handler, exists := c.registry[cmd.name]
    if !exists {
        return fmt.Errorf("this command has not been registered")
    }

    if err := handler(s, cmd); err != nil {
        return fmt.Errorf("%s: %v", cmd.name, err)
    }

    return nil
}
