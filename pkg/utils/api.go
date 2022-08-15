package utils

import (
	"io"
	"net/http"
)

var VERIFY_HOST = "http://127.0.0.1:85"

func ReadBody(url string) ([]byte, error) {
	resp, err := http.Get(VERIFY_HOST + url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	return body, err
}
