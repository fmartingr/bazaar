package clients

import (
	"fmt"
	"io"
	"net/http"
)

type HttpClient struct {
	http http.Client
}

func (c HttpClient) Get(url string) (io.Reader, error) {
	res, err := c.http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error retrieving url: %s", err)
	}

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("error retrieving url: %d %s", res.StatusCode, res.Status)
	}

	return res.Body, nil
}

func NewBasicHttpClient() HttpClient {
	return HttpClient{
		http: http.Client{},
	}
}
