package main

import (
	"crypto/tls"
	"fmt"
	"net/http"
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
	resp, err := goso.GetAnswers(client)
	if err != nil {
		println(err)
	}
	for i, item := range resp.Items {
		fmt.Printf("%d. %s\n\033[33m%s\033[0m\n", i+1, item.Title, item.Link)
	}
}
