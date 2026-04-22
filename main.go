package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	htmltomarkdown "github.com/JohannesKaufmann/html-to-markdown/v2"
	"github.com/mmcdole/gofeed"
)

func fetchRSSFeed(url string) (*gofeed.Feed, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	fp := gofeed.NewParser()

	feed, err := fp.ParseURLWithContext(url, ctx)
	if err != nil {
		return nil, err
	}

	return feed, nil
}

func toMarkdown(str string) (string, error) {
	markdown, err := htmltomarkdown.ConvertString(str)
	if err != nil {
		return "", err
	}
	return markdown, err
}

func htmlElementToMarkdown(elem string) string {
	md, err := toMarkdown(elem)
	if err != nil {
		log.Println(err)
	}
	return md
}
func stringFeed(feed *gofeed.Feed) string {
	result := ""
	result += fmt.Sprintf("# %s\n", htmlElementToMarkdown(feed.Title))
	result += fmt.Sprintf("**Description:** %s\n", htmlElementToMarkdown(feed.Description))
	result += fmt.Sprintf("**Link:** %s\n", htmlElementToMarkdown(feed.Link))

	for _, item := range feed.Items {
		result += fmt.Sprintf("#\n### %s\n", htmlElementToMarkdown(item.Title))
		result += fmt.Sprintf("%s\n", htmlElementToMarkdown(item.Description))
		result += fmt.Sprintf("%s\n", htmlElementToMarkdown(item.Link))
		result += fmt.Sprintf("%s\n", htmlElementToMarkdown(item.Published))
	}

	return result
}

func sendToDiscord(feed string) error {
	if len(feed) > 1900 {
		feed = feed[:1900]
	}

	payload := map[string]string{
		"content": feed,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	env, err := os.ReadFile(".env")
	channelId := strings.Split(string(env), "\n")[0]
	token := strings.Split(string(env), "\n")[1]

	req, err := http.NewRequest(
		"POST",
		"https://discord.com/api/v10/channels/"+channelId+"/messages",
		bytes.NewBuffer(body),
	)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bot "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return err
	}

	return nil
}

func main() {
	rssFeed := "https://cvedaily.com/feed-critical.xml"

	feed, err := fetchRSSFeed(rssFeed)
	if err != nil {
		log.Fatal(err)
	}

	strFeed := stringFeed(feed)

	sendToDiscord(strFeed)
}
