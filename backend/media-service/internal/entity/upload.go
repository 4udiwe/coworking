package entity

type UploadInput struct {
	FileName     string
	ContentType  string
	Data         []byte
}

type UploadResult struct {
	ID     string
	Status string
	URLs   map[string]string
}