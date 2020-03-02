package upload

import (
	"io"
	"mime/multipart"
)

type MIME string

func (m MIME) Equals(n string) bool {
	return string(m) == n
}

const CSV MIME = "text/csv"
const JPG MIME = "image/jpg"
const PNG MIME = "image/png"
const JPEG MIME = "image/jpeg"


type Upload struct {
	MIME          MIME
	ImportReader  io.Reader
	StorageReader io.Reader
	PipeCloser    io.Closer
	FileCloser    io.Closer
}

func NewUpload(header *multipart.FileHeader) (*Upload, error) {
	mime := MIME(header.Header.Get("Content-Type"))

	file, err := header.Open()
	if err != nil {
		return nil, err
	}

	pipeReader, pipeWriter := io.Pipe()
	teeReader := io.TeeReader(file, pipeWriter)

	upload := &Upload{
		MIME:          mime,
		ImportReader:  teeReader,
		StorageReader: pipeReader,
		PipeCloser:    pipeWriter,
		FileCloser:    file,
	}
	return upload, nil
}

func (u *Upload) Close() {
	u.PipeCloser.Close()
	u.FileCloser.Close()
}
