package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"

	"github.com/juandrzej/postgreSQL-blog-aggregator/internal/config"
	"github.com/juandrzej/postgreSQL-blog-aggregator/internal/database"
)

type command struct {
	name      string
	arguments []string
}

type state struct {
	db  *database.Queries
	cfg *config.Config
}

func handlerLogin(s *state, cmd command) error {
	// This func logs in given user while checking argument number and throwing error if applicable
	if len(cmd.arguments) == 0 {
		return fmt.Errorf("username is required")
	} else if len(cmd.arguments) > 1 {
		return fmt.Errorf("too many arguments, only username is required")
	}

	if err := s.cfg.SetUser(cmd.arguments[0]); err != nil {
		return err
	}

	fmt.Println("The user has been set.")
	return nil
}

func handlerRegister(s *state, cmd command) error {
	// This func registers new users in the database
	if len(cmd.arguments) == 0 {
		return fmt.Errorf("username is required")
	} else if len(cmd.arguments) > 1 {
		return fmt.Errorf("too many arguments, only username is required")
	}

	_, err := s.db.GetUser(context.Background(), cmd.arguments[0])
	if err == nil {
		os.Exit(1)
	}
	if err == sql.ErrNoRows {
	} else {
		return err
	}

	user, err := s.db.CreateUser(context.Background(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.arguments[0],
	})
	if err != nil {
		return err
	}

	s.cfg.CurrentUserName = user.Name
	fmt.Println("User was successfully created!")
	fmt.Println(user)

	return nil
}

type commands struct {
	commands map[string]func(*state, command) error
}

func (c *commands) run(s *state, cmd command) error {
	// This method runs a given command with the provided state if it exists.
	handler, ok := c.commands[cmd.name]
	if !ok {
		return fmt.Errorf("unknown command")
	}
	return handler(s, cmd)
}

func (c *commands) register(name string, f func(*state, command) error) {
	// This method registers a new handler function for a command name.
	c.commands[name] = f
}
