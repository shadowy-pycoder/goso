package main

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/shadowy-pycoder/goso"
)

func main() {
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
		Timeout: time.Duration(10) * time.Second,
	}
	resp, err := goso.GetText(client)
	if err != nil {
		println(err)
	}
	for _, item := range resp.Items {
		t := item.Body
		t = strings.ReplaceAll(t, "<p>", "")
		t = strings.ReplaceAll(t, "</p>", "")
		t = strings.ReplaceAll(t, "<strong>", "\033[1m")
		t = strings.ReplaceAll(t, "</strong>", "\033[0m")
		t = strings.ReplaceAll(t, "<em>", "\033[3m")
		t = strings.ReplaceAll(t, "</em>", "\033[0m")
		t = strings.ReplaceAll(t, "&lt;", "<")
		t = strings.ReplaceAll(t, "&gt;", ">")
		t = strings.ReplaceAll(t, "&quot;", "\"")
		t = strings.ReplaceAll(t, "<pre><code>", "\033[33m")
		t = strings.ReplaceAll(t, "</code></pre>", "\033[0m")
		t = strings.ReplaceAll(t, "<code>", "\033[32m")
		t = strings.ReplaceAll(t, "</code>", "\033[0m")
		t = strings.ReplaceAll(t, "&amp;", "&")
		t = strings.ReplaceAll(t, "<ul>", "")
		t = strings.ReplaceAll(t, "</ul>", "")
		t = strings.ReplaceAll(t, "<ol>", "")
		t = strings.ReplaceAll(t, "</ol>", "")
		t = strings.ReplaceAll(t, "<li>", " - ")
		t = strings.ReplaceAll(t, "</li>", "")
		t = strings.ReplaceAll(t, "<hr>", "_______________________________________________________________________________________________________")
		t = strings.ReplaceAll(t, "<b>", "\033[1m")
		t = strings.ReplaceAll(t, "</b>", "\033[0m")
		t = strings.ReplaceAll(t, "<br>", "\n")
		t = strings.ReplaceAll(t, "<blockquote>", "\033[3m")
		t = strings.ReplaceAll(t, "</blockquote>", "\033[0m")

		fmt.Println(t)
	}
}
