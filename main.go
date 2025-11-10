package main

import (
	"fmt"
	"log"

	"github.com/juandrzej/postgreSQL-blog-aggregator/internal/config"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatal(err)
	}

	if err := cfg.SetUser("juan"); err != nil {
		log.Fatal(err)
	}

	cfg2, err := config.Read()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(cfg2)

}

type state struct {
	cfg *config.Config
}

type command struct {
	name      string
	arguments []string
}
