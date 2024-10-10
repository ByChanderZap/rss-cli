package main

import (
	"context"
	"fmt"
)

func reset(s *state, _ command) error {
	err := s.db.DeleteAllUsers(context.Background())
	if err != nil {
		return fmt.Errorf("couldn't delete all users: %w", err)
	}
	return nil
}
