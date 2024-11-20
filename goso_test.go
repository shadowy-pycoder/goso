package goso

import (
	"encoding/json"
	"fmt"
	"maps"
	netUrl "net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"
)

func openFile(path string) (*os.File, func(), error) {
	f, err := os.Open(filepath.FromSlash(fmt.Sprintf("testdata/%s.json", path)))
	if err != nil {
		return nil, nil, err
	}
	return f, func() { f.Close() }, nil
}

func fetchGoogle(conf *Config, results map[int]*Result) error {
	var gsResp GoogleSearchResult
	f, close, err := openFile("goso")
	if err != nil {
		return err
	}
	defer close()
	err = json.NewDecoder(f).Decode(&gsResp)
	if err != nil {
		return err
	}
	for _, item := range gsResp.Items {
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
		questionId, _ := strconv.Atoi(strings.Split(u.Path, "/")[2])
		results[questionId] = &Result{
			Title:       item.Title,
			Link:        item.Link,
			QuestionId:  questionId,
			UpvoteCount: upvoteCount,
			Date:        dateCreated,
		}
	}
	return nil
}

func fetchStackOverflow(conf *Config, results map[int]*Result) error {
	questions := make([]string, len(results))
	var idx int
	for question := range maps.Keys(results) {
		questions[idx] = strconv.Itoa(question)
		idx++
	}
	_ = strings.Join(questions, ";")
	var soResp StackOverflowResult
	f, close, err := openFile("answers")
	if err != nil {
		return err
	}
	defer close()

	err = json.NewDecoder(f).Decode(&soResp)
	if err != nil {
		return err
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
	return nil
}

func BenchmarkGetAnswers(b *testing.B) {
	b.ResetTimer()
	conf := &Config{
		Style:       "onedark",
		Lexer:       "c",
		QuestionNum: 10,
		AnswerNum:   10,
	}
	for i := 0; i < b.N; i++ {
		_, err := GetAnswers(conf, fetchGoogle, fetchStackOverflow)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func TestOutput(t *testing.T) {
	conf := &Config{
		Style:       "onedark",
		Lexer:       "c",
		QuestionNum: 10,
		AnswerNum:   10,
	}
	answers, err := GetAnswers(conf, fetchGoogle, fetchStackOverflow)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(answers)
}
