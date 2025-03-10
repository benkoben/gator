package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/benkoben/gator/internal/config"
	"github.com/benkoben/gator/internal/database"

	_ "github.com/lib/pq"
)

func main() {
	// Set current state
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("could not read config file: %s", err)
	}

	// Create database connection
	db, err := sql.Open("postgres", cfg.DbUrl)
	if err != nil {
		log.Fatalf("could not connect to database: %s", err)
	}
	dbQueries := database.New(db)

	currentState := &state{
		db:     dbQueries,
		config: cfg,
	}

	// Register all necessary commands
	commands := commands{
		registry: map[string]func(*state, command) error{
			"login":    handlerLogin,
			"register": handlerRegister,
			"reset":    handlerResetUsers,
			"users":    handlerListUsers,
			"agg":      handlerFetchFeed,
            "addfeed":  handlerAddFeed,
            "feeds":    handlerFeeds,
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
