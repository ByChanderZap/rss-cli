package main

import (
	"context"
	"fmt"
)

func handlerAggregator(s *state, cmd command) error {
	rssFeed, err := fetchFeed(context.Background(), "https://www.wagslane.dev/index.xml")
	if err != nil {
		return err
	}

	fmt.Printf("result: %v", rssFeed)
	return nil
}
