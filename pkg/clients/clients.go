package clients

import (
	"io"
	"net/url"
)

type Client interface {
	Get(u *url.URL) (io.Reader, error)
}
