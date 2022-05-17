package registry

import (
	"bytes"
	"github.com/docker/distribution"
	"github.com/docker/distribution/manifest/manifestlist"
	"github.com/docker/distribution/manifest/schema1"
	"github.com/docker/distribution/manifest/schema2"
	digest "github.com/opencontainers/go-digest"
	"io/ioutil"
	"net/http"
)

var DefaultRequestedManifestMIMETypes = []string{
	schema1.MediaTypeManifest,
	schema2.MediaTypeManifest,
	schema1.MediaTypeSignedManifest,
	manifestlist.MediaTypeManifestList,
}

func (r *Registry) Manifest(repository, reference string) (distribution.Manifest, error) {
	url := r.url("/v2/%s/manifests/%s", repository, reference)
	r.Logf("registry.manifest.get url=%s repository=%s reference=%s", url, repository, reference)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	for _, mimeType := range DefaultRequestedManifestMIMETypes {
		req.Header.Add("Accept", mimeType)
	}
	resp, err := r.Client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	m, _, err := distribution.UnmarshalManifest(resp.Header.Get("Content-Type"), body)
	return m, err
}

func (r *Registry) ManifestDigest(repository, reference string) (digest.Digest, error) {
	url := r.url("/v2/%s/manifests/%s", repository, reference)
	r.Logf("registry.manifest.head url=%s repository=%s reference=%s", url, repository, reference)

	resp, err := r.Client.Head(url)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return "", err
	}
	return digest.Parse(resp.Header.Get("Docker-Content-Digest"))
}

func (r *Registry) DeleteManifest(repository string, digest digest.Digest) error {
	url := r.url("/v2/%s/manifests/%s", repository, digest)
	r.Logf("registry.manifest.delete url=%s repository=%s reference=%s", url, repository, digest)

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return err
	}
	resp, err := r.Client.Do(req)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return err
	}
	return nil
}

func (r *Registry) PutManifest(repository, reference string, manifest distribution.Manifest) error {
	url := r.url("/v2/%s/manifests/%s", repository, reference)
	r.Logf("registry.manifest.put url=%s repository=%s reference=%s", url, repository, reference)

	mediaType, payload, err := manifest.Payload()
	if err != nil {
		return err
	}

	buffer := bytes.NewBuffer(payload)
	req, err := http.NewRequest("PUT", url, buffer)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", mediaType)
	resp, err := r.Client.Do(req)
	if resp != nil {
		defer resp.Body.Close()
	}
	return err
}
