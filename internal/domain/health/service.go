package health

// Response is the data structure the API will return.
// We use JSON tags here so the Service defines the "Contract".
type Response struct {
	Name    string `json:"name"`
	Slug    string `json:"slug"`
	Version string `json:"version"`
	Status  string `json:"status"`
}

type Service interface {
	GetHealth() Response
}

type svc struct {
	name    string
	slug    string
	version string
}

func NewService(name, slug, version string) Service {
	return &svc{
		name:    name,
		slug:    slug,
		version: version,
	}
}

func (s *svc) GetHealth() Response {
	return Response{
		Name:    s.name,
		Slug:    s.slug,
		Version: s.version,
		Status:  "running",
	}
}
