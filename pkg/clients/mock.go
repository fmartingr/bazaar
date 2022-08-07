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

func (c MockClient) Get(u *url.URL) (io.Reader, error) {
	f, err := mockdata.Data.Open(u.Host + ".html")
	if err != nil {
		return nil, fmt.Errorf("can't open mock data for %s", u.Host)
	}

	return f, nil
}

func NewMockClient() MockClient {
	return MockClient{}
}
