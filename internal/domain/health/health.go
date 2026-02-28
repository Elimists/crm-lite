package health

type Health struct {
	Name    string `json:"name"`
	Slug    string `json:"slug"`
	Version string `json:"version"`
	Status  string `json:"status"`
}

type Service interface {
	GetHealth() Health
}
