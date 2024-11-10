package goso

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
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

func openFile(path string) (*os.File, func(), error) {
	f, err := os.Open(filepath.FromSlash(fmt.Sprintf("testdata/%s.json", path)))
	if err != nil {
		return nil, nil, err
	}
	return f, func() { f.Close() }, nil
}

func GetAnswers(client *http.Client) (GoogleSearchResult, error) {

	// apiKey := os.Getenv("GOOGLE_API_KEY")
	// se := os.Getenv("GOOGLE_SE")
	// query := "Make query in Golang"
	// url := fmt.Sprintf("https://www.googleapis.com/customsearch/v1?key=%s&cx=%s&q=%s",
	// 	apiKey, se, url.QueryEscape(query))
	// req, err := http.NewRequest(http.MethodGet, url, nil)
	// if err != nil {
	// 	return GoogleSearchResult{}, err
	// }
	// res, err := client.Do(req)
	// if err != nil {
	// 	return GoogleSearchResult{}, err
	// }
	// defer res.Body.Close()
	// if res.StatusCode > 299 {
	// 	return GoogleSearchResult{}, err
	// }
	// var gsResp GoogleSearchResult
	// err = json.NewDecoder(res.Body).Decode(&gsResp)
	var gsResp GoogleSearchResult
	f, close, err := openFile("goso")
	if err != nil {
		return GoogleSearchResult{}, err
	}
	defer close()
	err = json.NewDecoder(f).Decode(&gsResp)
	if err != nil {
		return GoogleSearchResult{}, err
	}
	return gsResp, nil
}
