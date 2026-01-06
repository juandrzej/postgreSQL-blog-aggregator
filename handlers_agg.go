package main

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/juandrzej/postgreSQL-blog-aggregator/internal/database"
)

func handlerAgg(s *state, cmd command) error {
	// This func fetches feed from RSS and prints it to the console
	feed, err := FetchFeed(context.Background(), "https://www.wagslane.dev/index.xml")
	if err != nil {
		return err
	}

	fmt.Println(feed)
	return nil
}
