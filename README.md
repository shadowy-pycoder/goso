# goso - [Stack Overflow](https://stackoverflow.com/) CLI Tool written in Go

![GitHub Actions Workflow Status](https://img.shields.io/github/actions/workflow/status/shadowy-pycoder/goso/go.yml)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/shadowy-pycoder/goso)
[![Go Reference](https://pkg.go.dev/badge/github.com/shadowy-pycoder/goso.svg)](https://pkg.go.dev/github.com/shadowy-pycoder/goso)



## Installation

```shell
CGO_ENABLED=0 go install -ldflags "-s -w" -trimpath github.com/shadowy-pycoder/goso/cmd/goso@latest
```
This will install the `goso` binary to your `$GOPATH/bin` directory.

This tool uses [Custom Search JSON API](https://developers.google.com/custom-search/v1/overview) from Google to obtain most relevant results from Stack Overflow. So, to make it work, you need to obtain an API key from Google and also a [Search Engine ID](https://developers.google.com/custom-search/v1/overview#search_engine_id).

Setup your `Search Engine ID` like this:

![Screenshot from 2024-11-12 13-17-26](https://github.com/user-attachments/assets/3dd798fb-d9de-438a-aeeb-81ffc47e488b)

Add variables to your environment:
```shell
echo "export GOSO_API_KEY=<YOUR_API_KEY>" >> $HOME/.profile
echo "export GOSO_SE=<YOUR_SEARCH_ENGINE_ID>" >> $HOME/.profile
source $HOME/.profile
```
## Usage

```shell
goso -h
                                                                  
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
  -a int
        The maximum number of answers for each result [min=1, max=10] (default 1)
  -l string
        The name of Chroma lexer. See https://github.com/alecthomas/chroma/tree/master/lexers/embedded (default "bash")
  -q int
        The number of results [min=1, max=10] (default 1)
  -s string
        The name of Chroma style. See https://xyproto.github.io/splash/docs/ (default "onedark")
``` 

It is possible to adjust default values for the number of questions and answers, lexer and style.
```shell
echo "export GOSO_LEXER=python" >> $HOME/.profile
echo "export GOSO_STYLE=onedark" >> $HOME/.profile
echo "export GOSO_ANSWERS=5" >> $HOME/.profile
echo "export GOSO_QUESTIONS=5" >> $HOME/.profile
source $HOME/.profile
```

## Example

```shell
goso  -l go -s onedark -q 1 -a 1 Sort maps in Golang
```
Output:

![Screenshot from 2024-11-14 10-16-52](https://github.com/user-attachments/assets/43282839-1719-44ae-a0e8-c2ed44d8e9e6)

