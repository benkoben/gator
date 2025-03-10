package main

import (
	"github.com/benkoben/gator/internal/config"
	"github.com/benkoben/gator/internal/database"
)

type state struct {
    db *database.Queries
    config *config.Config
}
