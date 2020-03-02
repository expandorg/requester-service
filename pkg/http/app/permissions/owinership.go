package permissions

import (
	"errors"
	"os"
	"strconv"

	"github.com/expandorg/requester-service/pkg/http/auth"

	"github.com/expandorg/requester-service/pkg/draftservice"
	"github.com/expandorg/requester-service/pkg/http/app/ctx"
	"github.com/expandorg/requester-service/pkg/taskdata"
	"github.com/gemsorg/svc-kit/http/api"
	"github.com/gin-gonic/gin"
)

// SetOwnDraft set own draft middleware
func SetOwnDraft(h *api.Handler, drafts draftservice.Service) gin.HandlerFunc {
	return h.Ensure(func(c *gin.Context) error {
		draft, err := drafts.GetDraft(ctx.ReqCtx(c), ctx.DraftID(c))
		if err != nil {
			return api.NotFound(err)
		}

		if draft.RequesterID != ctx.UserID(c) {
			return api.Forbidden(NotRequesterErr{})
		}

		c.Set(ctx.DraftKey, draft)
		return nil
	})
}

// SetOwnData set own data middleware
func SetOwnData(h *api.Handler, data taskdata.Service) gin.HandlerFunc {
	return h.Ensure(func(c *gin.Context) error {
		draft := ctx.Draft(c)
		dataID := ctx.DataID(c)

		if draft.DataID != dataID {
			return api.BadRequest(errors.New("Data id is not assigned to draft"))
		}

		data, err := data.Find(ctx.ReqCtx(c), draft.ID, dataID)
		if err != nil {
			return api.NotFound(err)
		}
		c.Set(ctx.TaskDataKey, data)
		return nil
	})
}

func IsAdmin(h *api.Handler) gin.HandlerFunc {
	return h.Ensure(func(c *gin.Context) error {
		userID, _ := auth.GetRequestUserID(c)
		if userID != 0 {
			if !IsAdminUser(userID) {
				return api.Forbidden(errors.New("Current user is not an admin"))
			}
			return nil
		}
		apiKey := c.GetHeader("Authorization")
		if apiKey != os.Getenv("ADMIN_API_KEY") {
			return api.Forbidden(errors.New("Current user is not an admin"))
		}
		return nil
	})
}

// Check weather user is admin or not
func IsAdminUser(userID uint64) bool {
	moderatorID, err := strconv.ParseUint(os.Getenv("MODERATOR_ID"), 10, 64)
	if err != nil {
		return false
	}
	return userID == moderatorID
}

type NotRequesterErr struct{}

func (err NotRequesterErr) Error() string {
	return "Current user is not the given job's owner"
}
