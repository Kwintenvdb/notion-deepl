package main

import (
	"os"

	"github.com/Kwintenvdb/notion-deepl/deepl"
)

func main() {
	apiKey := os.Getenv("DEEPL_API_KEY")
	client := deepl.NewClient(apiKey)
	translated, err := client.Translate("Hello", "EN", "DE")
	if err != nil {
		println("error")
		println(err)
	}
	println(translated.Translations[0].Text)
	println(translated.Translations[0].DetectedSourceLanguage)
}
