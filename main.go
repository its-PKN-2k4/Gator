package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/google/uuid"
	_ "github.com/lib/pq"

	"github.com/its-PKN-2k4/Gator/internal/config"
	"github.com/its-PKN-2k4/Gator/internal/database"
)

type state struct {
	db     *database.Queries
	cfgPtr *config.Config
}

type command struct {
	name string
	args []string
}

type commands struct {
	allCmds map[string]func(*state, command) error
}

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}
	fmt.Printf("Read config: %+v\n", cfg)

	currState := state{
		cfgPtr: &cfg,
	}

	cmds := commands{
		allCmds: make(map[string]func(*state, command) error),
	}

	if _, exist := cmds.allCmds["login"]; !exist {
		cmds.register("login", handlerLogin)
	}

	if _, exist := cmds.allCmds["register"]; !exist {
		cmds.register("register", handlerRegister)
	}

	if _, exist := cmds.allCmds["reset"]; !exist {
		cmds.register("reset", handlerResetUsers)
	}

	if _, exist := cmds.allCmds["users"]; !exist {
		cmds.register("users", handlerGetAllUsers)
	}

	if _, exist := cmds.allCmds["agg"]; !exist {
		cmds.register("agg", handlerFetchFeed)
	}

	if _, exist := cmds.allCmds["addfeed"]; !exist {
		cmds.register("addfeed", handlerCreateFeed)
	}

	args := os.Args
	if len(args) < 2 {
		log.Fatalf("Fewer than 2 arguments are provided. Needs at least 2 arguments")
		os.Exit(1)
	}

	cmd := command{
		name: args[1],
		args: args[2:],
	}

	db, err := sql.Open("postgres", currState.cfgPtr.DBURL)
	if err != nil {
		fmt.Printf("Encountered error while opening connection to database: %v", err)
		os.Exit(1)
	}
	dbQueries := database.New(db)

	currState.db = dbQueries

	err1 := cmds.run(&currState, cmd)
	if err1 != nil {
		fmt.Print(err1.Error())
		os.Exit(1)
	}

}

func (c *commands) run(s *state, cmd command) error {
	function, exists := c.allCmds[cmd.name]
	if !exists {
		return fmt.Errorf("The requested command does not exist\n")
	}

	err := function(s, cmd)
	if err != nil {
		return fmt.Errorf("Error encountered while executing requested command: %v\n", err)
	}
	return nil
}

func (c *commands) register(name string, f func(*state, command) error) {
	_, exists := c.allCmds[name]
	if !exists {
		c.allCmds[name] = f
		fmt.Printf("Successfully bind new command '%v' to its handler\n", name)
	} else {
		fmt.Printf("Cannot register command '%v' as new command since it already exists\n", name)
	}
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("This command needs 1 argument: username\n")
	}

	_, err0 := s.db.GetUser(context.Background(), cmd.args[0])
	switch err0 {
	case sql.ErrNoRows:
		return fmt.Errorf("No user with name <%v> exists to login", cmd.args[0])
	case nil:
		break
	default:
		return fmt.Errorf("Database operation malfunctioned: %v", err0)
	}

	err := s.cfgPtr.SetUser(cmd.args[0])
	if err != nil {
		return fmt.Errorf("Couldn't set current user for login: %v\n", err)
	}
	fmt.Printf("User has been set to: %v\n", s.cfgPtr.CurrentUserName)
	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("This command needs 1 argument: username\n")
	}

	_, err0 := s.db.GetUser(context.Background(), cmd.args[0])
	switch err0 {
	case nil:
		return fmt.Errorf("This username <%v> has already been registered", cmd.args[0])
	case sql.ErrNoRows:
		break
	default:
		return fmt.Errorf("Database operation malfunctioned: %v", err0)
	}

	newUser, err := s.db.CreateUser(context.Background(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.args[0],
	})
	if err != nil {
		return fmt.Errorf("Couldn't register new user with username <%v>", cmd.args[0])
	}

	s.cfgPtr.CurrentUserName = newUser.Name
	fmt.Printf("New user has been registered\n: %+v", newUser)

	err1 := s.cfgPtr.SetUser(cmd.args[0])
	if err1 != nil {
		return fmt.Errorf("Couldn't set current user for login: %v\n", err1)
	}
	fmt.Printf("User has been set to: %v\n", s.cfgPtr.CurrentUserName)
	return nil
}

func handlerResetUsers(s *state, cmd command) error {
	err := s.db.ResetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("Encountered error while deleting entries from [users] table: %v", err)
	}

	fmt.Print("Successfully delete all entries from [users] table")
	return nil
}

func handlerGetAllUsers(s *state, cmd command) error {
	users, err := s.db.GetAllUsers(context.Background())
	switch err {
	case sql.ErrNoRows:
		return fmt.Errorf("No entries exist in [users] table")
	case nil:
		break
	default:
		return fmt.Errorf("Error encountered while getting users from [users] table: %v", err)
	}

	for _, user := range users {
		if user.Name == s.cfgPtr.CurrentUserName {
			fmt.Printf("* %v (current)\n", user.Name)
		} else {
			fmt.Printf("* %v\n", user.Name)
		}
	}
	return nil
}

func handlerFetchFeed(s *state, cmd command) error {
	feed, err := fetchFeed(context.Background(), "https://www.wagslane.dev/index.xml")
	if err != nil {
		return fmt.Errorf("Error encountered while extracting XML content from given URL: %w", err)
	}

	fmt.Printf("Feed: %+v\n", feed)
	return nil
}

func handlerCreateFeed(s *state, cmd command) error {
	if len(cmd.args) < 2 {
		return fmt.Errorf("This command needs 2 argument: feed_name url\n")
	}

	currUser, err0 := s.db.GetUser(context.Background(), s.cfgPtr.CurrentUserName)
	switch err0 {
	case nil:
		break
	case sql.ErrNoRows:
		return fmt.Errorf("Current User's name: %v DOES NOT match with any entry", err0)
	default:
		return fmt.Errorf("Database operation malfunctioned: %v", err0)
	}

	feed, err := s.db.CreateFeed(context.Background(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.args[0],
		Url:       cmd.args[1],
		UserID:    currUser.ID,
	})

	if err != nil {
		return fmt.Errorf("Error encountered while creating feed: %w", err)
	}

	fmt.Printf("Created feed: %+v\n", feed)
	return nil
}
