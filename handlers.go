package main

import (
	"context"
	"database/sql"
	"fmt"
	"time"
    "log"

	"github.com/benkoben/gator/internal/database"
	"github.com/google/uuid"
)

func handlerLogin(s *state, cmd command) error {
    
    if len(cmd.args) == 0 {
        return fmt.Errorf("username is required")
    }
    user := cmd.args[0]
    if !usernameExists(s, user){
        return fmt.Errorf("username does not exist")
    }

    if err := s.config.SetUser(user); err != nil {
        return fmt.Errorf("could not set username to config: %s", err)
    }
    
    log.Printf("username %s has been set\n", user)
    return nil
}

func handlerRegister(s *state, cmd command) error {
    if len(cmd.args) == 0 {
        return fmt.Errorf("username is required")
    }

    if usernameExists(s, cmd.args[0]) {
        return fmt.Errorf("username already exists") 
    }

    ctx := context.Background()
    params := database.CreateUserParams{
        ID: uuid.New(),
        CreatedAt: time.Now(),
        UpdatedAt: time.Now(),
        Name: cmd.args[0],
    }
    if _, err := s.db.CreateUser(ctx, params); err != nil {
        return fmt.Errorf("failed to create user: %s", err)
    }

    log.Printf("user \"%s\" sucessfully created\n", cmd.args[0])
    return handlerLogin(s, cmd)
}

func handlerListUsers(s *state, cmd command) error {
    ctx := context.Background()
    users, err := s.db.GetUsers(ctx)
    
    if err != nil {
        return fmt.Errorf("could not get users: %s", err)
    }

    for i := range users{
        if users[i] == s.config.CurrentUsername {
            fmt.Printf("* %s (current)\n", users[i])
        } else {
            fmt.Printf("* %s\n", users[i])
        }
    }

    return nil
}

func handlerResetUsers(s *state, cmd command) error {
    return s.db.ResetUsers(context.Background())
}

// -- Helper functions
func usernameExists(s *state, name string) bool {
    ctx := context.Background()
    _, err := s.db.GetUser(ctx, name)

    if err != nil {
        if err == sql.ErrNoRows {
            return false
        }
    }

    return true 
}

