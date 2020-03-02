package ctx

import (
	"fmt"
	"mime/multipart"
	"strconv"

	"github.com/expandorg/requester-service/pkg/app"

	"github.com/expandorg/requester-service/pkg/http/auth"
	m "github.com/expandorg/requester-service/pkg/model"
	"github.com/gin-gonic/gin"
	"github.com/globalsign/mgo/bson"
)

const (
	DraftKey    = "draft"
	TaskDataKey = "task-data"
)

func UserID(c *gin.Context) uint64 {
	userID, _ := auth.GetUserID(c)
	return userID
}

func DraftID(c *gin.Context) bson.ObjectId {
	return bson.ObjectIdHex(c.Param("draftID"))
}

func Draft(c *gin.Context) *m.Draft {
	draft, _ := c.Get(DraftKey)
	return draft.(*m.Draft)
}

func TemplateID(c *gin.Context) bson.ObjectId {
	return bson.ObjectIdHex(c.Param("templateID"))
}

func TaskStatus(c *gin.Context) (m.StatusType, error) {
	status := m.StatusType(c.Param("taskStatus"))
	if !status.IsValid() {
		return "", fmt.Errorf("%s is not a valid task status", status)
	}
	return status, nil
}

func DataID(c *gin.Context) bson.ObjectId {
	return bson.ObjectIdHex(c.Param("dataID"))
}

func TaskData(c *gin.Context) *m.TaskData {
	data, _ := c.Get(TaskDataKey)
	return data.(*m.TaskData)
}

func JWT(c *gin.Context) string {
	cookie, _ := auth.GetJWTCookie(c)
	return cookie
}

func Page(c *gin.Context) uint64 {
	pageString := c.DefaultQuery("page", "0")
	pageUint, err := strconv.ParseInt(pageString, 10, 64)
	if err != nil {
		return 0
	}
	return uint64(pageUint)
}

func ReqCtx(c *gin.Context) *app.ReqCtx {
	ct, _ := c.Get("reqctx")
	return ct.(*app.ReqCtx)
}

func FileHeader(c *gin.Context, name string) (*multipart.FileHeader, error) {
	fileHeader, err := c.FormFile(name)
	if err != nil {
		return nil, err
	}
	return fileHeader, nil
}
