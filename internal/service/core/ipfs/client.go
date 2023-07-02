package ipfs

import "net/http"

type AuthHandler interface {
	RoundTrip(r *http.Request) (*http.Response, error)
}

type authTransport struct {
	http.RoundTripper
	ProjectId     string
	ProjectSecret string
}

func NewClient(projectId, projectSecret string) *http.Client {
	return &http.Client{
		Transport: authTransport{
			RoundTripper:  http.DefaultTransport,
			ProjectId:     projectId,
			ProjectSecret: projectSecret,
		},
	}
}

func (t authTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	r.SetBasicAuth(t.ProjectId, t.ProjectSecret)
	return t.RoundTripper.RoundTrip(r)
}
