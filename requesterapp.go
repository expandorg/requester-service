package main

import (
	"log"

	"github.com/expandorg/requester-service/mocks/templatemocks"
	"github.com/expandorg/requester-service/pkg/dashboard"
	"github.com/expandorg/requester-service/pkg/dataupload"
	"github.com/expandorg/requester-service/pkg/db"
	"github.com/expandorg/requester-service/pkg/draftservice"
	"github.com/expandorg/requester-service/pkg/drafttemplates"
	"github.com/expandorg/requester-service/pkg/http/app"
	"github.com/expandorg/requester-service/pkg/onboardingtemplate"
	"github.com/expandorg/requester-service/pkg/publisher"
	"github.com/expandorg/requester-service/pkg/svc-kit/cfg/env"
	"github.com/expandorg/requester-service/pkg/svc-kit/http/api"
	"github.com/expandorg/requester-service/pkg/svc-kit/mongo"
	"github.com/expandorg/requester-service/pkg/taskdata"
	"github.com/expandorg/requester-service/pkg/thumbnailupload"
	"github.com/expandorg/requester-service/pkg/upload"
)

func main() {
	environment, err := env.Detect()
	if err != nil {
		log.Fatalln(err)
	}

	// Connect to db
	session, err := db.Connect(environment)
	if err != nil {
		log.Fatalln(err)
	}
	defer session.Close()

	sf := mongo.NewSessionFactory(session)
	// Initialize db with test data if env is compose
	if environment == env.Compose {
		templatemocks.Populate(sf.GetDb(session))
	}

	app := newApp(sf, environment)
	app.Run(":8081")
}

func getStorage(e env.Env) upload.Storage {
	if e == env.Compose || e == env.Local {
		return new(upload.FileSystem)
	}
	return new(upload.GCloud)
}

var commit = ""

func newApp(sf *mongo.SessionFactory, e env.Env) *api.App {
	// repositories
	drafts := draftservice.NewRepository()
	data := taskdata.NewRepository()
	templates := drafttemplates.NewRepository()
	onboarding := onboardingtemplate.NewRepository()

	// stores
	store := getStorage(e)

	// services
	draftService := draftservice.NewService(drafts, data, templates)
	dashboardService := dashboard.NewService(drafts)
	dataService := taskdata.NewService(data)
	templateService := drafttemplates.NewService(templates)
	onboardingService := onboardingtemplate.NewService(onboarding)

	jobPublish := publisher.New(drafts, data)
	dataUpload := dataupload.New(drafts, data, store, new(upload.CSVImporter))
	thumbUpload := thumbnailupload.New(store, new(upload.ImageImporter))

	theApp := api.NewApp(commit)
	app.Setup(theApp, sf, draftService, dashboardService, dataService, templateService, onboardingService, jobPublish, dataUpload, thumbUpload)
	return theApp
}
