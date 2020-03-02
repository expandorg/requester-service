package publisher

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/gemsorg/svc-kit/svc"
	"github.com/globalsign/mgo/bson"

	"github.com/expandorg/requester-service/pkg/app"
	"github.com/expandorg/requester-service/pkg/backend"
	m "github.com/expandorg/requester-service/pkg/model"
	"github.com/expandorg/requester-service/pkg/nulls"
)

// Service interface.
type Service interface {
	Prepublish(ctx *app.ReqCtx, draft *m.Draft, token string) (*m.Draft, error)
	AdminPublish(ctx *app.ReqCtx, draftID bson.ObjectId, token string) (*m.Draft, error)
	AdminReject(ctx *app.ReqCtx, draftID bson.ObjectId, message string, token string) error
}

type publishService struct {
	drafts m.DraftRepository
	data   m.DataRepository
}

// New returns a new instance of a draft publisher.
func New(drafts m.DraftRepository, data m.DataRepository) Service {
	return &publishService{
		drafts: drafts,
		data:   data,
	}
}

func (ps *publishService) Prepublish(ctx *app.ReqCtx, draft *m.Draft, token string) (*m.Draft, error) {
	tasks, err := ps.getJobTasks(ctx, draft)
	if err != nil {
		return nil, err
	}

	result, err := ps.sendPreblishRequest(token, newBackendRequest(draft, tasks))
	if err != nil {
		return nil, err
	}

	// Update Draft Status to Pending
	if draft.Status == m.DraftStatus {
		draft.Status = m.PendingStatus
		draft.JobID = result.Job.ID
		err = ps.drafts.Update(ctx, draft)
		if err != nil {
			return nil, err
		}
	}
	return draft, nil
}

func (ps *publishService) getJobTasks(ctx *app.ReqCtx, draft *m.Draft) ([]*taskReq, error) {
	if len(draft.DataID) > 0 {
		data, err := ps.data.Find(ctx, draft.ID, draft.DataID)
		if err != nil {
			return nil, err
		}
		return makeRequestTasks(draft, data)
	}
	// if no variables declare a blank task
	if !draft.TaskForm.HasVariables() {
		return []*taskReq{&taskReq{IsActive: true, TaskData: json.RawMessage("null")}}, nil
	}
	return nil, nil
}

// sendPreblishRequest backend request
func (ps *publishService) sendPreblishRequest(token string, draft *publishBERequest) (*prepublishResponse, error) {
	reqBodyBytes, err := json.Marshal(draft)
	if err != nil {
		return nil, err
	}
	reqBody := bytes.NewBuffer(reqBodyBytes)

	res, err := backend.Request(http.MethodPost, "/prepublish/job", token, reqBody)
	if err != nil {
		return nil, err
	}

	resp := new(prepublishResponse)
	err = json.NewDecoder(res.Body).Decode(resp)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusCreated && res.StatusCode != http.StatusOK {
		if resp.Error != "" {
			return nil, &svc.ServiceError{Status: res.StatusCode, Err: errors.New(resp.Error)}
		}
		return nil, svc.ApplicationError(fmt.Errorf("Request to publish draft returned status %d", res.StatusCode))
	}
	return resp, nil
}

func findVariableBinding(vs []*m.TaskDataColumn, variable string) *variableBinding {
	for i, v := range vs {
		if !v.Skipped && v.Variable == variable {
			return &variableBinding{Variable: variable, DataIndex: i}
		}
	}
	return nil
}

func bindVariables(variables []string, columns []*m.TaskDataColumn) []*variableBinding {
	bindings := make([]*variableBinding, 0)
	numVars := len(variables)
	for i := 0; i < numVars; i++ {
		binding := findVariableBinding(columns, variables[i])
		if binding != nil {
			bindings = append(bindings, binding)
		}
	}
	return bindings
}

func makeRequestTasks(draft *m.Draft, taskData *m.TaskData) ([]*taskReq, error) {
	numRows := len(taskData.Values)
	bidnings := bindVariables(draft.Variables, taskData.Columns)
	values := make([]*taskReq, numRows)

	for i := 0; i < numRows; i++ {
		val, err := makeRequestTask(bidnings, taskData.Values[i])
		if err != nil {
			return nil, err
		}
		values[i] = val
	}
	return values, nil
}

func makeRequestTask(bindings []*variableBinding, values []string) (*taskReq, error) {
	l := len(bindings)
	m := make(map[string]string)
	for j := 0; j < l; j++ {
		m[bindings[j].Variable] = values[bindings[j].DataIndex]
	}
	data, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}
	return &taskReq{IsActive: true, TaskData: data}, nil
}

func (ps *publishService) AdminPublish(ctx *app.ReqCtx, draftID bson.ObjectId, token string) (*m.Draft, error) {
	draft, err := ps.drafts.FindByID(ctx, draftID)
	if err != nil {
		return nil, svc.ArgumentsErr(errors.New("Draft not found"))
	}
	if draft.Status != m.PendingStatus {
		return nil, svc.ArgumentsErr(errors.New("Draft should be in pedning state"))
	}

	_, err = ps.sendPublishRequest(token, draft)
	if err != nil {
		return nil, err
	}
	draft.Status = m.InProgressStatus
	err = ps.drafts.Update(ctx, draft)
	if err != nil {
		return nil, err
	}
	return draft, nil
}

func (ps *publishService) AdminReject(ctx *app.ReqCtx, draftID bson.ObjectId, message string, token string) error {
	draft, err := ps.drafts.FindByID(ctx, draftID)
	if err != nil {
		return svc.ArgumentsErr(errors.New("Draft not found"))
	}
	if draft.Status != m.PendingStatus {
		return svc.ArgumentsErr(errors.New("Draft should be in pedning state"))
	}

	body := RejectDraftRequest{
		DraftID:     draft.ID,
		DraftName:   draft.Name,
		RequesterID: draft.RequesterID,
		Message:     message,
	}

	err = ps.sendNorifyRequesterRequest(token, &body)
	if err != nil {
		return err
	}
	return nil
}

func (ps *publishService) sendPublishRequest(token string, draft *m.Draft) (*prepublishResponse, error) {
	reqBodyBytes, err := json.Marshal(draft)
	if err != nil {
		return nil, err
	}
	reqBody := bytes.NewBuffer(reqBodyBytes)

	res, err := backend.AdminRequest(http.MethodPost, "/admin/publish/job", token, reqBody)
	if err != nil {
		return nil, err
	}

	resp := new(prepublishResponse)
	err = json.NewDecoder(res.Body).Decode(resp)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusCreated && res.StatusCode != http.StatusOK {
		if resp.Error != "" {
			return nil, &svc.ServiceError{Status: res.StatusCode, Err: errors.New(resp.Error)}
		}
		return nil, svc.ApplicationError(fmt.Errorf("Request to publish draft returned status %d", res.StatusCode))
	}
	return resp, nil
}

func (ps *publishService) sendNorifyRequesterRequest(token string, body *RejectDraftRequest) error {
	reqBodyBytes, err := json.Marshal(body)
	if err != nil {
		return err
	}
	reqBody := bytes.NewBuffer(reqBodyBytes)
	res, err := backend.AdminRequest(http.MethodPost, "/admin/requester/notify", token, reqBody)
	if err != nil {
		return err
	}
	if res.StatusCode != http.StatusOK {
		return svc.ApplicationError(fmt.Errorf("Request to reject draft returned status %d", res.StatusCode))
	}
	return nil
}

type publishBERequest struct {
	m.Draft
	IsActive bool       `json:"isActive"`
	Tasks    []*taskReq `json:"tasks"`
}

type taskReq struct {
	IsActive bool            `json:"isActive"`
	TaskData json.RawMessage `json:"taskData"`
}

type variableBinding struct {
	Variable  string
	DataIndex int
}

type prepublishResponse struct {
	Job   Job    `json:"job"`
	Error string `json:"error"`
}

func newBackendRequest(draft *m.Draft, tasks []*taskReq) *publishBERequest {
	job := &publishBERequest{
		Draft: m.Draft{
			RequesterID:      draft.RequesterID,
			JobID:            draft.JobID,
			Status:           draft.Status,
			Name:             draft.Name,
			Description:      draft.Description,
			Logo:             draft.Logo,
			TaskForm:         draft.TaskForm,
			VerificationForm: draft.VerificationForm,
			Onboarding:       draft.Onboarding,
			Eligibility:      draft.Eligibility,
			Assignment:       draft.Assignment,
			Verification:     draft.Verification,
			Funding:          draft.Funding,
			CallbackURL:      draft.CallbackURL,
			StartDate:        draft.StartDate,
			EndDate:          draft.EndDate,
			Whitelist:        draft.Whitelist,
		},
		IsActive: false,
	}
	if tasks != nil {
		job.Tasks = tasks
	}
	return job
}

type RejectDraftRequest struct {
	DraftID     bson.ObjectId `json:"draftId,omitempty"`
	DraftName   string        `json:"draftName,omitempty"`
	RequesterID uint64        `json:"requesterId,omitempty"`
	Message     string        `json:"message,omitempty"`
}

type Job struct {
	ID              uint64       `json:"id"`
	Name            string       `json:"name"`
	Description     nulls.String `json:"description"`
	Logo            nulls.String `json:"logo"`
	RequesterID     uint64       `json:"requesterId"`
	IsActive        bool         `json:"isActive"`
	RequiresProfile bool         `json:"requiresProfile,omitempty"`
}
