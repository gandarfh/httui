package model

type Endpoint struct {
	Path    string `yaml:"path"`
	Method  string `yaml:"method"`
	Headers string `yaml:"headers"`
	Body    string `yaml:"body"`
}

func (e *Endpoint) UpdateEndpoint(Endpoint) {}
func (e *Endpoint) DeleteEndpoint(Endpoint) {}
