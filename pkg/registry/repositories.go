package registry

type repositoriesResponse struct {
	Repositories []string `json:"repositories"`
}

func (r *Registry) Repositories() ([]string, error) {
	url := r.url("/v2/_catalog")
	repos := make([]string, 0, 10)
	var err error //We create this here, otherwise url will be rescoped with :=
	var response repositoriesResponse
	for {
		r.Logf("registry.repositories url=%s", url)
		url, err = r.getPaginatedJSON(url, &response)
		switch err {
		case ErrNoMorePages:
			repos = append(repos, response.Repositories...)
			return repos, nil
		case nil:
			repos = append(repos, response.Repositories...)
			continue
		default:
			return nil, err
		}
	}
}
