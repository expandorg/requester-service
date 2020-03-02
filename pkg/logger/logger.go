package logger

import (
	"encoding/json"
	"fmt"

	"github.com/gin-gonic/gin"
)

type serverErrorLog struct {
	Method      string      `json:"method"`
	URL         string      `json:"url"`
	RequestBody interface{} `json:"requestBody,omitempty"`
	Response    int         `json:"response,omitempty"`
	Error       string      `json:"error"`
}

type Logger interface {
	LogError(error)
}

type ServerLogger struct {
	Context      *gin.Context
	RequestBody  interface{}
	ResponseCode int
}

func (s *ServerLogger) LogError(err error) {
	errLog := serverErrorLog{
		Method:      s.Context.Request.Method,
		URL:         s.Context.Request.URL.String(),
		RequestBody: s.RequestBody,
		Response:    s.ResponseCode,
		Error:       err.Error(),
	}
	data, _ := json.Marshal(errLog)
	fmt.Println(string(data))
}
