package main

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/ByChanderZap/rss-cli/internal/database"
	"github.com/google/uuid"
)

func handlerAddFeed(s *state, cmd command) error {
	if len(cmd.Args) != 2 {
		return fmt.Errorf("usage: %s <feed_name> <feed_url>", cmd.Name)
	} //cmd.Args[0] = name cmd.Args[1] = url

	u, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
	if err != nil {
		return fmt.Errorf("couldn't find user: %w", err)
	}

	fmt.Printf("DEBUG: apparent current user %v", u)

	dbFeed, err := s.db.CreateFeed(context.Background(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.Args[0],
		Url:       cmd.Args[1],
		UserID:    u.ID,
	})
	if err != nil {
		return fmt.Errorf("couldn't create feed: %w", err)
	}
	fmt.Printf("Feed created: %v\n", dbFeed)

	_, err = s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    u.ID,
		FeedID:    dbFeed.ID,
	})
	if err != nil {
		return fmt.Errorf("could not add feed follow: %w", err)
	}
	fmt.Printf("\n\n%s now follows the feed %s\n", s.cfg.CurrentUserName, dbFeed.Name)

	return nil
}

func handlerGetFeeds(s *state, _ command) error {
	dbFeeds, err := s.db.GetFeedsPopulated(context.Background())
	if err != nil {
		return fmt.Errorf("error getting feeds: %v", err)
	}

	fmt.Println(" ---- Feeds data: ---")
	for _, feed := range dbFeeds {

		fmt.Printf(" * name: %v\n", feed.Feedname)
		fmt.Printf(" * url: %v\n", feed.Url)
		fmt.Printf(" * feed creator: %v\n", feed.Username)
		fmt.Println(strings.Repeat("*", 48))
	}
	return nil
}
