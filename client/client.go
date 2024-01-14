package client

type Client struct {
	URL string
}

// NewClient returns a new instance of Client
func NewClient(u string) *Client {
	return &Client{
		URL: u,
	}
}
