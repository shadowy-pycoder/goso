package main

import (
	"crypto/tls"
	"html"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/alecthomas/chroma/v2/formatters"
	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/alecthomas/chroma/v2/styles"
	"github.com/shadowy-pycoder/goso"
)

const (
	codeStartTag string = "<pre><code>"
	codeEndTag   string = "</code></pre>"
)

var (
	codePattern  = regexp.MustCompile(`<pre\s.*?>`)
	aHrefPattern = regexp.MustCompile(`<a\s+(?:[^>]*?\s+)?href=(["'])?([^\'" >]+)(.*?)?</a>`)
)

func fmtText(text string) string {
	t := html.UnescapeString(text)
	t = strings.ReplaceAll(t, "<p>", "")
	t = strings.ReplaceAll(t, "</p>", "")
	t = strings.ReplaceAll(t, "<strong>", "\033[1m")
	t = strings.ReplaceAll(t, "</strong>", "\033[0m")
	t = strings.ReplaceAll(t, "<em>", "\033[3m")
	t = strings.ReplaceAll(t, "</em>", "\033[0m")
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
	t = strings.ReplaceAll(t, "<del>", "\033[9m")
	t = strings.ReplaceAll(t, "</del>", "\033[0m")
	t = strings.ReplaceAll(t, "<ins>", "")
	t = strings.ReplaceAll(t, "</ins>", "")
	t = codePattern.ReplaceAllString(t, "<pre>")
	t = aHrefPattern.ReplaceAllString(t, "\n- $2")
	return t
}

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
	var sb strings.Builder
	style := styles.Get("onedark")
	if style == nil {
		style = styles.Fallback
	}
	formatter := formatters.Get("terminal16m")
	if formatter == nil {
		formatter = formatters.Fallback
	}
	lexer := lexers.Get("go")
	if lexer == nil {
		lexer = lexers.Fallback
	}
	for _, item := range resp.Items {
		t := fmtText(item.Body)
		codeStartIdx := strings.Index(t, codeStartTag)
		var codeEndIdx int
		for codeStartIdx != -1 {
			codeEndIdx = strings.Index(t, codeEndTag)
			if codeEndIdx == -1 {
				println("malformed html")
				break
			}
			iterator, err := lexer.Tokenise(nil, t[codeStartIdx+len(codeStartTag):codeEndIdx])
			if err != nil {
				println(err)
			}
			err = formatter.Format(&sb, style, iterator)
			if err != nil {
				println(err)
			}
			t = t[:codeStartIdx] + sb.String() + t[codeEndIdx+len(codeEndTag):]
			codeStartIdx = strings.Index(t, codeStartTag)
			sb.Reset()
		}

		t = strings.ReplaceAll(t, "<code>", "\033[32m")
		t = strings.ReplaceAll(t, "</code>", "\033[0m")

		println(t)
	}
}
