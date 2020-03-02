package logerror

import (
	"net/http"
	"strconv"

	bugsnag "github.com/bugsnag/bugsnag-go"
	"github.com/expandorg/requester-service/pkg/http/auth"
	"github.com/expandorg/requester-service/pkg/logger"
	"github.com/gemsorg/svc-kit/cfg/env"
	"github.com/gemsorg/svc-kit/http/api"
	"github.com/gin-gonic/gin"
)

// Log error
func Log(c *gin.Context, e *api.Error) {
	l := &logger.ServerLogger{
		Context:      c,
		RequestBody:  e.Body,
		ResponseCode: e.Status,
	}
	l.LogError(e.Err)

	if env.Get() == env.Production {
		if e.Private == true {
			bugsnagError(c, e.Body, e.Status, e.Err)
		} else {
			if e.Status != http.StatusForbidden {
				bugsnagError(c, e.Body, e.Status, e.Err)
			}
		}
	}
}

// bs error
func bugsnagError(c *gin.Context, reqBody interface{}, status int, reportedErr error) {
	r := c.Copy().Request
	meta := bugsnag.MetaData{
		"meta": {
			"responseCode": status,
		},
	}
	if reqBody != nil {
		meta["meta"]["requestBody"] = reqBody
	}
	userID, err := auth.GetUserID(c)
	if err == nil {
		user := bugsnag.User{Id: strconv.FormatUint(userID, 10)}
		bugsnag.Notify(reportedErr, r, user, meta)
	} else {
		bugsnag.Notify(reportedErr, r, meta)
	}
}
