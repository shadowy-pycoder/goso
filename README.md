# goso - [Stack Overflow](https://stackoverflow.com/) CLI Tool written in Go

![GitHub Actions Workflow Status](https://img.shields.io/github/actions/workflow/status/shadowy-pycoder/goso/go.yml)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/shadowy-pycoder/goso)
[![Go Reference](https://pkg.go.dev/badge/github.com/shadowy-pycoder/goso.svg)](https://pkg.go.dev/github.com/shadowy-pycoder/goso)
[![Go Report Card](https://goreportcard.com/badge/github.com/shadowy-pycoder/goso)](https://goreportcard.com/report/github.com/shadowy-pycoder/goso)
![GitHub Downloads (all assets, all releases)](https://img.shields.io/github/downloads/shadowy-pycoder/goso/total)


## Demo

![goso - Animated gif demo](demo/demo.gif)



## Installation

You can download the binary for your platform from [Releases](https://github.com/shadowy-pycoder/goso/releases) page.

Alternatively, you can install it using `go install` command (requires Go [1.23](https://go.dev/doc/install) or later):

```shell
CGO_ENABLED=0 go install -ldflags "-s -w" -trimpath github.com/shadowy-pycoder/goso/cmd/goso@latest
```
This will install the `goso` binary to your `$GOPATH/bin` directory.

If none of the above works for you, you can use the [Makefile](https://github.com/shadowy-pycoder/goso/blob/master/Makefile) to build the binary from source.

```shell
git clone https://github.com/shadowy-pycoder/goso.git
cd goso
make build
```


## Search Engine Setup

###  Google Search JSON API
This approach employs [Custom Search JSON API](https://developers.google.com/custom-search/v1/overview) from Google to obtain most relevant results from Stack Overflow. So, to make it work, you need to get an API key from Google and also a [Search Engine ID](https://developers.google.com/custom-search/v1/overview#search_engine_id). That gives you `100 requests per day`, which I believe is enough for most use cases.

Setup your `Search Engine ID` like this:

![Screenshot from 2024-11-12 13-17-26](https://github.com/user-attachments/assets/3dd798fb-d9de-438a-aeeb-81ffc47e488b)

Add variables to your environment:
```shell
echo "export GOSO_API_KEY=<YOUR_API_KEY>" >> $HOME/.profile
echo "export GOSO_SE=<YOUR_SEARCH_ENGINE_ID>" >> $HOME/.profile
source $HOME/.profile
```

###  OpenSerp API
`goso` also supports [OpenSERP (Search Engine Results Page)](https://github.com/karust/openserp) from [Karust](https://github.com/karust). This is a completely *FREE* alternative to the Google Search JSON API, though it works a little bit slower, but gives basically the same results. 

So, to make it work, you need to run OperSERP server locally. You can do it like this:

With Docker:
```shell
docker run -p 127.0.0.1:7000:7000 -it karust/openserp serve -a 0.0.0.0 -p 7000
```

Or as a CLI command:

```shell
openserp serve 
```
You can learn more on how to install OpenSERP [here](https://github.com/karust/openserp).

Once you have it running, add variables to your environment:
```shell
echo "export GOSO_OS_HOST=127.0.0.1" >> $HOME/.profile
echo "export GOSO_OS_PORT=7000" >> $HOME/.profile
source $HOME/.profile
```
These variables will have priority over the `GOSO_API_KEY` and `GOSO_SE`.


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

You can also use `less` command for instance to page through the results:
```shell  
#!/bin/bash
goso "$@" | less -F -R -X
```

