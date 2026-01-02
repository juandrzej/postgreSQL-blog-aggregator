package main

import _ "github.com/lib/pq"

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/juandrzej/postgreSQL-blog-aggregator/internal/config"
	"github.com/juandrzej/postgreSQL-blog-aggregator/internal/database"
)

func main() {
	// start state and config
	st := state{}
	cfg, err := config.Read()
	if err != nil {
		log.Fatal(err)
	}
	st.cfg = &cfg

	// postgreSQL
	db, err := sql.Open("postgres", cfg.DbURL)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	dbQueries := database.New(db)
	st.db = dbQueries

	// register all commands
	cmds := commands{
		commands: make(map[string]func(*state, command) error),
	}
	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)
	cmds.register("reset", handlerReset)
	cmds.register("users", handlerUsers)
	cmds.register("agg", handlerAgg)
	cmds.register("addfeed", handlerAddFeed)

	// arrange arguments given in console
	arguments := os.Args
	if len(arguments) < 2 {
		fmt.Fprintln(os.Stderr, fmt.Errorf("not enough arguments"))
		os.Exit(1)
	}
	cmd := command{
		name:      arguments[1],
		arguments: arguments[2:],
	}

	// run command and output any errors
	err = cmds.run(&st, cmd)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
