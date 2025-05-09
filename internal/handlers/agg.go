package handlers

import (
	"context"
	"database/sql"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"net/http"
	"time"

	"github.com/lucashthiele/gator/internal/database"
	"github.com/lucashthiele/gator/internal/model"
)

func handleUnescapedSymbols(rss *model.RSSFeed) *model.RSSFeed {
	rss.Channel.Title = html.UnescapeString(rss.Channel.Title)
	rss.Channel.Description = html.UnescapeString(rss.Channel.Description)

	for i, item := range rss.Channel.Item {
		rss.Channel.Item[i].Title = html.UnescapeString(item.Title)
		rss.Channel.Item[i].Description = html.UnescapeString(item.Description)
	}

	return rss
}

func fetchFeed(ctx context.Context, feedUrl string) (*model.RSSFeed, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", feedUrl, nil)
	if err != nil {
		return &model.RSSFeed{}, err
	}

	req.Header.Set("User-Agent", "gator")

	client := &http.Client{}
	res, err := client.Do(req)

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return &model.RSSFeed{}, err
	}
	defer res.Body.Close()

	rssfeed := &model.RSSFeed{}
	err = xml.Unmarshal(body, rssfeed)
	if err != nil {
		return &model.RSSFeed{}, err
	}
	rssfeed = handleUnescapedSymbols(rssfeed)

	return rssfeed, err
}

func markFeedFetched(feed database.GetNextFeedToFetchRow, s *model.State) error {
	params := database.MarkFeedFetchedParams{
		LastFetchedAt: sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		},
		UpdatedAt: time.Now(),
		ID:        feed.ID,
	}
	return s.Db.MarkFeedFetched(context.Background(), params)
}

func printPosts(fetchedFeed *model.RSSFeed) {
	fmt.Println("List of post in this feed:")
	for _, post := range fetchedFeed.Channel.Item {
		fmt.Printf(" - %s\n", post.Title)
	}
}

func scrapeFeeds(s *model.State) {
	feeds, err := s.Db.GetNextFeedToFetch(context.Background())
	if err != nil {
		fmt.Printf("Error during scraping (getting next feed to fetch): %s\n", err.Error())
	}

	for _, feed := range feeds {
		fmt.Printf("Fetching feed: %s\n", feed.Name)

		fetchedFeed, err := fetchFeed(context.Background(), feed.Url)
		if err != nil {
			fmt.Printf("Error during scraping (fetching Feed): %s\n", err.Error())
		}

		printPosts(fetchedFeed)

		err = markFeedFetched(feed, s)
		if err != nil {
			fmt.Printf("Error during scraping (marking as fetched): %s\n", err.Error())
		}
	}
}

func HandlerAgg(s *model.State, cmd model.Command) error {
	expectedArguments := 1

	if len(cmd.Arguments) != expectedArguments {
		return fmt.Errorf(
			"expected %d arguments but got %d arguments\nyou need to provide the time between requests",
			expectedArguments,
			len(cmd.Arguments))
	}

	inputDuration := cmd.Arguments[0]
	timeBetweenRequests, err := time.ParseDuration(inputDuration)
	if err != nil {
		return err
	}

	fmt.Printf("Fetching feeds every %s\n", timeBetweenRequests.String())

	ticker := time.NewTicker(timeBetweenRequests)
	for ; ; <-ticker.C {
		scrapeFeeds(s)
	}
}
