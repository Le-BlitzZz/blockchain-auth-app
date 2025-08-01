package session

import "net/http"

var httpClient = http.DefaultClient

func SetHTTPClient(c *http.Client) {
	httpClient = c
}
