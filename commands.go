package main

import (
	"fmt"

	"github.com/juandrzej/postgreSQL-blog-aggregator/internal/config"
)

type state struct {
	cfg *config.Config
}

type command struct {
	name      string
	arguments []string
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
