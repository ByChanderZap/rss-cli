package main

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/ByChanderZap/rss-cli/internal/database"
	"github.com/google/uuid"
)

func handlerAddFeed(s *state, cmd command, user database.User) error {
	if len(cmd.Args) != 2 {
		return fmt.Errorf("usage: %s <feed_name> <feed_url>", cmd.Name)
	} //cmd.Args[0] = name cmd.Args[1] = url

	dbFeed, err := s.db.CreateFeed(context.Background(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.Args[0],
		Url:       cmd.Args[1],
		UserID:    user.ID,
	})
	if err != nil {
		return fmt.Errorf("couldn't create feed: %w", err)
	}
	feedFollow, err := s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    user.ID,
		FeedID:    dbFeed.ID,
	})
	if err != nil {
		return fmt.Errorf("could not add feed follow: %w", err)
	}

	fmt.Println("Feed created successfully:")
	printFeed(dbFeed, user)
	fmt.Println()
	fmt.Println("Feed followed successfully:")
	printFeedFollow(feedFollow.UserName, feedFollow.FeedName)
	fmt.Println("=====================================")

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

func printFeed(feed database.Feed, user database.User) {
	fmt.Printf("* ID:            %s\n", feed.ID)
	fmt.Printf("* Created:       %v\n", feed.CreatedAt)
	fmt.Printf("* Updated:       %v\n", feed.UpdatedAt)
	fmt.Printf("* Name:          %s\n", feed.Name)
	fmt.Printf("* URL:           %s\n", feed.Url)
	fmt.Printf("* User:          %s\n", user.Name)
}
