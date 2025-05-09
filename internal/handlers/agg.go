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

	"github.com/google/uuid"
	"github.com/lib/pq"
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

func parseCustomTime(input string) (time.Time, error) {
	return time.Parse(time.RFC1123Z, input)
}

func savePosts(s *model.State, feedId uuid.UUID, fetchedFeed *model.RSSFeed) error {
	for _, post := range fetchedFeed.Channel.Item {
		pubDate, err := parseCustomTime(post.PubDate)
		if err != nil {
			return err
		}
		postParams := database.CreatePostParams{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Title:     post.Title,
			Description: sql.NullString{
				String: post.Description,
				Valid:  true,
			},
			PublishedAt: sql.NullTime{
				Time:  pubDate,
				Valid: true,
			},
			Url:    post.Link,
			FeedID: feedId,
		}

		_, err = s.Db.CreatePost(context.Background(), postParams)
		if err != nil {
			if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
				continue
			}
			return err
		}
	}

	return nil
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

		err = savePosts(s, feed.ID, fetchedFeed)
		if err != nil {
			fmt.Printf("Error during scraping (saving post): %s\n", err.Error())
		}

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
