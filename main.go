package main

import (
	"context"
	"os"

	// "github.com/Kwintenvdb/notion-deepl/deepl"
	"github.com/dstotijn/go-notion"
)

func main() {
	// deeplApiKey := os.Getenv("DEEPL_API_KEY")
	// deeplClient := deepl.NewClient(deeplApiKey)
	// translated, err := deeplClient.Translate("Hello", "EN", "DE")
	// if err != nil {
	// 	println("error")
	// 	println(err)
	// }
	// println(translated.Translations[0].Text)
	// println(translated.Translations[0].DetectedSourceLanguage)

	notionApiKey := os.Getenv("NOTION_API_KEY")
	notionClient := notion.NewClient(notionApiKey)


	ctx := context.Background()
	println("fetching page...")
	page, err := notionClient.FindPageByID(ctx, "d6e50c5cce5a463483af0a7c2274d263")
	if err != nil {
		println(err.Error())
	}
	println(page.ID)

	block := recursivelyRetrieveChildBlocks(notionClient, "d6e50c5cce5a463483af0a7c2274d263")
	println(block.block.ID())
}

type BlockWrapper struct {
	block notion.Block
	children []*BlockWrapper
}

func recursivelyRetrieveChildBlocks(notionClient *notion.Client, blockId string) *BlockWrapper {
	block, _ := notionClient.FindBlockByID(context.Background(), blockId)
	blockWrapper := &BlockWrapper{
		block: block,
		children: []*BlockWrapper{},
	}

	retrieveChildBlocks(notionClient, blockWrapper)

	return blockWrapper
}

func retrieveChildBlocks(notionClient *notion.Client, blockWrapper *BlockWrapper) {
	if blockWrapper.block.HasChildren() {
		query := notion.PaginationQuery{
			StartCursor: "",
			PageSize:    100,
		}
		blocks, _ := notionClient.FindBlockChildrenByID(context.Background(), blockWrapper.block.ID(), &query)
		for _, block := range blocks.Results {
			childBlockWrapper := &BlockWrapper{
				block: block,
				children: []*BlockWrapper{},
			}

			retrieveChildBlocks(notionClient, childBlockWrapper)

			blockWrapper.children = append(blockWrapper.children, childBlockWrapper)
		}
	}
}
