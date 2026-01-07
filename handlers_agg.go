package main

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/juandrzej/postgreSQL-blog-aggregator/internal/database"
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

	var nullStringValue sql.NullString
	for _, item := range feed.Channel.Item {
		pubTime, err := time.Parse(time.RFC1123Z, item.PubDate)
		if err != nil {
			fmt.Println(err)
			continue
		}

		nullStringValue = sql.NullString{
			String: item.Description,
			Valid:  item.Description != "", // Valid is true if myString is not empty
		}

		_, err = s.db.CreatePost(context.Background(), database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			Title:       item.Title,
			Url:         item.Link,
			Description: nullStringValue,
			PublishedAt: pubTime,
			FeedID:      feedToFetch.ID,
		})
		if err != nil {
			if strings.Contains(err.Error(), "unique") {
				continue
			}
			fmt.Println(err)
		}
	}
	return nil
}

func handlerBrowse(s *state, cmd command, currentUser database.User) error {
	// This func prints the posts to terminal
	var limit int32
	var err error
	if len(cmd.arguments) == 0 {
		limit = 2
	} else {
		var parsedLimit int64
		parsedLimit, err = strconv.ParseInt(cmd.arguments[0], 10, 32)
		if err != nil {
			return err
		}
		limit = int32(parsedLimit)
	}

	posts, err := s.db.GetPostsForUser(context.Background(), database.GetPostsForUserParams{
		UserID: currentUser.ID,
		Limit:  limit,
	})
	if err != nil {
		return err
	}

	for _, post := range posts {
		fmt.Println(post)
	}

	return nil
}
