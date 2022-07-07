package model

type Uri struct {
	Alias     string      `yaml:"alias"`
	Endpoints *[]Endpoint `yaml:"endpoints"`
}

func (u *Uri) ListEndpoints() *[]Endpoint {

	return u.Endpoints
}

func (u *Uri) CreateEndpoint(in *Endpoint) {

}
func (u *Uri) UpdateUri() {}
func (u *Uri) DeleteUri() {}
