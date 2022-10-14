package main

import (
	"context"
	"fmt"
	"os"

	"github.com/Kwintenvdb/notion-deepl/deepl"
	"github.com/dstotijn/go-notion"
)

func main() {
	deeplApiKey := os.Getenv("DEEPL_API_KEY")
	deeplClient := deepl.NewClient(deeplApiKey)

	notionApiKey := os.Getenv("NOTION_API_KEY")
	notionClient := notion.NewClient(notionApiKey)

	rootBlock := recursivelyRetrieveChildBlocks(notionClient, "2f1e789217f64fc7901f867e8b3b9e26")

	recursivelyTranslateBlocks(notionClient, deeplClient, rootBlock)
}

// Utility struct to wrap notion.Block and its children
// Used for recursive retrieval of blocks
type BlockWrapper struct {
	notion.Block
	children []*BlockWrapper
}

func recursivelyRetrieveChildBlocks(notionClient *notion.Client, blockId string) *BlockWrapper {
	rootBlock, _ := notionClient.FindBlockByID(context.Background(), blockId)
	rootBlockWrapper := &BlockWrapper{
		Block:    rootBlock,
		children: []*BlockWrapper{},
	}

	retrieveChildBlocks(notionClient, rootBlockWrapper)

	return rootBlockWrapper
}

func retrieveChildBlocks(notionClient *notion.Client, parentBlockWrapper *BlockWrapper) {
	if parentBlockWrapper.HasChildren() {
		query := notion.PaginationQuery{
			StartCursor: "",
			PageSize:    100,
		}
		// TODO handle pagination
		blocks, _ := notionClient.FindBlockChildrenByID(context.Background(), parentBlockWrapper.ID(), &query)
		for _, block := range blocks.Results {
			childBlockWrapper := &BlockWrapper{
				Block:    block,
				children: []*BlockWrapper{},
			}

			retrieveChildBlocks(notionClient, childBlockWrapper)

			parentBlockWrapper.children = append(parentBlockWrapper.children, childBlockWrapper)
		}
	}
}

func recursivelyTranslateBlocks(notionClient *notion.Client, deeplClient *deepl.Client, blockWrapper *BlockWrapper) {
	translateBlockRichText(blockWrapper, notionClient, deeplClient)

	// Recursively translate children
	for _, childBlockWrapper := range blockWrapper.children {
		recursivelyTranslateBlocks(notionClient, deeplClient, childBlockWrapper)
	}
}

// If the block has rich text, translate it
func translateBlockRichText(blockWrapper *BlockWrapper, notionClient *notion.Client, deeplClient *deepl.Client) {
	switch b := blockWrapper.Block.(type) {
	case *notion.ParagraphBlock:
		translateRichTextAndUpdateBlock(notionClient, deeplClient, b, b.RichText)
	case *notion.Heading1Block:
		translateRichTextAndUpdateBlock(notionClient, deeplClient, b, b.RichText)
	case *notion.Heading2Block:
		translateRichTextAndUpdateBlock(notionClient, deeplClient, b, b.RichText)
	case *notion.Heading3Block:
		translateRichTextAndUpdateBlock(notionClient, deeplClient, b, b.RichText)
	case *notion.BulletedListItemBlock:
		translateRichTextAndUpdateBlock(notionClient, deeplClient, b, b.RichText)
	case *notion.NumberedListItemBlock:
		translateRichTextAndUpdateBlock(notionClient, deeplClient, b, b.RichText)
	case *notion.ToDoBlock:
		translateRichTextAndUpdateBlock(notionClient, deeplClient, b, b.RichText)
	case *notion.ToggleBlock:
		translateRichTextAndUpdateBlock(notionClient, deeplClient, b, b.RichText)
	case *notion.QuoteBlock:
		translateRichTextAndUpdateBlock(notionClient, deeplClient, b, b.RichText)
	case *notion.CalloutBlock:
		translateRichTextAndUpdateBlock(notionClient, deeplClient, b, b.RichText)
	case *notion.CodeBlock:
		translateRichTextAndUpdateBlock(notionClient, deeplClient, b, b.RichText)
	}
}

// Translate rich text and update the block
func translateRichTextAndUpdateBlock(notionClient *notion.Client, deeplClient *deepl.Client, block notion.Block, richText []notion.RichText) {
	fmt.Printf("Translating block %s...\n", block.ID())
	translateRichText(richText, deeplClient)
	notionClient.UpdateBlock(context.Background(), block.ID(), block)
}

// Translate rich text in-place
func translateRichText(richText []notion.RichText, deeplClient *deepl.Client) {
	for _, rt := range richText {
		if rt.Text != nil {
			translated, err := translateText(deeplClient, rt.Text.Content)
			if err != nil {
				fmt.Printf("Error translating text: %s\n", err)
				continue
			}
			rt.Text.Content = translated
		}
	}
}

func translateText(deeplClient *deepl.Client, text string) (string, error) {
	translated, err := deeplClient.Translate(text, "EN", "DE")
	if err != nil {
		return "", err
	}
	if len(translated.Translations) == 0 {
		return "", fmt.Errorf("no translations found for text: %s", text)
	}
	return translated.Translations[0].Text, nil
}
