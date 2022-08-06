package clients

import (
	"fmt"
	"io"
	"net/url"

	"github.com/fmartingr/bazaar/pkg/clients/mockdata"
)

type MockClient struct{}

func (c MockClient) Get(urlString string) (io.Reader, error) {
	parsedUrl, err := url.Parse(urlString)
	if err != nil {
		return nil, fmt.Errorf("error parsing url: %s", urlString)
	}

	f, err := mockdata.Data.Open(parsedUrl.Host + ".html")
	if err != nil {
		return nil, fmt.Errorf("can't open mock data for %s", parsedUrl.Host)
	}

	return f, nil
}

func NewMockClient() MockClient {
	return MockClient{}
}
