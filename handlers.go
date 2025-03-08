package main

import "fmt"

func handlerLogin(s *state, cmd command) error {
    
    if len(cmd.args) == 0 {
        return fmt.Errorf("username is required")
    }

    user := cmd.args[0]

    if err := s.config.SetUser(user); err != nil {
        return fmt.Errorf("could not set username to config: %s", err)
    }
    
    fmt.Printf("username %s has been set\n", user)
    return nil
}


