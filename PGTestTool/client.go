package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

func main() {

	postParam := url.Values{
		"params":     {"[8 1 18 4 49 50 51 50 24 0 32 2 42 3 50 51 50 48 232 1]\0"},
		"servername": {"pgLogin"},
		"funcname":   {"login"},
	}

	resp, err := http.PostForm("http://localhost:3001/start", postParam)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(string(body))
}
