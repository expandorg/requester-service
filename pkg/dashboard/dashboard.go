package dashboard

import (
	"github.com/expandorg/requester-service/pkg/app"
	m "github.com/expandorg/requester-service/pkg/model"
	"github.com/globalsign/mgo/bson"
)

// Service interface.
type Service interface {
	ListByUserID(ctx *app.ReqCtx, userID uint64) ([]*DashboardTask, error)
	ListByUserIDAndStatus(ctx *app.ReqCtx, userID uint64, status m.StatusType) ([]*DashboardTask, error)
	GetAdminDrafts(ctx *app.ReqCtx) ([]*PendingDraft, error)
}

type service struct {
	drafts m.DraftRepository
}

// New returns a new instance of a tasks service
func NewService(drafts m.DraftRepository) Service {
	return &service{drafts: drafts}
}

func (s *service) ListByUserID(ctx *app.ReqCtx, userID uint64) ([]*DashboardTask, error) {
	drafts, err := s.drafts.ListByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	return toTaskList(drafts), nil
}

func (s *service) ListByUserIDAndStatus(ctx *app.ReqCtx, userID uint64, status m.StatusType) ([]*DashboardTask, error) {
	drafts, err := s.drafts.ListByUserIDAndStatus(ctx, userID, status)
	if err != nil {
		return nil, err
	}
	return toTaskList(drafts), nil
}

func (s *service) GetAdminDrafts(ctx *app.ReqCtx) ([]*PendingDraft, error) {
	drafts, err := s.drafts.List(ctx)
	if err != nil {
		return nil, err
	}
	return toPendingDraftList(drafts), nil
}

type DashboardTask struct {
	ID     bson.ObjectId `json:"id"`
	Status m.StatusType  `json:"status"`
	Name   string        `json:"name,omitempty"`
	Logo   string        `json:"logo,omitempty"`
	JobID  uint64        `json:"jobId,omitempty"`
}

func taskFromDraft(draft *m.Draft) *DashboardTask {
	return &DashboardTask{
		ID:     draft.ID,
		Status: draft.Status,
		Name:   draft.Name,
		Logo:   draft.Logo,
		JobID:  draft.JobID,
	}
}

func toTaskList(drafts []*m.Draft) []*DashboardTask {
	n := len(drafts)
	list := make([]*DashboardTask, n)
	for i := 0; i < n; i++ {
		list[i] = taskFromDraft(drafts[i])
	}
	return list
}

type PendingDraft struct {
	ID          bson.ObjectId `json:"id"`
	Status      m.StatusType  `json:"status"`
	RequesterID uint64        `json:"requesterId"`
	Name        string        `json:"name,omitempty"`
	Logo        string        `json:"logo,omitempty"`
	JobID       uint64        `json:"jobId,omitempty"`
}

func toPendingDraftList(drafts []*m.Draft) []*PendingDraft {
	n := len(drafts)
	list := make([]*PendingDraft, n)
	for i := 0; i < n; i++ {
		list[i] = &PendingDraft{
			ID:          drafts[i].ID,
			Status:      drafts[i].Status,
			Name:        drafts[i].Name,
			Logo:        drafts[i].Logo,
			JobID:       drafts[i].JobID,
			RequesterID: drafts[i].RequesterID,
		}
	}
	return list
}
