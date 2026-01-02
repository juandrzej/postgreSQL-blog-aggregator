package main

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/juandrzej/postgreSQL-blog-aggregator/internal/database"
)

func handlerLogin(s *state, cmd command) error {
	// This func logs in given user while checking argument number and throwing error if applicable
	if len(cmd.arguments) == 0 {
		return fmt.Errorf("username is required")
	} else if len(cmd.arguments) > 1 {
		return fmt.Errorf("too many arguments, only username is required")
	}

	_, err := s.db.GetUser(context.Background(), cmd.arguments[0])
	if err == sql.ErrNoRows {
		return fmt.Errorf("user '%s' doesn't exist in the database", cmd.arguments[0])
	}
	if err != nil {
		return err
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
		return fmt.Errorf("user '%s' already exists", cmd.arguments[0])
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

	if err := s.cfg.SetUser(cmd.arguments[0]); err != nil {
		return err
	}
	fmt.Println("User was successfully created!")
	fmt.Println(user)

	return nil
}

func handlerReset(s *state, cmd command) error {
	// This func deletes all users from the database
	err := s.db.DeleteAllUsers(context.Background())
	if err != nil {
		return fmt.Errorf("failed to reset database: %w", err)
	}

	fmt.Println("Database has been reset successfully!")
	return nil
}

func handlerUsers(s *state, cmd command) error {
	// This func lists all users from the database
	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return err
	}

	for _, user := range users {
		if user.Name != s.cfg.CurrentUserName {
			fmt.Printf("* %s\n", user.Name)
		} else {
			fmt.Printf("* %s (current)\n", user.Name)
		}
	}
	return nil
}

func handlerAgg(s *state, cmd command) error {
	// This func fetches feed from RSS and prints it to the console
	feed, err := FetchFeed(context.Background(), "https://www.wagslane.dev/index.xml")
	if err != nil {
		return err
	}

	fmt.Println(feed)
	return nil
}

func handlerAddFeed(s *state, cmd command) error {
	// This func adds feed to feeds table
	if len(cmd.arguments) < 2 {
		return fmt.Errorf("feed name and url are both required")
	} else if len(cmd.arguments) > 2 {
		return fmt.Errorf("too many arguments, only feed name and url are required")
	}

	current_user, err := s.db.GetUser(context.Background(), s.cfg.GetCurrentUserName())
	if err != nil {
		return err
	}

	feed, err := s.db.CreateFeed(context.Background(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.arguments[0],
		Url:       cmd.arguments[1],
		UserID:    current_user.ID,
	})
	if err != nil {
		return err
	}

	fmt.Println("Feed was successfully created!")
	fmt.Println(feed)

	return nil
}
