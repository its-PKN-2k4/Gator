package main

import (
	"context"
	"fmt"

	"github.com/its-PKN-2k4/Gator/internal/database"
)

func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {
	return func(s *state, c command) error {
		currentLoggedIn, err := s.db.GetUser(context.Background(), s.cfgPtr.CurrentUserName)
		if err != nil {
			return fmt.Errorf("Error encountered: %v", err)
		}
		return handler(s, c, currentLoggedIn)
	}
}
