package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/url"
	"time"

	"github.com/benkoben/gator/internal/database"
	"github.com/benkoben/gator/internal/rss"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

func handlerFetchFeed(s *state, cmd command) error {
	sampleFeed := "https://www.wagslane.dev/index.xml"

	obj, err := rss.GetFeed(sampleFeed)
	if err != nil {
		return fmt.Errorf("failed to fetch feed: %s", err)
	}
	fmt.Println(obj)

	return nil
}

func handlerAddFeed(s *state, cmd command) error {
     
	if len(cmd.args) < 2 {
        return fmt.Errorf("missing url and name: < URL > < NAME >")
	}
    
    feedName := cmd.args[0]
    feedUrl := sql.NullString{
        String: cmd.args[1],
        Valid: true,
    }

    // Validate if the URL is in a valid URL format
    if _, err := url.ParseRequestURI(feedUrl.String); err != nil {
        return err
    }

    ctx := context.Background()
    currentUser, err := s.db.GetUser(ctx, s.config.CurrentUsername)
    if err != nil {
        return fmt.Errorf("cannot fetch feed, have you logged in?")
    }

    params := database.AddFeedParams{
        CreatedAt: time.Now(),
        UpdatedAt: time.Now(),
        Name: feedName,
        Url: feedUrl, 
        UserID: currentUser.ID,
    }
    feed, err :=  s.db.AddFeed(ctx, params)
    if err != nil {    
        if err.(*pq.Error).Code == "23505" {
            return fmt.Errorf("feed entry already exists in database")
        }
        return fmt.Errorf("could not add feed to database: %s", err)
    }
    
    fmt.Println(feed)
    return nil
}

/*
    Retrieves all the feeds in the database that are linked to the current user
*/
func handlerFeeds(s *state, cmd command) error {
    ctx := context.Background()

    feeds, err := s.db.GetFeeds(ctx)
    if err != nil {
        return fmt.Errorf("no rss feeds found for current user")
    }
    
    fmt.Println("Feeds")
    for i := range feeds{
        fmt.Printf("* Name: %s\n", feeds[i].FeedName)
        fmt.Printf("* Url: %s\n", feeds[i].Url.String)
        fmt.Printf("* User: %s\n", feeds[i].Username)
        fmt.Println("-------------------")
    } 

    return nil
}

func handlerLogin(s *state, cmd command) error {

	if len(cmd.args) == 0 {
		return fmt.Errorf("username is required")
	}
	user := cmd.args[0]
	if !usernameExists(s, user) {
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
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.args[0],
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

	for i := range users {
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

