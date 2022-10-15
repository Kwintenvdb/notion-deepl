package translator

import (
	"context"
	"errors"
	"fmt"

	"github.com/Kwintenvdb/notion-deepl/deepl"
	"github.com/dstotijn/go-notion"
)

type TranslationArgs struct {
	SourceLanguage string
	TargetLanguage string
	BlockId        string
}

type Translator interface {
	Translate(args TranslationArgs) error
}

type notionDeeplTranslator struct {
	deeplClient  *deepl.Client
	notionClient *notion.Client
}

func (t *notionDeeplTranslator) Translate(args TranslationArgs) error {
	if args.TargetLanguage == "" {
		return errors.New("target language is required")
	}

	if args.BlockId == "" {
		return errors.New("block id is required")
	}

	if args.SourceLanguage == args.TargetLanguage {
		return errors.New("source and target language are the same")
	}

	// A little slow, but parallelizing this makes the Notion API return 409 responses
	rootBlock := recursivelyRetrieveChildBlocks(t.notionClient, args.BlockId)
	t.recursivelyTranslateBlocks(rootBlock, args)

	return nil
}

type TranslatorOptions struct {
	DeeplApiKey  string
	NotionApiKey string
}

func NewTranslator(options TranslatorOptions) Translator {
	deeplClient := deepl.NewClient(options.DeeplApiKey)
	notionClient := notion.NewClient(options.NotionApiKey)
	return &notionDeeplTranslator{
		deeplClient:  deeplClient,
		notionClient: notionClient,
	}
}

func (t *notionDeeplTranslator) recursivelyTranslateBlocks(blockWrapper *BlockWrapper, args TranslationArgs) {
	t.translateBlockRichText(blockWrapper, args)

	// Recursively translate children
	for _, childBlockWrapper := range blockWrapper.Children {
		t.recursivelyTranslateBlocks(childBlockWrapper, args)
	}
}

// If the block has rich text, translate it
func (t *notionDeeplTranslator) translateBlockRichText(blockWrapper *BlockWrapper, args TranslationArgs) {
	switch block := blockWrapper.Block.(type) {
	case *notion.ParagraphBlock:
		t.translateRichTextAndUpdateBlock(block, block.RichText, args)
	case *notion.Heading1Block:
		t.translateRichTextAndUpdateBlock(block, block.RichText, args)
	case *notion.Heading2Block:
		t.translateRichTextAndUpdateBlock(block, block.RichText, args)
	case *notion.Heading3Block:
		t.translateRichTextAndUpdateBlock(block, block.RichText, args)
	case *notion.BulletedListItemBlock:
		t.translateRichTextAndUpdateBlock(block, block.RichText, args)
	case *notion.NumberedListItemBlock:
		t.translateRichTextAndUpdateBlock(block, block.RichText, args)
	case *notion.ToDoBlock:
		t.translateRichTextAndUpdateBlock(block, block.RichText, args)
	case *notion.ToggleBlock:
		t.translateRichTextAndUpdateBlock(block, block.RichText, args)
	case *notion.QuoteBlock:
		t.translateRichTextAndUpdateBlock(block, block.RichText, args)
	case *notion.CalloutBlock:
		t.translateRichTextAndUpdateBlock(block, block.RichText, args)
	case *notion.CodeBlock:
		t.translateRichTextAndUpdateBlock(block, block.RichText, args)
	}
}

// Translate rich text and update the block
func (t *notionDeeplTranslator) translateRichTextAndUpdateBlock(block notion.Block, richText []notion.RichText, args TranslationArgs) {
	fmt.Printf("Translating block %s...\n", block.ID())
	t.translateRichText(richText, args)
	fmt.Printf("Updating block %s...\n", block.ID())
	_, err := t.notionClient.UpdateBlock(context.Background(), block.ID(), block)
	if err != nil {
		fmt.Printf("Error updating block %s: %s\n", block.ID(), err)
	}
}

// Translate rich text in-place
func (t *notionDeeplTranslator) translateRichText(richText []notion.RichText, args TranslationArgs) {
	for _, rt := range richText {
		if rt.Text != nil {
			fmt.Printf("Translating text %s...\n", rt.Text.Content)
			translated, err := t.translateText(rt.Text.Content, args)
			if err != nil {
				fmt.Printf("Error translating text: %s\n", err)
				continue
			}
			rt.Text.Content = translated
		}
	}
}

func (t *notionDeeplTranslator) translateText(text string, args TranslationArgs) (string, error) {
	translated, err := t.deeplClient.Translate(text, args.SourceLanguage, args.TargetLanguage)
	if err != nil {
		return "", err
	}
	if len(translated.Translations) == 0 {
		return "", fmt.Errorf("no translations found for text: %s", text)
	}
	return translated.Translations[0].Text, nil
}