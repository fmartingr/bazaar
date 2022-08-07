package clients

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type HttpClient struct {
	http http.Client
}

func (c HttpClient) Get(u *url.URL) (io.Reader, error) {
	res, err := c.http.Get(u.String())
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
