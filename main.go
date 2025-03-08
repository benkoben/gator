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
    commands := commands{}
    commands.register("login", handlerLogin)
    args := os.Args

    // Parse user input
    if len(args) < 2 {
        log.Fatalf("At least two arguments where expected.")
    }
   
    command := command{
        name: args[0],
        args: args[1:],
    }

    // Execute the requested command
    commands.run(currentState, command)
}
