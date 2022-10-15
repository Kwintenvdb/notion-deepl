package main

import (
	"os"

	"github.com/Kwintenvdb/notion-deepl/translator"
)

func main() {
	deeplApiKey := os.Getenv("DEEPL_API_KEY")
	notionApiKey := os.Getenv("NOTION_API_KEY")

	t := translator.NewTranslator(translator.TranslatorOptions{
		DeeplApiKey:  deeplApiKey,
		NotionApiKey: notionApiKey,
	})

	t.Translate(translator.TranslationArgs{
		SourceLanguage: "DE",
		TargetLanguage: "NL",
		BlockId:        "91ada48a2a3b4a939f5e1f0e510c4d3d",
	})
}
