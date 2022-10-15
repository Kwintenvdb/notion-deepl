package deepl

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type Client struct {
	httpClient *http.Client
	apiKey     string
}

func NewClient(apiKey string) *Client {
	return &Client{
		httpClient: &http.Client{},
		apiKey:     apiKey,
	}
}

type TranslationResponse struct {
	Translations []Translation `json:"translations"`
}

type Translation struct {
	DetectedSourceLanguage string `json:"detected_source_language"`
	Text                   string `json:"text"`
}

func (c *Client) Translate(text, sourceLanguage, targetLanguage, formality string) (TranslationResponse, error) {
	data := url.Values{
		"text":        {text},
		"target_lang": {targetLanguage},
		"formality":   {formality},
	}

	encoded := data.Encode()
	req, err := http.NewRequest("POST", "https://api-free.deepl.com/v2/translate", strings.NewReader(encoded))
	if err != nil {
		return TranslationResponse{}, err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Authorization", "DeepL-Auth-Key " + c.apiKey)

	res, err := c.httpClient.Do(req)
	if err != nil {
		return TranslationResponse{}, err
	}

	if res.StatusCode != http.StatusOK {
		return TranslationResponse{}, errors.New("failed to translate text, status code: " + res.Status)
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return TranslationResponse{}, err
	}

	var translationResponse TranslationResponse
	err = json.Unmarshal(body, &translationResponse)
	if err != nil {
		return TranslationResponse{}, err
	}
	return translationResponse, nil
}
