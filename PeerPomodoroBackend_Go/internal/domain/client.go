package domain

type Client struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func NewClient(id string, name string) *Client {
	return &Client{
		ID:   id,
		Name: name,
	}
}
