package main

import (
	"errors"
	"fmt"
	"golang.org/x/net/html"
	"net/http"
	"os"
)

// gets url as the first cmd argument and parses it
func parseUrl() (*html.Node, error) {
	if len(os.Args[1:]) == 0 {
		return nil, errors.New("no URL provided")
	}

	res, err := http.Get(os.Args[1])
	if err != nil {
		return nil, err
	}

	if res.StatusCode > 300 {
		return nil, fmt.Errorf("bad status code: %d", res.StatusCode)
	}

	defer func(res *http.Response) {
		err = res.Body.Close()
	}(res)

	return html.Parse(res.Body)
}
