package goso

import (
	"encoding/json"
	"fmt"
	"html"
	"net/http"
	netUrl "net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/alecthomas/chroma/v2/formatters"
	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/alecthomas/chroma/v2/styles"
)

const (
	codeStartTag string = "<pre><code>"
	codeEndTag   string = "</code></pre>"
)

var (
	codeStartIdx int
	codeEndIdx   int
	codePattern  = regexp.MustCompile(`<pre\s.*?>`)
	aHrefPattern = regexp.MustCompile(`<a\s+(?:[^>]*?\s+)?href=(["'])?([^\'" >]+)(.*?)?</a>`)
	r            = strings.NewReplacer("<p>", "",
		"</p>", "",
		"<strong>", "\033[1m",
		"</strong>", "\033[0m",
		"<em>", "\033[3m",
		"</em>", "\033[0m",
		"<ul>", "",
		"</ul>", "",
		"<ol>", "",
		"</ol>", "",
		"<li>", " - ",
		"</li>", "",
		"<hr>", "_______________________________________________________________________________________________________",
		"<b>", "\033[1m",
		"</b>", "\033[0m",
		"<br>", "\n",
		"<blockquote>", "\033[3m",
		"</blockquote>", "\033[0m",
		"<del>", "\033[9m",
		"</del>", "\033[0m",
		"<ins>", "",
		"</ins>", "",
	)
)

type googleSearchResult struct {
	Kind string `json:"kind"`
	URL  struct {
		Type     string `json:"type"`
		Template string `json:"template"`
	} `json:"url"`
	Queries struct {
		Request []struct {
			Title          string `json:"title"`
			TotalResults   string `json:"totalResults"`
			SearchTerms    string `json:"searchTerms"`
			Count          int    `json:"count"`
			StartIndex     int    `json:"startIndex"`
			InputEncoding  string `json:"inputEncoding"`
			OutputEncoding string `json:"outputEncoding"`
			Safe           string `json:"safe"`
			Cx             string `json:"cx"`
		} `json:"request"`
		NextPage []struct {
			Title          string `json:"title"`
			TotalResults   string `json:"totalResults"`
			SearchTerms    string `json:"searchTerms"`
			Count          int    `json:"count"`
			StartIndex     int    `json:"startIndex"`
			InputEncoding  string `json:"inputEncoding"`
			OutputEncoding string `json:"outputEncoding"`
			Safe           string `json:"safe"`
			Cx             string `json:"cx"`
		} `json:"nextPage"`
	} `json:"queries"`
	Context struct {
		Title string `json:"title"`
	} `json:"context"`
	SearchInformation struct {
		SearchTime            float64 `json:"searchTime"`
		FormattedSearchTime   string  `json:"formattedSearchTime"`
		TotalResults          string  `json:"totalResults"`
		FormattedTotalResults string  `json:"formattedTotalResults"`
	} `json:"searchInformation"`
	Items []struct {
		Kind             string `json:"kind"`
		Title            string `json:"title"`
		HTMLTitle        string `json:"htmlTitle"`
		Link             string `json:"link"`
		DisplayLink      string `json:"displayLink"`
		Snippet          string `json:"snippet"`
		HTMLSnippet      string `json:"htmlSnippet"`
		FormattedURL     string `json:"formattedUrl"`
		HTMLFormattedURL string `json:"htmlFormattedUrl"`
		Pagemap          struct {
			CseThumbnail []struct {
				Src    string `json:"src"`
				Width  string `json:"width"`
				Height string `json:"height"`
			} `json:"cse_thumbnail"`
			Qapage []struct {
				Image              string `json:"image"`
				Primaryimageofpage string `json:"primaryimageofpage"`
				Name               string `json:"name"`
				Description        string `json:"description"`
			} `json:"qapage"`
			Question []struct {
				Image       string `json:"image"`
				Upvotecount string `json:"upvotecount"`
				Answercount string `json:"answercount"`
				Name        string `json:"name"`
				Datecreated string `json:"datecreated"`
				Text        string `json:"text"`
				URL         string `json:"url"`
			} `json:"question"`
			Answer []struct {
				Upvotecount  string `json:"upvotecount"`
				Commentcount string `json:"commentcount,omitempty"`
				Text         string `json:"text"`
				Datecreated  string `json:"datecreated"`
				URL          string `json:"url"`
			} `json:"answer"`
			Person []struct {
				Name string `json:"name"`
			} `json:"person"`
			Metatags []struct {
				OgImage            string `json:"og:image"`
				OgType             string `json:"og:type"`
				TwitterCard        string `json:"twitter:card"`
				TwitterTitle       string `json:"twitter:title"`
				OgSiteName         string `json:"og:site_name"`
				TwitterDomain      string `json:"twitter:domain"`
				Viewport           string `json:"viewport"`
				TwitterDescription string `json:"twitter:description"`
				Bingbot            string `json:"bingbot"`
				OgURL              string `json:"og:url"`
			} `json:"metatags"`
			CseImage []struct {
				Src string `json:"src"`
			} `json:"cse_image"`
		} `json:"pagemap"`
	} `json:"items"`
}

type stackOverlowResult struct {
	Items []struct {
		Owner struct {
			AccountID    int    `json:"account_id"`
			Reputation   int    `json:"reputation"`
			UserID       int    `json:"user_id"`
			UserType     string `json:"user_type"`
			AcceptRate   int    `json:"accept_rate"`
			ProfileImage string `json:"profile_image"`
			DisplayName  string `json:"display_name"`
			Link         string `json:"link"`
		} `json:"owner"`
		IsAccepted         bool   `json:"is_accepted"`
		Score              int    `json:"score"`
		LastActivityDate   int    `json:"last_activity_date"`
		LastEditDate       int    `json:"last_edit_date,omitempty"`
		CreationDate       int    `json:"creation_date"`
		AnswerID           int    `json:"answer_id"`
		QuestionID         int    `json:"question_id"`
		ContentLicense     string `json:"content_license"`
		Body               string `json:"body"`
		CommunityOwnedDate int    `json:"community_owned_date,omitempty"`
	} `json:"items"`
	HasMore        bool `json:"has_more"`
	QuotaMax       int  `json:"quota_max"`
	QuotaRemaining int  `json:"quota_remaining"`
}

type Config struct {
	ApiKey       string
	SearchEngine string
	Query        string
	Style        string
	Lexer        string
	QuestionNum  int
	AnswerNum    int
	Client       *http.Client
}

func fmtText(text string) string {
	t := html.UnescapeString(text)
	t = r.Replace(t)
	t = codePattern.ReplaceAllString(t, "<pre>")
	t = aHrefPattern.ReplaceAllString(t, "\n - $2")
	return t
}

func openFile(path string) (*os.File, func(), error) {
	f, err := os.Open(filepath.FromSlash(fmt.Sprintf("testdata/%s.json", path)))
	if err != nil {
		return nil, nil, err
	}
	return f, func() { f.Close() }, nil
}

func getText(conf *Config) (*stackOverlowResult, error) {

	url := fmt.Sprintf("https://www.googleapis.com/customsearch/v1?key=%s&cx=%s&q=%s",
		conf.ApiKey, conf.SearchEngine, netUrl.QueryEscape(conf.Query))
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	res, err := conf.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode > 299 {
		return nil, fmt.Errorf("failed connecting to Google API: %s", res.Status)
	}
	var gsResp googleSearchResult
	err = json.NewDecoder(res.Body).Decode(&gsResp)
	// var gsResp googleSearchResult
	// f, close, err := openFile("goso")
	// if err != nil {
	// 	return nil, err
	// }
	// defer close()
	// err = json.NewDecoder(f).Decode(&gsResp)
	if err != nil {
		return nil, err
	}
	var answers string
	for _, item := range gsResp.Items {
		u, _ := netUrl.Parse(item.Link)
		answers += strings.Split(u.Path, "/")[2]
		answers += ";"
	}
	answers = strings.TrimSuffix(answers, ";")
	url = fmt.Sprintf("https://api.stackexchange.com/2.3/questions/%s/answers?order=desc&sort=votes&site=stackoverflow&filter=withbody",
		netUrl.QueryEscape(answers))
	req, err = http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	res, err = conf.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode > 299 {
		return nil, fmt.Errorf("failed connecting to Stack Overflow API: %s", res.Status)
	}
	var soResp stackOverlowResult
	err = json.NewDecoder(res.Body).Decode(&soResp)
	// var soResp stackOverlowResult
	// f, close, err := openFile("answers")
	// if err != nil {
	// 	return stackOverlowResult{}, err
	// }
	// defer close()

	// err = json.NewDecoder(f).Decode(&soResp)
	if err != nil {
		return nil, err
	}
	return &soResp, nil
}

func GetAnswers(conf *Config) error {
	var sb strings.Builder
	style := styles.Get(conf.Style)
	if style == nil {
		style = styles.Fallback
	}
	formatter := formatters.Get("terminal16m")
	if formatter == nil {
		formatter = formatters.Fallback
	}
	lexer := lexers.Get(conf.Lexer)
	if lexer == nil {
		lexer = lexers.Fallback
	}
	resp, err := getText(conf)
	if err != nil {
		return err
	}
	for _, item := range resp.Items {
		t := fmtText(item.Body)
		codeStartIdx = strings.Index(t, codeStartTag)
		for codeStartIdx != -1 {
			codeEndIdx = strings.Index(t, codeEndTag)
			if codeEndIdx == -1 {
				break
			}
			iterator, err := lexer.Tokenise(nil, t[codeStartIdx+len(codeStartTag):codeEndIdx])
			if err != nil {
				return err
			}
			err = formatter.Format(&sb, style, iterator)
			if err != nil {
				return err
			}
			t = t[:codeStartIdx] + sb.String() + t[codeEndIdx+len(codeEndTag):]
			codeStartIdx = strings.Index(t, codeStartTag)
			sb.Reset()
		}

		t = strings.ReplaceAll(t, "<code>", "\033[32m")
		t = strings.ReplaceAll(t, "</code>", "\033[0m")

		println(t)
	}
	return nil
}