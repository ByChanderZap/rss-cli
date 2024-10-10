package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/ByChanderZap/rss-cli/internal/database"
	"github.com/google/uuid"
)

func handlerAggregator(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <time between requests>", cmd.Name)
	}
	timeBetweenRequests, err := time.ParseDuration(cmd.Args[0])
	if err != nil {
		return fmt.Errorf("error while parsing duration %v", err)
	}

	fmt.Printf(" * Collecting feeds every \n%v\n", timeBetweenRequests)

	ticker := time.NewTicker(timeBetweenRequests)

	for ; ; <-ticker.C {
		fmt.Println("Starting scrapeFeeds...")
		scrapeFeeds(s)
	}
}

func scrapeFeeds(s *state) {
	feedToScrape, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		log.Println("could not get next feeds to fetch", err)
		return
	}

	processFeed(s.db, feedToScrape)

}

func processFeed(db *database.Queries, feed database.Feed) {
	log.Println("Found a feed to fetch")

	dbFeed, err := db.MarkFeedFetched(context.Background(), feed.ID)
	if err != nil {
		fmt.Println("error marking feed to fetch", err)
		return
	}

	rssFeed, err := fetchFeed(context.Background(), dbFeed.Url)
	if err != nil {
		fmt.Println("error marking feed to fetch", err)
		return
	}
	for _, item := range rssFeed.Channel.Item {
		publishedAt := sql.NullTime{}
		if t, err := time.Parse(time.RFC1123Z, item.PubDate); err == nil {
			publishedAt = sql.NullTime{
				Time:  t,
				Valid: true,
			}
		}

		_, err := db.CreatePost(context.Background(), database.CreatePostParams{
			ID:        uuid.New(),
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
			Title:     item.Title,
			Url:       item.Link,
			Description: sql.NullString{
				String: item.Description,
				Valid:  true,
			},
			PublishedAt: publishedAt,
			FeedID:      feed.ID,
		})
		if err != nil {
			if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
				continue
			}
			log.Printf("Unable to create the post: %v\n", err)
			continue
		}
		// fmt.Printf(" ~ Item title: %v\n", item.Title)
	}
	log.Printf(" \n\n* Feed `%s` collected\n * %v posts found\n", feed.Name, len(rssFeed.Channel.Item))
}
