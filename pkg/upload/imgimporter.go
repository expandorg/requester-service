package upload

import (
	"io"
  "bufio"
)

type ImgImporter interface {
	Import(*Upload) (error)
}

type ImageImporter struct{}

func (i *ImageImporter) Import(up *Upload) (error) {
  imgReader := bufio.NewReader(up.ImportReader)
  // To Do: Handle file too large
  maxFileSize := 200000
  buffer := make([]byte, maxFileSize)
  io.ReadFull(imgReader, buffer);
  // To Do: Handle Error
  return nil
}
