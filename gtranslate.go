// Package gtranslate provides free access to the Google Translate API.
//
// Usage example:
//
//   import "github.com/frei-0xff/gtranslate"
//   ...
//   ctx := context.Background()
// 	 results, err := gtranslate.Translate(ctx, []string{"Hello World!"}, language.English, language.French)

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

// Translate strings of text from a source language to a target language.
// All inputs must be in the same language.
//
// The returned strings appear in the same order as the inputs.
func Translate(ctx context.Context, inputs []string, source language.Tag, target language.Tag) ([]string, error) {
	tk, err := generateToken(ctx, strings.Join(inputs, ""))
	if err != nil {
		return nil, err
	}

	data := url.Values{}
	data.Set("sl", source.String())
	data.Set("tl", target.String())
	data.Set("tk", tk)
	for i := range inputs {
		data.Add("q", inputs[i])
	}

	req, err := http.NewRequest(http.MethodPost, "https://translate.googleapis.com/translate_a/t?client=te", strings.NewReader(data.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)

	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var results []string
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(body, &results)
	if len(results) != len(inputs) {
		return nil, fmt.Errorf("Bad response format: %s", body)
	}
	return results, nil
}
