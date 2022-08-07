package clients

import (
	"fmt"
	"io"
	"net/url"

	"github.com/fmartingr/bazaar/internal/mockdata"
)

// MockClient A simple client used for test the shops which will load an HTML from the mockdata present
// in this same package based on the requested host.
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
