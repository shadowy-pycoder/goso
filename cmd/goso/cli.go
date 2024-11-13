package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strconv"
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
	var (
		lex, style string
		qn, an     int
		set        bool
		err        error
	)
	lex, set = os.LookupEnv("GOSO_LEXER")
	if !set {
		lex = "bash"
	}
	style, set = os.LookupEnv("GOSO_STYLE")
	if !set {
		style = "onedark"
	}
	q, set := os.LookupEnv("GOSO_QUESTIONS")
	if !set {
		qn = 1
	} else {
		qn, err = strconv.Atoi(q)
		if err != nil {
			return fmt.Errorf("-q should be within [min=1, max=10], please check your env")
		}
		if qn < 1 || qn > 10 {
			return fmt.Errorf("-q should be within [min=1, max=10], please check your env")
		}
	}
	a, set := os.LookupEnv("GOSO_ANSWERS")
	if !set {
		an = 1
	} else {
		an, err = strconv.Atoi(a)
		if err != nil {
			return fmt.Errorf("-a should be within [min=1, max=10], please check your env")
		}
		if an < 1 || an > 10 {
			return fmt.Errorf("-a should be within [min=1, max=10], please check your env")
		}
	}
	flags := flag.NewFlagSet(app, flag.ExitOnError)
	flags.StringVar(&conf.Lexer, "l", lex, "The name of Chroma lexer. See https://github.com/alecthomas/chroma/tree/master/lexers/embedded")
	flags.StringVar(&conf.Style, "s", style, "The name of Chroma style. See https://xyproto.github.io/splash/docs/")
	qNum := flags.Int("q", qn, "The number of questions [min=1, max=10]")
	aNum := flags.Int("a", an, "The number of answers for each result [min=1, max=10]")

	flags.Usage = func() {
		fmt.Print(usagePrefix)
		flags.PrintDefaults()
	}

	if err := flags.Parse(args); err != nil {
		return err
	}
	if *qNum < 1 || *qNum > 10 {
		return fmt.Errorf("-q should be within [min=1, max=10]")
	}
	conf.QuestionNum = *qNum
	if *aNum < 1 || *aNum > 10 {
		return fmt.Errorf("-a should be within [min=1, max=10]")
	}
	conf.AnswerNum = *aNum
	apiKey, set := os.LookupEnv("GOSO_API_KEY")
	if !set {
		return fmt.Errorf("api key is not set")
	}
	conf.ApiKey = apiKey
	se, set := os.LookupEnv("GOSO_SE")
	if !set {
		return fmt.Errorf("search engine is not set")
	}
	conf.SearchEngine = se
	conf.Query = strings.Join(flags.Args(), " ")
	if conf.Query == "" {
		return fmt.Errorf("query is empty")
	}
	answers, err := goso.GetAnswers(conf, goso.FetchGoogle, goso.FetchStackOverflow)
	if err != nil {
		return err
	}
	fmt.Println(answers)
	return nil
}
