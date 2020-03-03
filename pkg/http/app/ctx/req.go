package ctx

import (
	"github.com/expandorg/requester-service/pkg/app"
	"github.com/expandorg/requester-service/pkg/svc-kit/mongo"
	"github.com/gin-gonic/gin"
)

func CtxContainerMiddleware(sf *mongo.SessionFactory) gin.HandlerFunc {
	return func(c *gin.Context) {
		sess := sf.Get()
		defer sess.Close()
		db := sf.GetDb(sess)
		c.Set("reqctx", app.NewReqCtx(db))
		c.Next()
	}
}
