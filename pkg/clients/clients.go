package clients

import "io"

type Client interface {
	Get(url string) (io.Reader, error)
}
