package cmd

import (
	"fmt"
	"os"

	"github.com/Kwintenvdb/notion-deepl/translator"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "notion-deepl",
	Short: "A CLI tool to translate Notion pages using DeepL",
	Run: func(cmd *cobra.Command, args []string) {
		blockId, _ := cmd.PersistentFlags().GetString("block-id")

		deeplApiKey := viper.GetString("DEEPL_API_KEY")
		notionApiKey := viper.GetString("NOTION_API_KEY")

		t := translator.NewTranslator(translator.TranslatorOptions{
			DeeplApiKey:  deeplApiKey,
			NotionApiKey: notionApiKey,
		})

		getFlag := func(flagName string) string {
			flag, _ := cmd.PersistentFlags().GetString(flagName)
			return flag
		}

		sourceLanguage := getFlag("source-language")
		targetLanguage := getFlag("target-language")
		formality := getFlag("formality")
		err := t.Translate(translator.TranslationArgs{
			SourceLanguage: sourceLanguage,
			TargetLanguage: targetLanguage,
			BlockId:        blockId,
			Formality:      formality,
		})
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(func() {
		viper.BindEnv("DEEPL_API_KEY")
		viper.BindEnv("NOTION_API_KEY")
	})

	rootCmd.PersistentFlags().StringP("deepl-api-key", "d", "", "DeepL API key (required: if not provided, the DEEPL_API_KEY environment variable will be used)")
	viper.BindPFlag("DEEPL_API_KEY", rootCmd.PersistentFlags().Lookup("deepl-api-key"))

	rootCmd.PersistentFlags().StringP("notion-api-key", "n", "", "Notion API key (required: if not provided, the NOTION_API_KEY environment variable will be used)")
	viper.BindPFlag("NOTION_API_KEY", rootCmd.PersistentFlags().Lookup("notion-api-key"))

	rootCmd.PersistentFlags().StringP("block-id", "b", "", "The Notion block ID to translate (required)")
	rootCmd.MarkFlagRequired("block-id")

	rootCmd.PersistentFlags().StringP("source-language", "s", "", "The source language (optional: will be automatically detected for each block if not specified)")
	rootCmd.PersistentFlags().StringP("target-language", "t", "", "The target language (required)")
	rootCmd.MarkFlagRequired("target-language")

	rootCmd.PersistentFlags().StringP("formality", "f", "default", "The formality of the translation (options: default, more, less, prefer_more, prefer_less)")
}
