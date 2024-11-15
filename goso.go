package goso

import (
	"cmp"
	"encoding/json"
	"fmt"
	"html"
	"maps"
	"net/http"
	netUrl "net/url"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/alecthomas/chroma/v2/formatters"
	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/alecthomas/chroma/v2/styles"
	"golang.org/x/term"
)

const (
	codeStartTag     string = "<pre><code>"
	codeEndTag       string = "</code></pre>"
	reset            string = "\033[0m"
	bold             string = "\033[1m"
	italic           string = "\033[3m"
	strikethrough    string = "\033[9m"
	gray             string = "\033[37m"
	blue             string = "\033[36m"
	green            string = "\033[32m"
	yellow           string = "\033[33m"
	magenta          string = "\033[35m"
	questionColor    string = "\033[38;5;204m"
	answerColor      string = "\033[38;5;255m"
	downvoted        string = "\033[38;5;160m"
	lightgray        string = "\033[38;5;248m"
	urlColor         string = "\033[38;5;248m"
	terminalMaxWidth int    = 80
)

var (
	codeStartIdx  int
	codeEndIdx    int
	terminalWidth int
	// https://meta.stackexchange.com/questions/1777/what-html-tags-are-allowed-on-stack-exchange-sites
	codePattern  = regexp.MustCompile(`<pre\s.*?>`)
	aHrefPattern = regexp.MustCompile(`(?s)<a\s+(?:[^>]*?\s+)?href=(["'])?([^\'" >]+)(.*?)?</a>`)
	divPattern   = regexp.MustCompile(`<div.*?>`)
	bqPattern    = regexp.MustCompile(`<blockquote.*?>`)
	r            = strings.NewReplacer(
		"<li><p>", "",
		"<li><a href", "<a href",
		"<p>", "",
		"</p>", "",
		"<strong>", bold,
		"</strong>", reset,
		"<em>", italic,
		"</em>", reset,
		"<i>", italic,
		"</i>", reset,
		"<ul>", "",
		"</ul>", "",
		"<ol>", "",
		"</ol>", "",
		"<li>", " - ",
		"</li>", "",
		"<b>", bold,
		"</b>", reset,
		"<h1>", bold,
		"</h1>", reset,
		"<h2>", bold,
		"</h2>", reset,
		"<h3>", bold,
		"</h3>", reset,
		"<h4>", bold,
		"</h4>", reset,
		"<h5>", bold,
		"</h5>", reset,
		"<h6>", bold,
		"</h6>", reset,
		"<br>", "\n",
		"<blockquote>", italic,
		"</blockquote>", reset,
		"<del>", strikethrough,
		"</del>", reset,
		"<s>", strikethrough,
		"</s>", reset,
		"</div>", "",
		"<code>", green,
		"</code>", reset,
		"<br />", "",
		"<br/>", "",
		"<hr />", "",
		"<hr/>", "",
		"<sup>", "",
		"</sup>", "",
		"<sub>", "",
		"</sub>", "",
		"<dl>", "",
		"</dl>", "",
		"<dt>", "",
		"</dt>", "",
		"<dd>", " - ",
		"</dd>", "",
		"<kbd>", bold,
		"</kbd>", reset,
	)
)

type GoogleSearchResult struct {
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

type StackOverflowResult struct {
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
type Answer struct {
	Author     string
	Score      int
	Body       string
	Link       string
	IsAccepted bool
	Date       time.Time
}

func (a *Answer) String() string {
	line := strings.Repeat("─", terminalWidth)
	color := yellow
	if a.IsAccepted {
		color = green
	} else if a.Score < 0 {
		color = downvoted
	}
	return fmt.Sprintf(`
%s
%s[%d]%s %sAnswer from %s%s%s
%sDate: %s%s
%sLink: %s%s
%s

`,
		line,
		color, a.Score, reset, answerColor, bold, a.Author, reset,
		lightgray, a.Date.Format(time.RFC822), reset,
		lightgray, a.Link, reset,
		line)
}

type Result struct {
	Title       string
	Link        string
	QuestionId  int
	UpvoteCount int
	Date        time.Time
	Answers     []*Answer
}

func (r *Result) String() string {
	line := strings.Repeat("─", terminalWidth)
	color := yellow
	if r.UpvoteCount < 0 {
		color = downvoted
	}

	return fmt.Sprintf(`
%s
%s[%d]%s %s%s%s%s
%sDate: %s%s
%sLink: %s%s
%s`,
		line,
		color, r.UpvoteCount, reset, bold, questionColor, r.Title, reset,
		lightgray, r.Date.Format(time.RFC822), reset,
		lightgray, r.Link, reset,
		line)
}

func prepareText(text string) string {
	return codePattern.ReplaceAllString(text, "<pre>")
}

func fmtText(text string) string {
	t := r.Replace(html.UnescapeString(text))
	t = strings.ReplaceAll(t, "<hr>", strings.Repeat("─", terminalWidth))
	t = divPattern.ReplaceAllString(t, "")
	t = aHrefPattern.ReplaceAllString(t, fmt.Sprintf("\n %s- $2%s", urlColor, reset))
	t = bqPattern.ReplaceAllString(t, italic)
	return t
}

func FetchGoogle(conf *Config) (*GoogleSearchResult, error) {
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
	var gsResp GoogleSearchResult
	err = json.NewDecoder(res.Body).Decode(&gsResp)
	if err != nil {
		return nil, err
	}
	return &gsResp, nil
}

func FetchStackOverflow(conf *Config, gr *GoogleSearchResult) (map[int]*Result, error) {

	results := make(map[int]*Result)
	questions := make([]string, 0, len(gr.Items))
	for _, item := range gr.Items {
		var upvoteCount int
		var dateCreated time.Time
		if len(item.Pagemap.Question) > 0 {
			question := item.Pagemap.Question[0]
			answerCount, _ := strconv.Atoi(question.Answercount)
			if answerCount == 0 {
				continue
			}
			upvoteCount, _ = strconv.Atoi(question.Upvotecount)
			dateCreated, _ = time.Parse("2006-01-02T15:04:05", question.Datecreated)
		}
		u, _ := netUrl.Parse(item.Link)
		questionStr := strings.Split(u.Path, "/")[2]
		questions = append(questions, questionStr)
		questionId, _ := strconv.Atoi(questionStr)
		results[questionId] = &Result{
			Title:       item.Title,
			Link:        item.Link,
			QuestionId:  questionId,
			UpvoteCount: upvoteCount,
			Date:        dateCreated,
		}
	}
	question_ids := strings.Join(questions, ";")
	url := fmt.Sprintf("https://api.stackexchange.com/2.3/questions/%s/answers?order=desc&sort=votes&site=stackoverflow&filter=withbody",
		netUrl.QueryEscape(question_ids))
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
		return nil, fmt.Errorf("failed connecting to Stack Overflow API: %s", res.Status)
	}
	var soResp StackOverflowResult
	err = json.NewDecoder(res.Body).Decode(&soResp)
	if err != nil {
		return nil, err
	}
	for _, item := range soResp.Items {
		result, ok := results[item.QuestionID]
		if !ok {
			continue
		}
		result.Answers = append(result.Answers,
			&Answer{
				Author:     item.Owner.DisplayName,
				Score:      item.Score,
				Body:       item.Body,
				Link:       fmt.Sprintf("https://stackoverflow.com/a/%d", item.AnswerID),
				IsAccepted: item.IsAccepted,
				Date:       time.Unix(int64(item.CreationDate), 0).UTC(),
			})
	}
	return results, nil
}

func GetAnswers(conf *Config,
	fetchResults func(*Config) (*GoogleSearchResult, error),
	fetchAnswers func(*Config, *GoogleSearchResult) (map[int]*Result, error),
) (string, error) {
	var err error
	if term.IsTerminal(0) {
		terminalWidth, _, err = term.GetSize(0)
		if err != nil {
			return "", err
		}
		terminalWidth = min(terminalWidth, terminalMaxWidth)
	} else {
		terminalWidth = terminalMaxWidth
	}
	var answers strings.Builder
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
	gsResp, err := fetchResults(conf)
	if err != nil {
		return "", err
	}
	results, err := fetchAnswers(conf, gsResp)
	if err != nil {
		return "", err
	}
	var qIdx int
	for _, res := range slices.Backward(slices.SortedStableFunc(maps.Values(results), func(a, b *Result) int {
		return cmp.Compare(a.UpvoteCount, b.UpvoteCount)
	})) {
		if qIdx >= conf.QuestionNum {
			break
		}
		if len(res.Answers) == 0 {
			continue
		}
		qIdx++
		slices.SortStableFunc(res.Answers, func(a, b *Answer) int {
			return cmp.Compare(a.Score, b.Score)
		})
		answers.WriteString(res.String())
		var aIdx int
		for _, ans := range slices.Backward(res.Answers) {
			if aIdx >= conf.AnswerNum {
				break
			}
			aIdx++
			answers.WriteString(ans.String())
			t := prepareText(ans.Body)
			codeStartIdx = strings.Index(t, codeStartTag)
			if codeStartIdx == -1 {
				answers.WriteString(fmtText(t))
			}
			for codeStartIdx != -1 {
				codeEndIdx = strings.Index(t, codeEndTag)
				if codeEndIdx == -1 {
					break
				}
				answers.WriteString(fmtText(t[:codeStartIdx]))
				iterator, err := lexer.Tokenise(nil, html.UnescapeString(t[codeStartIdx+len(codeStartTag):codeEndIdx]))
				if err != nil {
					return "", err
				}
				err = formatter.Format(&answers, style, iterator)
				if err != nil {
					return "", err
				}
				t = t[codeEndIdx+len(codeEndTag):]
				codeStartIdx = strings.Index(t, codeStartTag)
				if codeStartIdx == -1 {
					answers.WriteString(fmtText(t))
				}
			}
		}
	}
	return answers.String(), nil
}
