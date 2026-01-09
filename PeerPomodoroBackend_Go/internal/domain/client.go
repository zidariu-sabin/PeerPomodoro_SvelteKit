package domain

type Client struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func newClient(id string, name string) *Client {
	return &Client{
		ID:   id,
		Name: name,
	}
}
