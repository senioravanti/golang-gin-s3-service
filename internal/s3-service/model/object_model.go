package model

type ObjectResponse struct {
	Key string `json:"key"`
	Size int64 `json:"size"`
	ContentType string `json:"contentType"`
	LastModified string `json:"lastModified"`
}
