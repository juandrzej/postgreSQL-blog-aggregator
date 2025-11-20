package main

import (
	"fmt"
	"log"
	"os"

	"github.com/juandrzej/postgreSQL-blog-aggregator/internal/config"
)

func main() {
	st := state{}
	cfg, err := config.Read()
	if err != nil {
		log.Fatal(err)
	}
	st.cfg = &cfg

	cmds := commands{
		commands: make(map[string]func(*state, command) error),
	}

	cmds.register("login", handlerLogin)
	arguments := os.Args
	if len(arguments) < 2 {
		fmt.Fprintln(os.Stderr, fmt.Errorf("not enough arguments"))
		os.Exit(1)
	}

	cmd := command{
		name:      arguments[1],
		arguments: arguments[2:],
	}

	err = cmds.run(&st, cmd)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
