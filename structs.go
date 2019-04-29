package iiifutils

type Canvas struct {
	ID     string `json:"@id"`
	Images []struct {
		ID       string `json:"@id"`
		Resource struct {
			ID      string `json:"@id"`
			Service struct {
				ID string `json:"@id"`
			} `json:"service"`
		} `json:"resource"`
	} `json:"images"`
}

type Manifest struct {
	Sequences []struct {
		Canvases []Canvas `json:"canvases"`
		ID       string   `json:"@id"`
		Label    string   `json:"label"`
		Type     string   `json:"@type"`
	} `json:"sequences"`
}
