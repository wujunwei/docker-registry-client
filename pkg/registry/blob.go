package registry

import (
	"io"
	"net/http"
	"net/url"

	"github.com/docker/distribution"
	digest "github.com/opencontainers/go-digest"
)

func (r *Registry) DownloadBlob(repository string, digest digest.Digest) (io.ReadCloser, error) {
	url := r.url("/v2/%s/blobs/%s", repository, digest)
	r.Logf("registry.blob.download url=%s repository=%s digest=%s", url, repository, digest)

	resp, err := r.Client.Get(url)
	if err != nil {
		return nil, err
	}

	return resp.Body, nil
}

func (r *Registry) UploadBlob(repository string, digest digest.Digest, content io.Reader) error {
	uploadURL, err := r.initiateUpload(repository)
	if err != nil {
		return err
	}
	q := uploadURL.Query()
	q.Set("digest", digest.String())
	uploadURL.RawQuery = q.Encode()

	r.Logf("registry.blob.upload url=%s repository=%s digest=%s", uploadURL, repository, digest)

	upload, err := http.NewRequest("PUT", uploadURL.String(), content)
	if err != nil {
		return err
	}
	upload.Header.Set("Content-Type", "application/octet-stream")

	_, err = r.Client.Do(upload)
	return err
}

func (r *Registry) HasBlob(repository string, digest digest.Digest) (bool, error) {
	checkURL := r.url("/v2/%s/blobs/%s", repository, digest)
	r.Logf("registry.blob.check url=%s repository=%s digest=%s", checkURL, repository, digest)

	resp, err := r.Client.Head(checkURL)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err == nil {
		return resp.StatusCode == http.StatusOK, nil
	}

	urlErr, ok := err.(*url.Error)
	if !ok {
		return false, err
	}
	httpErr, ok := urlErr.Err.(*HTTPStatusError)
	if !ok {
		return false, err
	}
	if httpErr.Response.StatusCode == http.StatusNotFound {
		return false, nil
	}

	return false, err
}

func (r *Registry) BlobMetadata(repository string, digest digest.Digest) (distribution.Descriptor, error) {
	checkURL := r.url("/v2/%s/blobs/%s", repository, digest)
	r.Logf("registry.blob.check url=%s repository=%s digest=%s", checkURL, repository, digest)

	resp, err := r.Client.Head(checkURL)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return distribution.Descriptor{}, err
	}

	return distribution.Descriptor{
		Digest: digest,
		Size:   resp.ContentLength,
	}, nil
}

func (r *Registry) initiateUpload(repository string) (*url.URL, error) {
	initiateURL := r.url("/v2/%s/blobs/uploads/", repository)
	r.Logf("registry.blob.initiate-upload url=%s repository=%s", initiateURL, repository)

	resp, err := r.Client.Post(initiateURL, "application/octet-stream", nil)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return nil, err
	}

	location := resp.Header.Get("Location")
	locationURL, err := url.Parse(location)
	if err != nil {
		return nil, err
	}
	return locationURL, nil
}
