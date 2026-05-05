package entity

type UploadInput struct {
	OwnerType    string
	OwnerID      string
	Purpose      MediaPurpose
	FileName     string
	ContentType  string
	Data         []byte
	UploadedBy   string
}

type UploadResult struct {
	ID     string
	Status string
	URLs   map[string]string
}