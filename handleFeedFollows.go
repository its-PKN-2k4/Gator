package main

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/its-PKN-2k4/Gator/internal/database"
)

func handlerFollow(s *state, cmd command, user database.User) error {

	if len(cmd.args) != 1 {
		return fmt.Errorf("Usage: %s <feed_url>", cmd.name)
	}

	feed, err := s.db.GetFeedByUrl(context.Background(), cmd.args[0])
	if err != nil {
		return fmt.Errorf("Couldn't get feed: %w", err)
	}

	ffRow, err := s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	})
	if err != nil {
		return fmt.Errorf("Couldn't create feed follow: %w", err)
	}

	fmt.Println("Feed follow created:")
	printFeedFollow(ffRow.UserName, ffRow.FeedName)
	return nil
}

func handlerListFeedFollows(s *state, cmd command, user database.User) error {

	feedFollows, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return fmt.Errorf("Couldn't get feed follows: %w", err)
	}

	if len(feedFollows) == 0 {
		fmt.Println("No feed follows found for this user.")
		return nil
	}

	fmt.Printf("Feed follows for user %s:\n", user.Name)
	for _, ff := range feedFollows {
		fmt.Printf("* %s\n", ff.FeedName)
	}

	return nil
}

func printFeedFollow(username, feedname string) {
	fmt.Printf("* User:          %s\n", username)
	fmt.Printf("* Feed:          %s\n", feedname)
}

func handlerRemoveUserFollowForFeed(s *state, cmd command, user database.User) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("Usage: %s <feed_url>", cmd.name)
	}

	feed, err := s.db.GetFeedByUrl(context.Background(), cmd.args[0])
	if err != nil {
		return fmt.Errorf("Couldn't get feed: %w", err)
	}

	err0 := s.db.RemoveUserFollowForFeed(context.Background(), database.RemoveUserFollowForFeedParams{
		UserID: user.ID,
		FeedID: feed.ID,
	})
	if err0 != nil {
		return fmt.Errorf("Error encountered while removing feed follow: %v", err0)
	}
	fmt.Print("Successfully remove user follow for provided feed")
	return nil
}
