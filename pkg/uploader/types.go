package uploader

type File struct {
	Data     []byte
	Filename string
}

// Uploader is an interface for uploading files to a remote service
type Uploader interface {
	Upload(message string, file *File) error
}
