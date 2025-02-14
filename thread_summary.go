package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jomei/notionapi"
	"github.com/slack-go/slack"
)

func runThreadSummary() {
	api := slack.New(os.Getenv("SLACK_BOT_TOKEN"))
	channelID := os.Getenv("SLACK_CHANNEL_ID")

	// Fetch conversation history with increased limit for threads
	params := &slack.GetConversationHistoryParameters{
		ChannelID: channelID,
		Limit:     50,
	}
	history, err := api.GetConversationHistory(params)
	if err != nil {
		log.Fatalf("Error fetching conversation history: %v", err)
	}

	var summaryAll string

	// Iterate through messages to find threads (parent messages with a thread timestamp)
	for _, message := range history.Messages {
		if message.ThreadTimestamp != "" {
			// Fetch all messages in the thread using GetConversationRepliesContext
			replies, _, _, err := api.GetConversationRepliesContext(context.Background(), &slack.GetConversationRepliesParameters{
				ChannelID: channelID,
				Timestamp: message.ThreadTimestamp,
			})
			if err != nil {
				log.Printf("Error fetching thread replies: %v", err)
				continue
			}

			// Build a structured summary for the thread
			summary := "Thread Summary:\n"
			summary += fmt.Sprintf("Parent Message: %s\n", message.Text)
			replyCount := len(replies) - 1
			summary += fmt.Sprintf("Number of Replies: %d\n", replyCount)
			if replyCount > 0 {
				summary += "Replies:\n"
				for i, reply := range replies {
					if i == 0 {
						continue
					}
					summary += fmt.Sprintf("  - %s\n", reply.Text)
				}
			}
			summary += "-----\n"

			// Print the summary to the console
			fmt.Print(summary)
			summaryAll += summary
		}
	}

	// Save the aggregated thread summaries to Notion if any summary exists
	if summaryAll != "" {
		notionClient := notionapi.NewClient(notionapi.Token(os.Getenv("NOTION_API_TOKEN")))
		addToNotion(notionClient, "Thread Summary", summaryAll)
	} else {
		log.Println("No thread summaries found")
	}
}
