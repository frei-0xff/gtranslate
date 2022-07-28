package gtranslate

import (
	"context"
	"testing"

	"golang.org/x/text/language"
)

func TestTranslate(t *testing.T) {
	type test struct {
		text   string
		source language.Tag
		target language.Tag
		want   string
	}

	tests := []test{
		{"Hello World!", language.English, language.Russian, "Привет, мир!"},
		{"Hello World!", language.English, language.German, "Hallo Welt!"},
		{"Привет, мир!", language.Russian, language.Ukrainian, "Привіт світ!"},
	}
	ctx := context.Background()
	for _, tc := range tests {
		got, _ := Translate(ctx, tc.text, tc.source, tc.target)
		if tc.want != got {
			t.Fatalf("expected: %v, got: %v", tc.want, got)
		}
	}
}
