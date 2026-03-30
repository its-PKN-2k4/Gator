package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"

	"github.com/its-PKN-2k4/Gator/internal/config"
	"github.com/its-PKN-2k4/Gator/internal/database"
)

type state struct {
	db     *database.Queries
	cfgPtr *config.Config
}

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}

	currState := state{
		cfgPtr: &cfg,
	}

	cmds := commands{
		allCmds: make(map[string]func(*state, command) error),
	}

	db, err := sql.Open("postgres", currState.cfgPtr.DBURL)
	if err != nil {
		fmt.Printf("Encountered error while opening connection to database: %v", err)
		os.Exit(1)
	}
	dbQueries := database.New(db)

	currState.db = dbQueries

	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)
	cmds.register("reset", handlerResetUsers)
	cmds.register("users", handlerGetAllUsers)
	cmds.register("agg", handlerFetchFeed)
	cmds.register("addfeed", middlewareLoggedIn(handlerCreateFeed))
	cmds.register("feeds", handlerGetFeeds)
	cmds.register("follow", middlewareLoggedIn(handlerFollow))
	cmds.register("following", middlewareLoggedIn(handlerListFeedFollows))
	cmds.register("unfollow", middlewareLoggedIn(handlerRemoveUserFollowForFeed))

	args := os.Args
	if len(args) < 2 {
		log.Fatalf("Fewer than 2 arguments are provided. Needs at least 2 arguments")
		os.Exit(1)
	}

	cmd := command{
		name: args[1],
		args: args[2:],
	}

	err1 := cmds.run(&currState, cmd)
	if err1 != nil {
		fmt.Print(err1.Error())
		os.Exit(1)
	}

}
