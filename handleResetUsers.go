package main

import (
	"context"
	"fmt"
)

func handlerResetUsers(s *state, cmd command) error {
	err := s.db.ResetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("Encountered error while deleting entries from [users] table: %v", err)
	}

	fmt.Print("Successfully delete all entries from [users] table")
	return nil
}
