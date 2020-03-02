package ctx

import (
	"github.com/expandorg/requester-service/pkg/draftservice"
	"github.com/expandorg/requester-service/pkg/publisher"
	"github.com/expandorg/requester-service/pkg/taskdata"
	"github.com/gin-gonic/gin"
)

func BindDataColumnsReq(c *gin.Context) (*taskdata.ColumnsRequest, error) {
	var body taskdata.ColumnsRequest
	err := c.BindJSON(&body)
	if err != nil {
		return nil, err
	}
	return &body, nil
}

func BindCreateDraftRequest(c *gin.Context) (*draftservice.CreateRequest, error) {
	var body draftservice.CreateRequest
	err := c.BindJSON(&body)
	if err != nil {
		return nil, err
	}
	return &body, nil
}

func BindUpdateDraftSettingsRequest(c *gin.Context) (*draftservice.UpdateSettingsRequest, error) {
	var body draftservice.UpdateSettingsRequest
	err := c.BindJSON(&body)
	if err != nil {
		return nil, err
	}
	return &body, nil
}

func BindUpdateDraftVariablesRequest(c *gin.Context) (*draftservice.UpdateVariablesRequest, error) {
	var body draftservice.UpdateVariablesRequest
	err := c.BindJSON(&body)
	if err != nil {
		return nil, err
	}
	return &body, nil
}

func BindUpdateDraftVerificationRequest(c *gin.Context) (*draftservice.UpdateVerificationRequest, error) {
	var body draftservice.UpdateVerificationRequest
	err := c.BindJSON(&body)
	if err != nil {
		return nil, err
	}
	return &body, nil
}

func BindUpdateDraftRequest(c *gin.Context) (*draftservice.UpdateRequest, error) {
	var body draftservice.UpdateRequest
	err := c.BindJSON(&body)
	if err != nil {
		return nil, err
	}
	return &body, nil
}

func BindRejectDraftRequest(c *gin.Context) (*publisher.RejectDraftRequest, error) {
	var body publisher.RejectDraftRequest
	err := c.BindJSON(&body)
	if err != nil {
		return nil, err
	}
	return &body, nil
}
