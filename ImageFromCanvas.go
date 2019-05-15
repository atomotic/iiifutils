package iiifutils

import (
	"encoding/json"
	"errors"
	"hash/fnv"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/peterbourgon/diskv"
)

var cache *diskv.Diskv

func init() {
	flatTransform := func(s string) []string { return []string{} }
	cache = diskv.New(diskv.Options{
		BasePath:  "/tmp/iiifutils/cache",
		Transform: flatTransform,
	})
}

func hash(s string) string {
	h := fnv.New64a()
	h.Write([]byte(s))
	return strconv.FormatUint(h.Sum64(), 16)
}

func ReadManifest(manifest string) (*Manifest, error) {
	var b []byte

	b, err := cache.Read(hash(manifest))

	if err != nil {
		resp, err := http.Get(manifest)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		if resp.StatusCode != 404 {
			b, err = ioutil.ReadAll(resp.Body)
			if err != nil {
				return nil, err
			}
			cache.Write(hash(manifest), b)
		} else {
			return nil, errors.New("404 manifest")
		}

	}
	var m Manifest
	err = json.Unmarshal(b, &m)
	if err != nil {
		return nil, err
	}

	return &m, nil
}

func ImageFromCanvas(manifestURL string, canvasURL string) (string, int, error) {
	var image string
	var page int

	m, err := ReadManifest(manifestURL)
	if err != nil {
		return "", 0, err
	}
	for sequence, canvas := range m.Sequences[0].Canvases {
		if canvas.ID == canvasURL {
			image = canvas.Images[0].Resource.Service.ID
			page = sequence + 1
		}
	}

	return image, page, nil
}
