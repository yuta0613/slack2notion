# Slack2Notion

A command line tool to integrate Slack and Notion.

## Overview

This project fetches conversations and thread summaries from Slack and saves them to a Notion database.

## Prerequisites

- Set the necessary environment variables:
  - SLACK_BOT_TOKEN: Your Slack Bot token.
  - SLACK_CHANNEL_ID: The Slack channel ID to fetch the messages.
  - NOTION_API_TOKEN: Your Notion API token.
  - NOTION_DATABASE_ID: The target Notion database ID.

## Usage

### Running Slack Messages to Notion
To run the default mode which fetches Slack messages and saves them to Notion, execute:
```
go run main.go
```

### Running Thread Summaries
To fetch thread summaries from Slack and save them to Notion, run:
```
go run main.go thread-summary
```

Ensure you have set the necessary environment variables before running the commands.

## Notes

- The `thread-summary` mode aggregates thread summaries and saves the result to the Notion database.
- If you encounter any issues, verify that your environment variables are set correctly and that you have access to the required API endpoints.
