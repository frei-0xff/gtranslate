// Package gtranslate provides free access to the Google Translate API.
//
// Usage example:
//
//   import "github.com/frei-0xff/gtranslate"
//   ...
//   ctx := context.Background()
// 	 result, err := gtranslate.Translate(ctx, "Hello World!", language.English, language.French)

package gtranslate

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"golang.org/x/text/language"
)

// Translate string from a source language to a target language.
//
// The source and target parameter supply languages to translate from and to respectively.
func Translate(ctx context.Context, input string, source language.Tag, target language.Tag) (string, error) {
	tk, err := generateToken(ctx, input)
	if err != nil {
		return "", err
	}

	data := url.Values{}
	data.Set("sl", source.String())
	data.Set("tl", target.String())
	data.Set("tk", tk)
	data.Set("q", input)

	req, err := http.NewRequest(http.MethodPost, "https://translate.googleapis.com/translate_a/t?client=te", strings.NewReader(data.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	if err != nil {
		return "", err
	}
	req = req.WithContext(ctx)

	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var value []string
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	err = json.Unmarshal(body, &value)
	if len(value) < 1 {
		return "", fmt.Errorf("Bad response format: %s", body)
	}
	return value[0], nil
}
