package gtranslate

import (
	"context"
	"reflect"
	"testing"

	"golang.org/x/text/language"
)

func TestTranslate(t *testing.T) {
	type test struct {
		text   []string
		source language.Tag
		target language.Tag
		want   []string
	}

	tests := []test{
		{[]string{"Hello World!", "Test sentence"}, language.English, language.Russian, []string{"Привет, мир!", "Тестовое предложение"}},
		{[]string{"Привет, мир!", "Тестовая строка"}, language.Russian, language.Ukrainian, []string{"Привіт світ!", "Тестовий рядок"}},
	}
	ctx := context.Background()
	for _, tc := range tests {
		got, _ := Translate(ctx, tc.text, tc.source, tc.target)
		if !reflect.DeepEqual(tc.want, got) {
			t.Fatalf("expected: %v, got: %v", tc.want, got)
		}
	}
}
