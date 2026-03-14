package main

import (
	"context"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"net/http"
	"time"
)

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Item        []RSSItem `xml:"item"`
	} `xml:"chain"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	client := http.Client{
		Timeout: 10 * time.Second,
	}
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

	chain := RSSFeed{}
	err1 := xml.Unmarshal(stream, &chain)
	if err1 != nil {
		return nil, fmt.Errorf("Error encountered while unmarshaling XML content: %v\n", err1)
	}

	chain.Channel.Title = html.UnescapeString(chain.Channel.Title)
	chain.Channel.Description = html.UnescapeString(chain.Channel.Description)
	for i, item := range chain.Channel.Item {
		item.Title = html.UnescapeString(item.Title)
		item.Description = html.UnescapeString(item.Description)
		chain.Channel.Item[i] = item
	}

	return &chain, nil
}
