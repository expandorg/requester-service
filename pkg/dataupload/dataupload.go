package dataupload

import (
	"errors"
	"mime/multipart"

	"github.com/expandorg/requester-service/pkg/app"
	m "github.com/expandorg/requester-service/pkg/model"
	"github.com/expandorg/requester-service/pkg/svc-kit/svc"
	"github.com/expandorg/requester-service/pkg/upload"
)

type Service interface {
	Upload(ctx *app.ReqCtx, userID uint64, draft *m.Draft, fileHeader *multipart.FileHeader) (*m.Draft, error)
}

type dataUploadService struct {
	drafts   m.DraftRepository
	data     m.DataRepository
	importer *upload.CSVImporter
	storage  upload.Storage
}

// New returns a new instance of a ata upload service
func New(drafts m.DraftRepository, data m.DataRepository, storage upload.Storage, importer *upload.CSVImporter) Service {
	return &dataUploadService{
		drafts:   drafts,
		data:     data,
		importer: importer,
		storage:  storage,
	}
}

func (ds *dataUploadService) Upload(ctx *app.ReqCtx, userID uint64, draft *m.Draft, fileHeader *multipart.FileHeader) (*m.Draft, error) {
	up, err := upload.NewUpload(fileHeader)
	if err != nil {
		return nil, err
	}

	if up.MIME != upload.CSV {
		return nil, svc.ArgumentsErr(errors.New("Uploaded file must be a CSV"))
	}

	urlChan := make(chan fuChan)

	var data *upload.Data
	go func(ec chan fuChan) {
		defer up.Close()
		data, err = ds.importer.Import(up)
		ec <- fuChan{"", err}
	}(urlChan)

	go func(ec chan fuChan) {
		u, e := ds.storage.UploadTaskData(up, userID, draft.ID)
		ec <- fuChan{u, e}
	}(urlChan)

	for i := 0; i < 2; i++ {
		fileUpload := <-urlChan
		if fileUpload.Error != nil {
			return nil, err
		}
	}

	return ds.updateData(ctx, draft, data)
}

func (ds *dataUploadService) updateData(ctx *app.ReqCtx, draft *m.Draft, d *upload.Data) (*m.Draft, error) {
	dataID, err := ds.data.Create(ctx, draft.ID, d)
	if err != nil {
		return nil, err
	}
	draft.DataID = dataID
	err = ds.drafts.Update(ctx, draft)
	if err != nil {
		return nil, err
	}
	return ds.drafts.FindByID(ctx, draft.ID)
}

type fuChan struct {
	FileName string
	Error    error
}
