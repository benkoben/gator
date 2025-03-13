package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/url"
    "strconv"
	"time"

	"github.com/benkoben/gator/internal/database"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

func handlerBrowse(s *state, cmd command, user database.User) error {
    var limit int32
    if len(cmd.args) == 0 {
        limit = 2
    } else {
        // Check if the argument is valid
        l, err := strconv.Atoi(cmd.args[0]); 
        if err != nil {
            return fmt.Errorf("limit must be a number")
        }
        limit = int32(l)
    }

    params := database.GetPostsForUserParams{
       UserID: user.ID,
       Limit: limit,
    }

    ctx := context.Background()
    userPosts, err := s.db.GetPostsForUser(ctx, params)
    if err != nil {
        return fmt.Errorf("could not retrieve posts for user: %s", err)
    }

    fmt.Printf("Latest posts for %s\n", user.Name)
    for _, post := range userPosts{
        fmt.Printf("Title: %s\n", post.Title) 
        fmt.Printf("Description: %s\n", post.Description)
        fmt.Printf("Url: %s\n", post.Url)
        fmt.Println("-----------------")
    }

    return nil
}

func handlerFetchFeed(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("duration is missing: <1s, 1m, 10m, 1h>")
	}

	duration, err := time.ParseDuration(cmd.args[0])
	if err != nil {
		return fmt.Errorf("could not parse duration: %s", err)
	}

	ticker := time.NewTicker(duration)

	for ; ; <-ticker.C {
		fmt.Println("scraping next feed")
		if err := scrapeFeeds(s); err != nil {
			fmt.Println(err)
		}
	}
}

func handlerAddFeed(s *state, cmd command, user database.User) error {

	if len(cmd.args) < 2 {
		return fmt.Errorf("missing url and name: < URL > < NAME >")
	}

	feedName := cmd.args[0]
	feedUrl := cmd.args[1]

	// Validate if the URL is in a valid URL format
	if _, err := url.ParseRequestURI(feedUrl); err != nil {
		return err
	}

	ctx := context.Background()
	params := database.AddFeedParams{
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      feedName,
		Url:       feedUrl,
		UserID:    user.ID,
	}
	feed, err := s.db.AddFeed(ctx, params)
	if err != nil {
		if err.(*pq.Error).Code == "23505" {
			return fmt.Errorf("feed entry already exists in database")
		}
		return fmt.Errorf("could not add feed to database: %s", err)
	}

	fmt.Println(feed)
	return handlerFollow(s, command{args: cmd.args[1:]}, user)
}

/*
Unfollow takes a url argument and deletes all feed follow entries that match the current user ID and the feed_id that corresponds to the URL
Any database related errors are returned.
*/
func handlerUnfollow(s *state, cmd command, user database.User) error {
	if len(cmd.args[0]) == 0 {
		return fmt.Errorf("missing url: unfollow < URL >")
	}

	// Validate if the URL is in a valid URL format
	if _, err := url.ParseRequestURI(cmd.args[0]); err != nil {
		return err
	}

	ctx := context.Background()
	feed, err := s.db.LookupFeedByUrl(ctx, cmd.args[0])

	if err != nil {
		return fmt.Errorf("could not lookup feed in database: %s", err)
	}
	params := database.DeleteFeedFollowForUserParams{
		UserID: user.ID,
		FeedID: feed.ID,
	}
	if _, err := s.db.DeleteFeedFollowForUser(ctx, params); err != nil {
		return fmt.Errorf("could not delete feed follow entry: %s", err)
	}

	fmt.Printf("Unfollowed %s", feed.Name)
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
	for i := range feeds {
		fmt.Printf("* Name: %s\n", feeds[i].FeedName)
		fmt.Printf("* Url: %s\n", feeds[i].Url)
		fmt.Printf("* User: %s\n", feeds[i].Username)
		fmt.Println("-------------------")
	}

	return nil
}

/*
returns the feeds the current user is following
*/
func handlerFollowing(s *state, cmd command, user database.User) error {
	ctx := context.Background()

	follows, err := s.db.GetFeedFollowsForUser(ctx, user.ID)
	if err != nil {
		fmt.Errorf("could not lookup feed follows in database: %s", err)
	}

	fmt.Println("Following:")
	for i := range follows {
		fmt.Printf("* %s\n", follows[i].FeedName)
	}

	return nil
}

func handlerFollow(s *state, cmd command, user database.User) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("url is required")
	}

	// Validate if the URL is in a valid URL format
	if _, err := url.ParseRequestURI(cmd.args[0]); err != nil {
		return err
	}

	// Run two queries to save the requested feed and current user objects to memory
	ctx := context.Background()

	feed, err := s.db.LookupFeedByUrl(ctx, cmd.args[0])
	if err != nil {
		return fmt.Errorf("could not lookup requested feed: %s", err)
	}

	// Compose parameters and query to create a new feed follow entry
	params := database.CreateFeedFollowParams{
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	}

	follow, err := s.db.CreateFeedFollow(ctx, params)
	if err != nil {
		return fmt.Errorf("could not create new database entry: %s", err)
	}

	fmt.Printf("%s is now following %s", follow.UserName, follow.FeedName)

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
