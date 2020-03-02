package app

import (
	"os"
	"strings"

	"github.com/gemsorg/svc-kit/mongo"

	"github.com/gemsorg/svc-kit/cfg/env"
	"github.com/gemsorg/svc-kit/http/api"
	"github.com/gemsorg/svc-kit/http/middleware"

	"github.com/expandorg/requester-service/pkg/dashboard"
	"github.com/expandorg/requester-service/pkg/dataupload"
	"github.com/expandorg/requester-service/pkg/draftservice"
	"github.com/expandorg/requester-service/pkg/drafttemplates"
	"github.com/expandorg/requester-service/pkg/http/app/ctx"
	"github.com/expandorg/requester-service/pkg/http/app/endpoints"
	"github.com/expandorg/requester-service/pkg/http/auth"
	"github.com/expandorg/requester-service/pkg/http/logerror"
	"github.com/expandorg/requester-service/pkg/onboardingtemplate"
	"github.com/expandorg/requester-service/pkg/publisher"
	"github.com/expandorg/requester-service/pkg/taskdata"
	"github.com/expandorg/requester-service/pkg/thumbnailupload"
)

func split(srt string) []string {
	var result = []string{}
	arr := strings.Split(srt, ",")
	for _, s := range arr {
		result = append(result, strings.TrimSpace(s))
	}
	return result
}

// Setup router
func Setup(
	app *api.App,
	sf *mongo.SessionFactory,
	drafts draftservice.Service,
	tasks dashboard.Service,
	data taskdata.Service,
	templates drafttemplates.Service,
	onboarding onboardingtemplate.Service,
	publisher publisher.Service,
	dataUploader dataupload.Service,
	thumbs thumbnailupload.Service,
) {
	e := env.Get()
	if e == env.Production {
		app.Engine.Use(middleware.Bugsnag())
	}

	if e != env.Testing {
		app.Engine.Use(middleware.CORS(split(os.Getenv("REQUESTER_FRONTEND_ADDRESS"))))
	}

	v1 := app.Engine.Group("/v1")
	v1.GET("/version", app.GetVersion)
	v1.Use(ctx.CtxContainerMiddleware(sf))

	private := v1.Group("/")
	handler := api.NewHandler(logerror.Log)
	private.Use(handler.Ensure(auth.Middleware))

	endpoints.AdminRoutes(v1, handler, publisher, tasks, drafts, data, os.Getenv("ADMIN_API_KEY"))

	endpoints.TemplatesRoutes(private, handler, templates)
	endpoints.TasksRoutes(private, handler, tasks)
	endpoints.OnboardingRoutes(private, handler, onboarding)
	endpoints.DraftRoutes(private, handler, drafts, templates, data, publisher)
	endpoints.ThumbnailRoutes(private, handler, thumbs)
	endpoints.DataRoutes(private, handler, drafts, data, dataUploader)
}
