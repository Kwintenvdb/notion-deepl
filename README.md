# notion-deepl
notion-deepl is a CLI tool to translate Notion pages with Deepl

# Install

```sh
go install github.com/Kwintenvdb/notion-deepl
```

# Overview

You can use this CLI tool to translate Notion pages (or blocks) using the DeepL API.

## Requirements

* Go 1.18+
* A [DeepL API key](https://www.deepl.com/pro-api/). A free API subscription is sufficient.
* A [Notion API key](https://developers.notion.com/). You will need to create your own Notion integration and enable it for the pages you wish to translate.

## Limitations

This tool was developed mostly for personal usage. As a result, several limitations exist at this moment:

* Formatting and rich text markup are mostly not preserved.
* Databases cannot be translated, though individual database pages can be.
* Tables cannot be translated.
* Blocks are currently translated in sequence. This may be further optimized in the future, but is currently limited by Notion's rate-limiting.

# Usage

> #### :warning: This tool will translate your Notion pages or blocks *in-place*. In other words, the text will be overwritten. If you want to test out the tool, or want to make changes safely, I recommend you make a copy of your page or block first before using this.

```
Usage:
  notion-deepl [flags]

Flags:
  -b, --block-id string          The Notion block ID to translate (required)
  -d, --deepl-api-key string     DeepL API key (required: if not provided, the DEEPL_API_KEY environment variable will be used)
  -f, --formality string         The formality of the translation (options: default, more, less, prefer_more, prefer_less) (default "default")
  -h, --help                     help for notion-deepl
  -n, --notion-api-key string    Notion API key (required: if not provided, the NOTION_API_KEY environment variable will be used)
  -s, --source-language string   The source language (optional: will be automatically detected for each block if not specified)
  -t, --target-language string   The target language (required)
```

See the [DeepL API docs](https://www.deepl.com/docs-api/translate-text/) for which source language and target language options are supported.

The page you want to translate (or blocks within this page) must be connected to your Notion integration before executing the translation. See the [Notion docs](https://developers.notion.com/docs#step-1-create-an-integration) for how to create and integrate a personal integration with a page.

## Examples

[Read here](https://stackoverflow.com/a/67652092/5089252) for how to get a block ID of a Notion page or block.

#### Translate a page using API keys from environment variables
```sh
notion-deepl --source-language=EN --target-language=NL --block-id=91ada48a2a3b4a939f5e1f0e510c4d3d
```

#### Translate a page by specifying API keys and formality via the CLI
```sh
notion-deepl --notion-api-key=[YOUR NOTION API KEY]
    --deepl-api-key=[YOUR DEEPL API KEY]
    --source-language=EN
    --target-language=NL
    --block-id=91ada48a2a3b4a939f5e1f0e510c4d3d
    --formality=prefer_less
```
