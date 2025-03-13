package main

import (
    "context"
    "fmt"
    "github.com/benkoben/gator/internal/database"
)

func middlewareLoggedIn(handler func(s *state, cmd command, u database.User) error) func(*state, command) error {

   return func(s *state, c command) error {
        ctx := context.Background()

        user, err := s.db.GetUser(ctx, s.config.CurrentUsername);
        if err != nil {
            return fmt.Errorf("Not logged in.")
        }

        if err := handler(s, c, user); err != nil {
            return fmt.Errorf("could not run handler: %s", err)
        }
        
        return nil

   }
}

