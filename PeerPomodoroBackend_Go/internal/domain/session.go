package domain

type Session struct {
	ID      string   `json:"id"`
	Clients []Client `json:"clients"`
	Timer   *Timer   `json:"timer"`
}


