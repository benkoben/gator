package main

import (
	"log"
	"os"

	"github.com/benkoben/gator/internal/config"
)

func main(){
    // Set current state
    currentState := &state{}
    cfg, err := config.Read()
    if err != nil {
        log.Fatalf("could not read config file: %s", err)
    }
    currentState.config = cfg
   
    // Register all necessary commands
    commands := commands{
        registry: map[string]func(*state, command) error{
            "login": handlerLogin,
        },
    }
    args := os.Args

    // Parse user input
    if len(args) < 2 {
        log.Fatalf("At least two arguments where expected.")
    }
   
    command := command{
        name: args[1],
        args: args[2:],
    }

    // Execute the requested command
    if err := commands.run(currentState, command); err != nil {
        log.Fatalf("failed to run command: %s", err)
    }
}
