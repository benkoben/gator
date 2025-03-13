package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/benkoben/gator/internal/database"
	"github.com/benkoben/gator/internal/rss"
	"github.com/lib/pq"
)


func scrapeFeeds(s *state) error {
    ctx := context.Background()
    next, err := s.db.GetNextFeedToFetch(ctx)
    if err != nil {
        return err
    }

    if err := s.db.MarkFeedFetched(ctx, next.ID); err != nil {
        return fmt.Errorf("failed to mark feed as fetched: s", err)
    }

    feed, err := rss.GetFeed(next.Url)
    if err != nil {
        return fmt .Errorf("could not fetch %s: %s", next.Url, err)
    }

    for _, item := range feed.Channel.Item {
        
        pubDate, err := time.Parse(time.RFC1123Z, item.PubDate)

        if err != nil {
            // Log the error but dont exit
            log.Printf("could not parse publish date: %s", err)
            continue
        }

        ctx := context.Background()
        feed, err := s.db.LookupFeedByUrl(ctx, next.Url)
        if err != nil {
            log.Printf("could not lookup feed by url(%s): %s", item.Link, err)
            continue
        }
        params := database.CreatePostParams{
            CreatedAt: time.Now(),
            UpdatedAt: time.Now(),
            Title: item.Title,
            Description: item.Description,
            Url: item.Link,
            PublishedAt: pubDate,
            FeedID: feed.ID,
        }

        if _, err := s.db.CreatePost(ctx, params); err != nil {
            if err.(*pq.Error).Code != "23505" {
                log.Printf("could not add post to database: %s", err)
            }
            continue
        }

        log.Printf("added post %s %s to database", item.Title, item.Link)
    }

    return nil
}
