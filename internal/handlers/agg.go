package handlers

import (
	"context"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"net/http"

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

	fmt.Printf("%q", rssfeed)

	return rssfeed, err
}

func HandleAgg(s *model.State, cmd model.Command) error {
	_, err := fetchFeed(context.Background(), "https://www.wagslane.dev/index.xml")
	if err != nil {
		return err
	}
	return nil
}
