package main

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/juandrzej/postgreSQL-blog-aggregator/internal/database"
)

func middlewareLoggedIn(
	handler func(s *state, cmd command, user database.User) error,
) func(*state, command) error {

	return func(s *state, cmd command) error {
		currentUser, err := s.db.GetUser(context.Background(), s.cfg.GetCurrentUserName())
		if err != nil {
			return err
		}

		return handler(s, cmd, currentUser)
	}
}

func handlerAddFeed(s *state, cmd command, currentUser database.User) error {
	// This func adds feed to feeds table
	if len(cmd.arguments) < 2 {
		return fmt.Errorf("feed name and url are both required")
	} else if len(cmd.arguments) > 2 {
		return fmt.Errorf("too many arguments, only feed name and url are required")
	}

	feed, err := s.db.CreateFeed(context.Background(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.arguments[0],
		Url:       cmd.arguments[1],
		UserID:    currentUser.ID,
	})
	if err != nil {
		return err
	}
	feedFollow, err := s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    currentUser.ID,
		FeedID:    feed.ID,
	})
	if err != nil {
		return err
	}

	fmt.Println("Feed and feed follow was successfully created!")
	fmt.Printf("Feed %s followed by %s\n", feedFollow.FeedName, feedFollow.UserName)

	return nil
}

func handlerFeeds(s *state, cmd command) error {
	// This func lists all feeds from the database
	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return err
	}

	for _, feed := range feeds {
		feedUser, err := s.db.GetUserByID(context.Background(), feed.UserID)
		if err != nil {
			return err
		}
		fmt.Printf("* %s\n", feed.Name)
		fmt.Printf("* %s\n", feed.Url)
		fmt.Printf("* %s\n", feedUser.Name)
	}
	return nil
}

func handlerFollow(s *state, cmd command, currentUser database.User) error {
	// This func takes a single url argument and creates a new feed follow record for the current user
	if len(cmd.arguments) == 0 {
		return fmt.Errorf("url is required")
	} else if len(cmd.arguments) > 1 {
		return fmt.Errorf("too many arguments, only url is required")
	}

	feed, err := s.db.GetFeed(context.Background(), cmd.arguments[0])
	if err != nil {
		return err
	}

	feedFollow, err := s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    currentUser.ID,
		FeedID:    feed.ID,
	})
	if err != nil {
		return err
	}

	fmt.Println("Feed follow was successfully created!")
	fmt.Printf("Feed %s followed by %s\n", feedFollow.FeedName, feedFollow.UserName)

	return nil
}

func handlerFollowing(s *state, cmd command, currentUser database.User) error {
	// This func print all the names of the feeds the current user is following
	feedFollows, err := s.db.GetFeedFollowsForUser(context.Background(), currentUser.ID)
	if err != nil {
		return err
	}

	fmt.Printf("Current user - %s is following these feeds:\n", currentUser.Name)
	for _, feedFollow := range feedFollows {
		fmt.Printf("* %s\n", feedFollow.FeedName)
	}
	return nil
}

func handlerUnfollow(s *state, cmd command, currentUser database.User) error {
	// This func accepts a feed's URL as an argument and unfollows it for the current user
	if len(cmd.arguments) == 0 {
		return fmt.Errorf("url is required")
	} else if len(cmd.arguments) > 1 {
		return fmt.Errorf("too many arguments, only url is required")
	}

	feed, err := s.db.GetFeed(context.Background(), cmd.arguments[0])
	if err != nil {
		return err
	}

	err = s.db.DeleteFeedFollow(context.Background(), database.DeleteFeedFollowParams{
		UserID: currentUser.ID,
		FeedID: feed.ID,
	})
	if err != nil {
		return err
	}

	fmt.Println("Feed follow was successfully unfollowed!")

	return nil
}
