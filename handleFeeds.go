package main

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/its-PKN-2k4/Gator/internal/database"
)

func handlerFetchFeed(s *state, cmd command) error {
	feed, err := fetchFeed(context.Background(), "https://www.wagslane.dev/index.xml")
	if err != nil {
		return fmt.Errorf("Error encountered while extracting XML content from given URL: %w", err)
	}

	fmt.Printf("Feed: %+v\n", feed)
	return nil
}

func handlerCreateFeed(s *state, cmd command) error {
	if len(cmd.args) < 2 {
		return fmt.Errorf("This command needs 2 argument: feed_name url\n")
	}

	currUser, err0 := s.db.GetUser(context.Background(), s.cfgPtr.CurrentUserName)
	switch err0 {
	case nil:
		break
	case sql.ErrNoRows:
		return fmt.Errorf("Current User's name: %v DOES NOT match with any entry", err0)
	default:
		return fmt.Errorf("Database operation malfunctioned: %v", err0)
	}

	feed, err := s.db.CreateFeed(context.Background(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.args[0],
		Url:       cmd.args[1],
		UserID:    currUser.ID,
	})

	if err != nil {
		return fmt.Errorf("Error encountered while creating feed: %w", err)
	}

	feedFollow, err := s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    currUser.ID,
		FeedID:    feed.ID,
	})
	if err != nil {
		return fmt.Errorf("Error encounteres while creating feed follow: %w", err)
	}

	fmt.Println("Feed created successfully:")
	printFeed(feed, currUser)
	fmt.Println()
	fmt.Println("Feed followed successfully:")
	printFeedFollow(feedFollow.UserName, feedFollow.FeedName)
	fmt.Println("=====================================")
	return nil
}

func handlerGetFeeds(s *state, cmd command) error {
	feedsList, err := s.db.GetFeeds(context.Background())
	switch err {
	case sql.ErrNoRows:
		return fmt.Errorf("No entries exist in [feeds] table")
	case nil:
		break
	default:
		return fmt.Errorf("Error encountered while getting feeds from [feeds] table: %v", err)
	}

	for _, feed := range feedsList {
		fmt.Printf("%+v\n", feed)
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
