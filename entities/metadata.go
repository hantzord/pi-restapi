package entities

type Metadata struct {
	Page  int `json:"page"`
	Limit int `json:"limit"`
}

type MetadataFull struct {
	Page   int    `json:"page"`
	Limit  int    `json:"limit"`
	Sort   string `json:"sort"`
	Order  string `json:"order"`
	Search string `json:"search"`
}

func (metadata *Metadata) Offset() int {
	return (metadata.Page - 1) * metadata.Limit
}
