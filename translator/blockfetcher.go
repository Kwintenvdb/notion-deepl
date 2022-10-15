package translator

import (
	"context"

	"github.com/dstotijn/go-notion"
)

// Utility struct to wrap notion.Block and its children
// Used for recursive retrieval of blocks
type BlockWrapper struct {
	notion.Block
	Children []*BlockWrapper
}

func recursivelyRetrieveChildBlocks(notionClient *notion.Client, blockId string) *BlockWrapper {
	rootBlock, _ := notionClient.FindBlockByID(context.Background(), blockId)
	rootBlockWrapper := &BlockWrapper{
		Block:    rootBlock,
		Children: []*BlockWrapper{},
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
				Children: []*BlockWrapper{},
			}

			retrieveChildBlocks(notionClient, childBlockWrapper)

			parentBlockWrapper.Children = append(parentBlockWrapper.Children, childBlockWrapper)
		}
	}
}
