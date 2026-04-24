package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
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

func sendFeed(feed *gofeed.Feed, channelId, token, filePath string) error {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}
	lastGuid := string(content)
	newGuid := ""

	feedTitle := htmlElementToMarkdown(feed.Title)

	for i, item := range feed.Items {
		if i == 0 {
			newGuid = item.GUID
		}
		if item.GUID == lastGuid {
			break
		}

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

	if newGuid != lastGuid {
		return os.WriteFile(filePath, []byte(newGuid), 0644)
	}
	return nil
}

func initOutputFile(filePath string) error {
	_, err := os.Stat(filePath)
	if errors.Is(err, os.ErrNotExist) {
		_, err := os.Create(filePath)
		if err != nil {
			return err
		}

	} else if err != nil {
		return err
	}

	return nil
}

func main() {
	dir := "logs"

	err := os.MkdirAll(dir, 0755)
	if err != nil {
		log.Fatal(err)
	}

	feedUrl := "https://cvedaily.com/feed-critical.xml"
	feedName := "cve-daily"

	filePath := dir + "/" + feedName
	err = initOutputFile(filePath)
	if err != nil {
		log.Fatal(err)
	}

	feed, err := fetchRSSFeed(feedUrl)
	if err != nil {
		log.Fatal(err)
	}

	env, err := os.ReadFile(".env")
	channelId := strings.Split(string(env), "\n")[0]
	token := strings.Split(string(env), "\n")[1]

	err = sendFeed(feed, channelId, token, filePath)
	if err != nil {
		log.Fatal(err)
	}
}
