package main

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestConvertMarkdownToNotionBlocks(t *testing.T) {
	// Given markdown with a heading, bold text, and a regular paragraph.
	markdown := "## Test Heading\n\n**Bold Text**\n\nRegular text"
	blocks := convertMarkdownToNotionBlocks(markdown)

	// Expecting 3 blocks: one heading, one paragraph with bold text, one regular paragraph.
	if len(blocks) != 3 {
		t.Fatalf("Expected 3 blocks, got %d", len(blocks))
	}

	// Test the first block is a heading block with correct content.
	headingBlock := blocks[0]
	if headingBlock["type"] != "heading_2" {
		t.Errorf("Expected first block type to be 'heading_2', got %v", headingBlock["type"])
	}
	headingData, ok := headingBlock["heading_2"].(map[string]interface{})
	if !ok {
		t.Fatal("Expected 'heading_2' field to be a map")
	}
	richText, ok := headingData["rich_text"].([]interface{})
	if !ok || len(richText) == 0 {
		t.Fatal("Expected non-empty 'rich_text' for heading block")
	}
	firstText, ok := richText[0].(map[string]interface{})
	if !ok {
		t.Fatal("Expected first rich_text element to be a map")
	}
	textField, ok := firstText["text"].(map[string]interface{})
	if !ok {
		t.Fatal("Expected 'text' field in rich_text element to be a map")
	}
	if textField["content"] != "Test Heading" {
		t.Errorf("Expected heading content 'Test Heading', got %v", textField["content"])
	}

	// Test the second block is a paragraph with bold text.
	boldBlock := blocks[1]
	if boldBlock["type"] != "paragraph" {
		t.Errorf("Expected second block type to be 'paragraph', got %v", boldBlock["type"])
	}
	paraData, ok := boldBlock["paragraph"].(map[string]interface{})
	if !ok {
		t.Fatal("Expected 'paragraph' field to be a map")
	}
	richText, ok = paraData["rich_text"].([]interface{})
	if !ok || len(richText) == 0 {
		t.Fatal("Expected non-empty 'rich_text' for paragraph block")
	}
	boldText, ok := richText[0].(map[string]interface{})
	if !ok {
		t.Fatal("Expected rich_text element to be a map")
	}
	textField, ok = boldText["text"].(map[string]interface{})
	if !ok {
		t.Fatal("Expected 'text' field in rich_text element to be a map")
	}
	if textField["content"] != "Bold Text" {
		t.Errorf("Expected bold content 'Bold Text', got %v", textField["content"])
	}
	annotations, ok := boldText["annotations"].(map[string]interface{})
	if !ok {
		t.Fatal("Expected 'annotations' field to be a map")
	}
	if bold, ok := annotations["bold"].(bool); !ok || !bold {
		t.Errorf("Expected annotations.bold to be true, got %v", annotations["bold"])
	}

	// Test the third block is a regular paragraph.
	paraBlock := blocks[2]
	if paraBlock["type"] != "paragraph" {
		t.Errorf("Expected third block type to be 'paragraph', got %v", paraBlock["type"])
	}
	paraData, ok = paraBlock["paragraph"].(map[string]interface{})
	if !ok {
		t.Fatal("Expected 'paragraph' field to be a map")
	}
	richText, ok = paraData["rich_text"].([]interface{})
	if !ok || len(richText) == 0 {
		t.Fatal("Expected non-empty 'rich_text' for regular paragraph block")
	}
	regularText, ok := richText[0].(map[string]interface{})
	if !ok {
		t.Fatal("Expected rich_text element to be a map")
	}
	textField, ok = regularText["text"].(map[string]interface{})
	if !ok {
		t.Fatal("Expected 'text' field in rich_text element to be a map")
	}
	if textField["content"] != "Regular text" {
		t.Errorf("Expected paragraph content 'Regular text', got %v", textField["content"])
	}
}

func TestConvertToNotionBlocks(t *testing.T) {
	// Create a raw block map
	rawBlocks := []map[string]interface{}{
		{
			"object": "block",
			"type":   "paragraph",
			"paragraph": map[string]interface{}{
				"rich_text": []map[string]interface{}{
					{
						"type": "text",
						"text": map[string]interface{}{
							"content": "Sample text",
						},
					},
				},
			},
		},
	}

	blocks := convertToNotionBlocks(rawBlocks)
	if len(blocks) != 1 {
		t.Fatalf("Expected 1 block, got %d", len(blocks))
	}

	// Marshal back to JSON and compare with original rawBlocks for structural equality.
	dataRaw, err := json.Marshal(rawBlocks)
	if err != nil {
		t.Fatalf("Error marshaling rawBlocks: %v", err)
	}
	dataBlocks, err := json.Marshal(blocks)
	if err != nil {
		t.Fatalf("Error marshaling blocks: %v", err)
	}
	if !reflect.DeepEqual(dataRaw, dataBlocks) {
		t.Errorf("Expected converted blocks to match rawBlocks, got %s vs %s", dataRaw, dataBlocks)
	}
}
