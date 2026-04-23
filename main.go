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

func sendToDiscord(message, channelId, token string) error {
	if len(message) > 1900 {
		message = message[:1900] + "..."
	}

	payload := map[string]string{
		"content": message,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

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

func sendFeed(feed *gofeed.Feed, channelId, token string) error {
	feedTitle := htmlElementToMarkdown(feed.Title)

	for _, item := range feed.Items {
		message := ""

		message += fmt.Sprintf("#\n### %s\n", htmlElementToMarkdown(item.Title))
		message += fmt.Sprintf("%s\n", htmlElementToMarkdown(item.Description))
		message += fmt.Sprintf("%s\n", htmlElementToMarkdown(item.Link))
		message += fmt.Sprintf("%s\n", htmlElementToMarkdown(item.Published))
		message += feedTitle

		err := sendToDiscord(message, channelId, token)
		if err != nil {
			return err
		}
	}

	return nil
}

func main() {
	rssFeed := "https://cvedaily.com/feed-critical.xml"

	feed, err := fetchRSSFeed(rssFeed)
	if err != nil {
		log.Fatal(err)
	}

	env, err := os.ReadFile(".env")
	channelId := strings.Split(string(env), "\n")[0]
	token := strings.Split(string(env), "\n")[1]

	err = sendFeed(feed, channelId, token)
	if err != nil {
		log.Fatal(err)
	}
}
