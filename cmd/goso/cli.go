package main

import (
	"crypto/tls"
	"net/http"
	"os"
	"time"

	"github.com/shadowy-pycoder/goso"
)

const app string = "goso"

func root(args []string) error {
	conf := &goso.Config{
		ApiKey:       os.Getenv("GOOGLE_API_KEY"),
		SearchEngine: os.Getenv("GOOGLE_SE"),
		Query:        "How to create CLI tool in Golang",
		Style:        "onedark",
		Lexer:        "go",
		QuestionNum:  1,
		AnswerNum:    1,
		Client: &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			},
			Timeout: time.Duration(10) * time.Second,
		},
	}
	return goso.GetAnswers(conf)
}
