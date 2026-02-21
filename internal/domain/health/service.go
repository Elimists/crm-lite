package health

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

func (s *svc) GetHealth() Health {
	return Health{
		Name:    s.name,
		Slug:    s.slug,
		Version: s.version,
		Status:  "running",
	}
}
