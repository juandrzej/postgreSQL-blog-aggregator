package main

import (
	"context"
	"fmt"
	"time"
)

func handlerAgg(s *state, cmd command) error {
	// This function reads duration and calls scrapeFeeds

	if len(cmd.arguments) == 0 {
		return fmt.Errorf("duration is required")
	} else if len(cmd.arguments) > 1 {
		return fmt.Errorf("too many arguments, only duration is required")
	}

	timeBetweenRequests, err := time.ParseDuration(cmd.arguments[0])
	if err != nil {
		return err
	}

	ticker := time.NewTicker(timeBetweenRequests)
	for ; ; <-ticker.C {
		if err := scrapeFeeds(s); err != nil {
			fmt.Println("error scraping feeds:", err)
		}
	}
}

func scrapeFeeds(s *state) error {
	// This func fetches feed from RSS and prints it to the console
	feedToFetch, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		return err
	}

	feed, err := fetchFeed(context.Background(), feedToFetch.Url)
	if err != nil {
		return err
	}

	err = s.db.MarkFeedFetched(context.Background(), feedToFetch.ID)
	if err != nil {
		return err
	}

	for _, item := range feed.Channel.Item {
		fmt.Println(item.Title)
	}
	return nil
}
