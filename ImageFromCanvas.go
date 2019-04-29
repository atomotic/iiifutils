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

func ReadCanvas(canvas string) (*Canvas, error) {
	var b []byte

	b, err := cache.Read(hash(canvas))

	if err != nil {
		resp, err := http.Get(canvas)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()
		if resp.StatusCode != 404 {
			b, err = ioutil.ReadAll(resp.Body)
			if err != nil {
				return nil, err
			}

			cache.Write(hash(canvas), b)
		} else {
			return nil, errors.New("404 canvas")
		}

	}

	var c Canvas
	err = json.Unmarshal(b, &c)
	if err != nil {
		return nil, err
	}

	return &c, nil
}

func ImageFromCanvas(manifestURL string, canvasURL string) (string, error) {
	var image string

	c, err := ReadCanvas(canvasURL)

	// if errors in reading Canvas (uri not resolvable) or if Images[] empty
	// then read the whole manifest
	if err != nil || len(c.Images) <= 0 {
		m, err := ReadManifest(manifestURL)
		if err != nil {
			return "", err
		}
		for _, canvas := range m.Sequences[0].Canvases {
			if canvas.ID == canvasURL {
				image = canvas.Images[0].Resource.Service.ID

			}
		}
	} else {
		image = c.Images[0].Resource.ID
	}

	return image, nil
}
