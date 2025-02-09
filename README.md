# slack2notion

## Overview

This CLI tool fetches messages from a specified Slack channel and adds them to a Notion database.

## Setup

1. **Install Go**: Ensure you have Go installed on your system. You can download it from [golang.org](https://golang.org/dl/).

2. **Clone the Repository**: Clone this repository to your local machine.

3. **Set Environment Variables**:
   - Use the `.envrc` file to set your environment variables. Add the following lines to your `.envrc` file:
     ```bash
     export SLACK_BOT_TOKEN='your-slack-bot-token'
     export NOTION_API_TOKEN='your-notion-api-token'
     export SLACK_DATABASE_ID='your-slack-database-id'
     export NOTION_DATABASE_ID='your-notion-database-id'
     ```
   - After editing the `.envrc` file, run `direnv allow` to apply the changes.

   - `SLACK_BOT_TOKEN`: Your Slack bot token.
   - `NOTION_API_TOKEN`: Your Notion API token.

4. **Update Configuration**:
   - Replace `"your-channel-id"` in `main.go` with your Slack channel ID.
   - Replace `"your-database-id"` in `main.go` with your Notion database ID.

## Running the Tool

Navigate to the project directory and run the following command:

```bash
go run main.go
```

This command will execute the tool, fetching messages from the specified Slack channel and adding them to the Notion database.

## Running the Tool with Environment Variables

You can set the Slack and Notion database IDs directly in the command line when running the tool:

```bash
SLACK_DATABASE_ID='your-slack-database-id' NOTION_DATABASE_ID='your-notion-database-id' go run main.go
```

Ensure that the environment variables are set correctly before running the tool.

## Example

```bash
export SLACK_BOT_TOKEN='xoxb-your-slack-bot-token'
export NOTION_API_TOKEN='your-notion-api-token'
go run main.go
```

Ensure that the environment variables are set correctly before running the tool.
