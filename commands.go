package main

import (
	"fmt"

	"github.com/juandrzej/postgreSQL-blog-aggregator/internal/config"
)

type command struct {
	name      string
	arguments []string
}

type state struct {
	cfg *config.Config
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.arguments) != 1 {
		return fmt.Errorf("Invalid arguments number.")
	}

	if err := s.cfg.SetUser(cmd.arguments[0]); err != nil {
		return err
	}

	fmt.Println("The user has been set.")
	return nil
}

type commands struct {
	commands map[string]func(*state, command) error
}

func (c *commands) run(s *state, cmd command) error {

	return nil
}

func (c *commands) register(name string, f func(*state, command) error) {
}
