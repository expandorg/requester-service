package thumbnailupload

import (
	"errors"
	"mime/multipart"

	"github.com/expandorg/requester-service/pkg/svc-kit/svc"
	"github.com/expandorg/requester-service/pkg/upload"
	"github.com/globalsign/mgo/bson"
)

type Service interface {
	Upload(userID uint64, fileHeader *multipart.FileHeader) (string, error)
}

type thumbnailUploader struct {
	storage  upload.Storage
	importer upload.ImgImporter
}

// New returns a new instance of a thumbnail upload service.
func New(storage upload.Storage, importer upload.ImgImporter) Service {
	return &thumbnailUploader{
		storage:  storage,
		importer: importer,
	}
}

func (tc *thumbnailUploader) Upload(userID uint64, fileHeader *multipart.FileHeader) (string, error) {
	up, err := upload.NewUpload(fileHeader)
	if err != nil {
		return "", err
	}

	if up.MIME != upload.PNG && up.MIME != upload.JPEG && up.MIME != upload.JPG {
		return "", svc.ArgumentsErr(errors.New("Uploaded file must be a png, jpg or jpeg image file"))
	}

	urlChan := make(chan fuChan)

	go func(ec chan fuChan) {
		defer up.Close()
		err = tc.importer.Import(up)
		ec <- fuChan{"", err}
	}(urlChan)

	go func(ec chan fuChan) {
		u, e := tc.storage.UploadTaskData(up, userID, bson.NewObjectId())
		ec <- fuChan{u, e}
	}(urlChan)

	var urlHost = "https://storage.cloud.google.com/task-data-uploads/"
	var fileName string

	for i := 0; i < 2; i++ {

		fileUpload := <-urlChan
		if fileUpload.Error != nil {
			return "", err
		}
		if fileUpload.FileName != "" {
			fileName = fileUpload.FileName
		}
	}
	return urlHost + fileName, nil
}

type fuChan struct {
	FileName string
	Error    error
}
