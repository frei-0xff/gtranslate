# gtranslate #

[![GoDoc](https://img.shields.io/static/v1?label=godoc&message=reference&color=blue)](https://pkg.go.dev/github.com/frei-0xff/gtranslate)
[![Test Status](https://github.com/google/go-github/workflows/tests/badge.svg)](https://github.com/frei-0xff/gtranslate/actions?query=workflow%3Atest)

gtranslate is a Go client library for FREE and unlimited access to Google Translate API :dollar::no_entry_sign:  
Multiple strings could be translated with one HTTP request to the API endpoint.

## Installation ##

```bash
go get github.com/frei-0xff/gtranslate
```

## Usage ##

```go
import "github.com/frei-0xff/gtranslate"
ctx := context.Background()
results, err := gtranslate.Translate(ctx, []string{"Hello World!"}, language.English, language.French)
```


## Example ##

```go
package main

import (
	"context"
	"fmt"
	"time"

	"github.com/frei-0xff/gtranslate"
	"golang.org/x/text/language"
)

func main() {
	inputs := []string{"Hello World!", "What a wonderful world!"}
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Duration(time.Second*5))
	defer cancel()
	results, err := gtranslate.Translate(ctx, inputs, language.English, language.French)
	if err != nil {
		panic(err)
	}
	fmt.Println(results)
}
```
