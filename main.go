// main.go
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jomei/notionapi"
	"github.com/slack-go/slack"
)

func addToNotion(client *notionapi.Client, pageTitle string, content string) {
	// Example: Adding content to a Notion database
	page := notionapi.PageCreateRequest{
		Parent: notionapi.Parent{
			DatabaseID: notionapi.DatabaseID(os.Getenv("NOTION_DATABASE_ID")),
		},
		Properties: notionapi.Properties{
			"Name": notionapi.TitleProperty{
				Title: []notionapi.RichText{
					{
						Text: &notionapi.Text{
							Content: pageTitle,
						},
					},
				},
			},
		},
		Children: []notionapi.Block{
			&notionapi.ParagraphBlock{
				Paragraph: notionapi.Paragraph{
					RichText: []notionapi.RichText{
						{
							Text: &notionapi.Text{
								Content: content,
							},
						},
					},
				},
			},
		},
	}

	_, err := client.Page.Create(context.TODO(), &page)
	if err != nil {
		log.Fatalf("Error adding page to Notion: %v", err)
	}
}

func main() {
	api := slack.New(os.Getenv("SLACK_BOT_TOKEN"))

	// Example: Fetching Slack messages
	channelID := os.Getenv("SLACK_DATABASE_ID")
	params := &slack.GetConversationHistoryParameters{
		ChannelID: channelID,
		Limit:     10,
	}
	history, err := api.GetConversationHistory(params)
	if err != nil {
		log.Fatalf("Error fetching conversation history: %v", err)
	}

	for _, message := range history.Messages {
		fmt.Printf("Message: %s\n", message.Text)
	}

	// Notion integration logic
	notionClient := notionapi.NewClient(notionapi.Token(os.Getenv("NOTION_API_TOKEN")))

	// Example: Adding Slack messages to Notion
	for _, message := range history.Messages {
		addToNotion(notionClient, "Slack Message", message.Text)
	}
	fmt.Println("Slack to Notion CLI Tool")
}
