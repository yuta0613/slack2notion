package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

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

			// Build a structured summary for the thread using markdown
			summary := "## Thread Summary\n\n"
			summary += fmt.Sprintf("**Parent Message:** %s\n\n", message.Text)
			replyCount := len(replies) - 1
			summary += fmt.Sprintf("**Number of Replies:** %d\n\n", replyCount)
			if replyCount > 0 {
				summary += "**Replies:**\n"
				for i, reply := range replies {
					if i == 0 {
						continue
					}
					summary += fmt.Sprintf("- %s\n", reply.Text)
				}
			}
			summary += "\n---\n\n"

			// Print the summary to the console
			fmt.Print(summary)
			summaryAll += summary
		}
	}

	// Save the aggregated thread summaries to Notion if any summary exists
	if summaryAll != "" {
		notionClient := notionapi.NewClient(notionapi.Token(os.Getenv("NOTION_API_TOKEN")))
		rawBlocks := convertMarkdownToNotionBlocks(summaryAll)
		// Convert raw blocks (map based) into typed Notion blocks by converting each block individually.
		blocks := convertToNotionBlocks(rawBlocks)
		// Send the converted blocks to Notion
		addToNotionBlocks(notionClient, "Thread Summary", blocks)
	} else {
		log.Println("No thread summaries found")
	}
}

func parseRichText(s string) []interface{} {
	parts := strings.Split(s, "**")
	// if no bold formatting is present, return a single item
	if len(parts) == 1 {
		return []interface{}{
			map[string]interface{}{
				"type": "text",
				"text": map[string]interface{}{
					"content": s,
				},
			},
		}
	}
	var richTexts []interface{}
	for i, part := range parts {
		if part == "" {
			continue
		}
		item := map[string]interface{}{
			"type": "text",
			"text": map[string]interface{}{
				"content": part,
			},
		}
		if i%2 == 1 {
			item["annotations"] = map[string]interface{}{
				"bold": true,
			}
		}
		richTexts = append(richTexts, item)
	}
	return richTexts
}

func convertMarkdownToNotionBlocks(markdown string) []map[string]interface{} {
	lines := strings.Split(markdown, "\n")
	var blocks []map[string]interface{}

	for _, line := range lines {
		trimmedLine := strings.TrimSpace(line)
		if trimmedLine == "" {
			continue
		}
		if strings.HasPrefix(trimmedLine, "## ") {
			headingText := strings.TrimSpace(strings.TrimPrefix(trimmedLine, "## "))
			block := map[string]interface{}{
				"object": "block",
				"type":   "heading_2",
				"heading_2": map[string]interface{}{
					"rich_text": parseRichText(headingText),
				},
			}
			blocks = append(blocks, block)
		} else if strings.HasPrefix(trimmedLine, "- ") {
			content := strings.TrimPrefix(trimmedLine, "- ")
			block := map[string]interface{}{
				"object": "block",
				"type":   "bulleted_list_item",
				"bulleted_list_item": map[string]interface{}{
					"rich_text": parseRichText(content),
				},
			}
			blocks = append(blocks, block)
		} else {
			block := map[string]interface{}{
				"object": "block",
				"type":   "paragraph",
				"paragraph": map[string]interface{}{
					"rich_text": parseRichText(trimmedLine),
				},
			}
			blocks = append(blocks, block)
		}
	}
	return blocks
}

// unmarshalBlock converts JSON data into an appropriate notionapi.Block based on its "type" field.
func unmarshalBlock(data []byte) (notionapi.Block, error) {
	var raw map[string]interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, err
	}
	blockType, ok := raw["type"].(string)
	if !ok {
		return nil, fmt.Errorf("block type not found in JSON")
	}
	switch blockType {
	case "heading_2":
		var block notionapi.Heading2Block
		if err := json.Unmarshal(data, &block); err != nil {
			return nil, err
		}
		return block, nil
	case "paragraph":
		var block notionapi.ParagraphBlock
		if err := json.Unmarshal(data, &block); err != nil {
			return nil, err
		}
		return block, nil
	case "bulleted_list_item":
		var block notionapi.BulletedListItemBlock
		if err := json.Unmarshal(data, &block); err != nil {
			return nil, err
		}
		return block, nil
	default:
		return nil, fmt.Errorf("unsupported block type: %s", blockType)
	}
}

// convertToNotionBlocks converts raw block maps into a slice of notionapi.Block by converting each map individually.
func convertToNotionBlocks(rawBlocks []map[string]interface{}) []notionapi.Block {
	var blocks []notionapi.Block
	for _, rb := range rawBlocks {
		data, err := json.Marshal(rb)
		if err != nil {
			log.Printf("Error marshaling raw block: %v", err)
			continue
		}
		b, err := unmarshalBlock(data)
		if err != nil {
			log.Printf("Error unmarshaling block: %v", err)
			continue
		}
		blocks = append(blocks, b)
	}
	return blocks
}

func addToNotionBlocks(client *notionapi.Client, title string, blocks []notionapi.Block) {
	databaseID := os.Getenv("NOTION_DATABASE_ID")
	if databaseID == "" {
		log.Fatal("NOTION_DATABASE_ID environment variable is not set")
	}

	req := &notionapi.PageCreateRequest{
		Parent: notionapi.Parent{
			DatabaseID: notionapi.DatabaseID(databaseID),
		},
		Properties: notionapi.Properties{
			"Name": notionapi.TitleProperty{
				Title: []notionapi.RichText{
					{
						Type: "text",
						Text: &notionapi.Text{
							Content: title,
						},
					},
				},
			},
		},
		Children: blocks,
	}

	page, err := client.Page.Create(context.Background(), req)
	if err != nil {
		log.Printf("Error creating Notion page: %v", err)
	} else {
		fmt.Println("Notion page created successfully. Page ID:", page.ID)
	}
}
