package main

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/its-PKN-2k4/Gator/internal/database"
)

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
