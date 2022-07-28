package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/bluele/gcache"
)

var gc gcache.Cache

func init() {
	gc = gcache.New(1).LRU().Build()
}

func getCTKK(ctx context.Context) (string, error) {
	ctkk, err := gc.Get("ctkk")
	if err == nil {
		return ctkk.(string), nil
	}
	req, err := http.NewRequest("GET", "https://translate.google.com/translate_a/element.js", nil)
	if err != nil {
		return "", err
	}
	req = req.WithContext(ctx)
	client := http.DefaultClient
	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	r := regexp.MustCompile(`c._ctkk='(.*?)'`)
	m := r.FindSubmatch(body)
	if len(m) == 2 {
		ctkk := string(m[1])
		gc.SetWithExpire("ctkk", ctkk, time.Hour*24)
		return ctkk, nil
	}
	return "", errors.New("Failed to get _ctkk")
}

func crypt(num, op string) string {
	iNum64, err := strconv.ParseInt(num, 10, 32)
	if err != nil {
		return ""
	}
	iNum := int32(iNum64)
	bOp := []byte(op)
	for i := 0; i < len(op)-2; i += 3 {
		c := int32(bOp[i+2])
		if 97 <= c {
			c -= 87
		} else {
			c -= 48
		}
		if 43 == bOp[i+1] {
			c = int32(uint32(iNum) >> c)
		} else {
			c = iNum << c
		}
		if 43 == bOp[i] {
			iNum = iNum + c
		} else {
			iNum = iNum ^ c
		}
	}
	return strconv.FormatInt(int64(iNum), 10)
}

func getTK(text, ctkk string) string {
	/*
			function Dn(text, ctkk) {
		  console.log(text, ctkk);
		  var parts = ctkk.split(".");
		  var t = Number(parts[0]) || 0;
		  var buf = [];
		  var j = 0;
		  var i = 0;
		  for (; i < text.length; i++) {
		    var ch = text.charCodeAt(i);
		    if (128 > ch) {
		      buf[j++] = ch;
		    } else {
		      if (2048 > ch) {
		        buf[j++] = ch >> 6 | 192;
		      } else {
		        if (55296 == (ch & 64512) && i + 1 < text.length && 56320 == (text.charCodeAt(i + 1) & 64512)) {
		          ch = 65536 + ((ch & 1023) << 10) + (text.charCodeAt(++i) & 1023);
		          buf[j++] = ch >> 18 | 240;
		          buf[j++] = ch >> 12 & 63 | 128;
		        } else {
		          buf[j++] = ch >> 12 | 224;
		        }
		        buf[j++] = ch >> 6 & 63 | 128;
		      }
		      buf[j++] = ch & 63 | 128;
		    }
		  }
		  text = t;
		  j = 0;
		  for (; j < buf.length; j++) {
		    text = text + buf[j];
		    text = Cn(text, "+-a^+6");
		  }
		  text = Cn(text, "+-3^+b+-f");
		  text = text ^ (Number(parts[1]) || 0);
		  if (0 > text) {
		    text = (text & 2147483647) + 2147483648;
		  }
		  parts = text % 1E6;
		  return parts.toString() + "." + (parts ^ t);
		}
	*/
	fmt.Println(text, ctkk)
	parts := strings.Split(ctkk, ".")
	if len(parts) != 2 {
		return ""
	}
	p1, _ := strconv.ParseInt(parts[0], 10, 32)
	rText := []rune(text)
	buf := make([]rune, len(rText)*3)
	for i, j := 0, 0; i < len(rText); i++ {
		ch := rText[i]
		if 128 > ch {
			buf[j] = ch
			j++
		} else {
			if 2048 > ch {
				buf[j] = ch>>6 | 192
				j++
			} else {
				if 55296 == (ch&64512) && i+1 < len(rText) && 56320 == (rText[i+1]&64512) {
					i++
					ch = 65536 + ((ch & 1023) << 10) + (rText[i] & 1023)
					buf[j] = ch>>18 | 240
					j++
					buf[j] = ch>>12&63 | 128
					j++
				} else {
					buf[j] = ch>>12 | 224
					j++
				}
				buf[j] = ch>>6&63 | 128
				j++
			}
			buf[j] = ch&63 | 128
			j++
		}
	}
	text = parts[0]
	for i := range buf {
		if buf[i] == 0 {
			buf = buf[:i]
			break
		}
	}
	for j := 0; j < len(buf); j++ {
		intText, _ := strconv.ParseInt(text, 10, 32)
		text = strconv.FormatInt(intText+int64(buf[j]), 10)
		text = crypt(text, "+-a^+6")
	}
	text = crypt(text, "+-3^+b+-f")
	p2, _ := strconv.ParseInt(parts[1], 10, 64)
	intText, _ := strconv.ParseInt(text, 10, 32)
	intText ^= p2
	if 0 > intText {
		intText = int64((uint32(intText) & 2147483647) + 2147483648)
	}
	mod := intText % 1e6
	return strconv.FormatInt(mod, 10) + "." + strconv.FormatInt(mod^p1, 10)
}

func translate(ctx context.Context, text string) (string, error) {
	ctkk, err := getCTKK(ctx)
	if err != nil {
		return "", err
	}
	tk := getTK(text, ctkk)
	resp, err := http.PostForm("https://translate.googleapis.com/translate_a/t?anno=3&client=te_lib&format=html&v=1.0&key&logld=vTE_20220719&sl=ru&tl=uk&tc=1&dom=1&sr=1&tk="+tk+"&mode=1",
		url.Values{"q": {text}})
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	return string(body), nil
}
