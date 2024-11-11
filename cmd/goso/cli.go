package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/shadowy-pycoder/goso"
)

const app string = "goso"

const usagePrefix string = `                                                                  
 .d88b.   .d88b.  .d8888b   .d88b.  
d88P"88b d88""88b 88K      d88""88b 
888  888 888  888 "Y8888b. 888  888 
Y88b 888 Y88..88P      X88 Y88..88P 
 "Y88888  "Y88P"   88888P'  "Y88P"  
     888                            
Y8b d88P  Stack Overlow CLI Tool by shadowy-pycoder                         
 "Y88P"   GitHub: https://github.com/shadowy-pycoder/goso                        
                                                                                                                                                                                              
Usage: goso [OPTIONS] QUERY
Options:
  -h    Show this help message and exit.
`

func root(args []string) error {

	conf := &goso.Config{
		Client: &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			},
			Timeout: time.Duration(10) * time.Second,
		},
	}
	flags := flag.NewFlagSet(app, flag.ExitOnError)
	flags.StringVar(&conf.Lexer, "l", "bash", "The name of Chroma lexer. See https://github.com/alecthomas/chroma/tree/master/lexers/embedded")
	flags.StringVar(&conf.Style, "s", "onedark", "The name of Chroma style. See https://xyproto.github.io/splash/docs/")
	qNum := flags.Int("q", 1, "The number of results [min=1, max=10]")
	if *qNum < 1 || *qNum > 10 {
		*qNum = 1
	}
	conf.QuestionNum = *qNum
	aNum := flags.Int("a", 1, "The maximum number of answers for each result [min=1, max=10]")
	if *aNum < 1 || *aNum > 10 {
		*aNum = 1
	}
	conf.AnswerNum = *aNum
	flags.Usage = func() {
		fmt.Print(usagePrefix)
		flags.PrintDefaults()
	}

	if err := flags.Parse(args); err != nil {
		return err
	}
	apiKey, set := os.LookupEnv("GOOGLE_API_KEY")
	if !set {
		return fmt.Errorf("google api key is not set")
	}
	conf.ApiKey = apiKey
	se, set := os.LookupEnv("GOOGLE_SE")
	if !set {
		return fmt.Errorf("google search engine is not set")
	}
	conf.SearchEngine = se
	conf.Query = strings.Join(flags.Args(), " ")
	if conf.Query == "" {
		return fmt.Errorf("query is empty")
	}
	return goso.GetAnswers(conf)
}
