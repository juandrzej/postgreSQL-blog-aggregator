package main

import (
	"fmt"

	"github.com/juandrzej/postgreSQL-blog-aggregator/internal/config"
)

func main() {
	config, err := config.Read()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(config)

}
