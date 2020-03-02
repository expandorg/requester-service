package endpoints

import (
	"github.com/expandorg/requester-service/pkg/http/app/ctx"
	"github.com/expandorg/requester-service/pkg/thumbnailupload"
	"github.com/gemsorg/svc-kit/http/api"
	"github.com/gin-gonic/gin"
)

// ThumbnailRoutes setup routes
func ThumbnailRoutes(r *gin.RouterGroup, h *api.Handler, t thumbnailupload.Service) {
	c := thumbsController{
		thumbs: t,
	}
	r.POST("/thumbnails/upload", h.Handle(c.ThumbnailUpload))
}

type thumbsController struct {
	thumbs thumbnailupload.Service
}

func (tc *thumbsController) ThumbnailUpload(c *gin.Context) (interface{}, error) {
	fileHeader, err := ctx.FileHeader(c, "thumbnail")
	if err != nil {
		return nil, api.BadRequest(err)
	}
	url, err := tc.thumbs.Upload(ctx.UserID(c), fileHeader)
	if err != nil {
		return nil, err
	}
	return gin.H{"url": url}, nil
}
