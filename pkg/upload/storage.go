package upload

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"time"

	"cloud.google.com/go/storage"
	"github.com/globalsign/mgo/bson"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
)

const (
	taskDataStorageBucket = "task-data-uploads"
)

type Storage interface {
	UploadTaskData(*Upload, uint64, bson.ObjectId) (string, error)
}

type FileSystem struct{}

func (s *FileSystem) UploadTaskData(up *Upload, userID uint64, draftID bson.ObjectId) (string, error) {
	filename := makeFilename(userID, draftID)
	destination, err := os.Create(filename)
	if err != nil {
		return "", err
	}

	_, err = io.Copy(destination, up.StorageReader)
	return filename, err
}

type GCloud struct{}

func (s *GCloud) UploadTaskData(up *Upload, userID uint64, draftID bson.ObjectId) (string, error) {
	ctx := context.Background()

	// load credentials
	data, err := ioutil.ReadFile("./private.json")
	if err != nil {
		return "", err
	}

	creds, err := google.CredentialsFromJSON(ctx, data, "https://www.googleapis.com/auth/devstorage.read_write")
	if err != nil {
		return "", err
	}

	client, err := storage.NewClient(ctx, option.WithCredentials(creds))
	if err != nil {
		return "", err
	}

	bucket := client.Bucket(taskDataStorageBucket)
	filename := makeFilename(userID, draftID)
	object := bucket.Object(filename)
	w := object.NewWriter(ctx)
	defer w.Close()

	_, err = io.Copy(w, up.StorageReader)
	return filename, err
}

func makeFilename(userID uint64, draftID bson.ObjectId) string {
	return fmt.Sprintf("%d-%s-%d", userID, draftID.Hex(), time.Now().Unix())
}
