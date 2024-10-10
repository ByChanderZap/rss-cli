package main

import (
	"context"
	"fmt"
	"time"

	"github.com/ByChanderZap/rss-cli/internal/database"
	"github.com/google/uuid"
)

func handlerRemoveFollow(s *state, cmd command, user database.User) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <feed_url>", cmd.Name)
	}

	feed, err := s.db.GetFeedByURL(context.Background(), cmd.Args[0])
	if err != nil {
		return fmt.Errorf("couldn't get feed: %w", err)
	}

	err = s.db.DeleteFeedFollow(context.Background(), database.DeleteFeedFollowParams{
		UserID: user.ID,
		FeedID: feed.ID,
	})

	if err != nil {
		return fmt.Errorf("error trying to remove follow: %w", err)
	}

	fmt.Println("Follow removed successfully")
	return nil

	//

	// if len(cmd.Args) != 1 {
	// 	return fmt.Errorf("usage: %s <feed_url>", cmd.Name)
	// }

	// feed, err := s.db.GetFeedByURL(context.Background(), cmd.Args[0])
	// if err != nil {
	// 	return fmt.Errorf("couldn't get feed: %w", err)
	// }

	// err = s.db.DeleteFeedFollow(context.Background(), database.DeleteFeedFollowParams{
	// 	UserID: user.ID,
	// 	FeedID: feed.ID,
	// })
	// if err != nil {
	// 	return fmt.Errorf("couldn't delete feed follow: %w", err)
	// }

	// fmt.Printf("%s unfollowed successfully!\n", feed.Name)
	// return nil
}

func handlerAddFollow(s *state, cmd command, user database.User) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <feed_url>", cmd.Name)
	}

	dbFeed, err := s.db.GetFeedByURL(context.Background(), cmd.Args[0])
	if err != nil {
		return fmt.Errorf("error fetching feed by url: %v", err)
	}

	ffRow, err := s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    user.ID,
		FeedID:    dbFeed.ID,
	})
	if err != nil {
		return fmt.Errorf("could not add feed follow: %w", err)
	}

	fmt.Println("Following feed created:")
	printFeedFollow(ffRow.UserName, ffRow.FeedName)

	return nil
}

func handlerListFeedFollows(s *state, cmd command, user database.User) error {
	followedByUser, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return fmt.Errorf("error retrieving user `%s` follows: %w", user.Name, err)
	}
	fmt.Println("DEBUG: ")
	fmt.Printf("currentUserName: %s ", s.cfg.CurrentUserName)
	if len(followedByUser) == 0 {
		fmt.Println("No feed follows found for this user.")
		return nil
	}

	fmt.Printf("Feed follows for user %s:\n", user.Name)
	for _, ff := range followedByUser {
		fmt.Printf("* %s\n", ff.FeedName)
	}
	return nil
}

func printFeedFollow(username, feedname string) {
	fmt.Printf("* User:          %s\n", username)
	fmt.Printf("* Feed:          %s\n", feedname)
}
