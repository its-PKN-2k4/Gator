package main

import (
	"context"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"net/http"
)

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	client := http.Client{}
	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, http.NoBody)
	if err != nil {
		return nil, fmt.Errorf("Error encountered while creating netwrok request: %v\n", err)
	}
	req.Header.Set("User-Agent", "gator")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Error encountered while sending network request: %v\n", err)
	}
	defer resp.Body.Close()

	stream, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Error encountered while reading request response: %v\n", err)
	}

	channel := RSSFeed{}
	err1 := xml.Unmarshal(stream, &channel)
	if err1 != nil {
		return nil, fmt.Errorf("Error encountered while unmarshaling XML content: %v\n", err1)
	}

	channel.Channel.Title = html.UnescapeString(channel.Channel.Title)
	channel.Channel.Description = html.UnescapeString(channel.Channel.Description)
	for i, item := range channel.Channel.Item {
		item.Title = html.UnescapeString(item.Title)
		item.Description = html.UnescapeString(item.Description)
		channel.Channel.Item[i] = item
	}

	return &channel, nil
}
