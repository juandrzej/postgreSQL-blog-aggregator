package main

import (
	"fmt"

	"github.com/juandrzej/postgreSQL-blog-aggregator/internal/config"
	"github.com/juandrzej/postgreSQL-blog-aggregator/internal/database"
)

type state struct {
	db  *database.Queries
	cfg *config.Config
}

type command struct {
	name      string
	arguments []string
}

type commands struct {
	commands map[string]func(*state, command) error
}

func (c *commands) register(name string, f func(*state, command) error) {
	// This method registers a new handler function for a command name.
	c.commands[name] = f
}

func (c *commands) run(s *state, cmd command) error {
	// This method runs a given command with the provided state if it exists.
	handler, ok := c.commands[cmd.name]
	if !ok {
		return fmt.Errorf("unknown command")
	}
	return handler(s, cmd)
}
