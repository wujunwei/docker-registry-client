package main

import (
	"fmt"
	"github.com/docker/distribution/manifest/manifestlist"
	"github.com/heroku/docker-registry-client/pkg/registry"
)

var source *registry.Registry
var destination *registry.Registry

func init() {
	source, _ = registry.New("https://docker.io", "testuser", "testpassword")
	destination, _ = registry.New("https://docker.io", "testuser", "testpassword")
}

func main() {
	fmt.Println(Copy("demo", "1.0.0"))
}

// Copy copy one image reference may be name or digest of the image
func Copy(name, reference string) error {
	src, err := source.Manifest(name, reference)
	if err != nil {
		return err
	}
	metaType, _, _ := src.Payload()
	// multiply architecture
	if metaType == manifestlist.MediaTypeManifestList {
		for _, manifest := range src.References() {
			fmt.Printf("start to copy %s \n", manifest.Digest.String())
			err = Copy(name, manifest.Digest.String())
			if err != nil {
				return err
			}
		}

		fmt.Println("start to copy multi manifest list")
		err = destination.PutManifest(name, reference, src)
		if err != nil {
			return err
		}
		fmt.Println("successfully !")
		return nil
	}
	for _, layer := range src.References() {
		blob, _ := destination.HasBlob(name, layer.Digest)
		if blob {
			fmt.Printf("manifest %s exist skip upload.\n", layer.Digest)
			continue
		}
		downloadBlob, err := source.DownloadBlob(name, layer.Digest)
		if err != nil {
			return err
		}
		err = destination.UploadBlob(name, layer.Digest, downloadBlob)
		if err != nil {
			return err
		}
		fmt.Printf("upload blob %s succeed.\n", layer.Digest)
		_ = downloadBlob.Close()
	}
	return destination.PutManifest(name, reference, src)
}
