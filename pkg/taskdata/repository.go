package taskdata

import (
	"math"

	"github.com/globalsign/mgo/bson"

	"github.com/expandorg/requester-service/pkg/app"
	m "github.com/expandorg/requester-service/pkg/model"
	"github.com/expandorg/requester-service/pkg/upload"
	"github.com/gemsorg/svc-kit/mongo"
)

const (
	rowPageSize = 12
)

type dataRepository struct {
	mongo.Repository
}

// NewRepository returns a new instance of a ata repository.
func NewRepository() m.DataRepository {
	return &dataRepository{
		Repository: mongo.Repository{Name: "taskData"},
	}
}

func (r *dataRepository) Find(ctx *app.ReqCtx, draftID bson.ObjectId, dataID bson.ObjectId) (*m.TaskData, error) {
	selector := bson.M{"_id": dataID, "draftId": draftID}
	data := new(m.TaskData)
	err := r.Col(ctx).Find(selector).One(data)
	return data, err
}

func (r *dataRepository) FindPaged(ctx *app.ReqCtx, draftID bson.ObjectId, dataID bson.ObjectId, page uint64) (*m.TaskData, uint64, error) {
	fullData, err := r.Find(ctx, draftID, dataID)

	if err != nil {
		return nil, 0, err
	}

	// Determines how many rows in the dataset
	numRows := len(fullData.Values)
	// How many columns to be displayed
	numColumns := len(fullData.Columns)
	// How many pages will be displayed
	numPages := uint64(math.Ceil(float64(numRows) / rowPageSize))
	// this calculates the first id of the daya
	first := page * rowPageSize
	// calculates the last column header
	lastColumn := numColumns

	data := new(m.TaskData)
	data.ID = fullData.ID
	data.DraftID = fullData.DraftID
	data.Columns = fullData.Columns[0:lastColumn]

	if numRows == 0 {
		data.Values = [][]string{}
		return data, numPages, nil
	}

	// find number of rows on page because the last page might now be full
	numRowsResponded := rowPageSize
	if page == numPages-1 && numRows%rowPageSize != 0 {
		numRowsResponded = (numRows % rowPageSize)
	}
	values := make([][]string, numRowsResponded)
	// Get data until we've hit the limit for the page
	for i := 0; i < numRowsResponded; i++ {
		values[i] = fullData.Values[uint64(i)+first][0:numColumns]
	}
	data.Values = values

	return data, numPages, nil
}

func (r *dataRepository) Copy(ctx *app.ReqCtx, data *m.TaskData, draftID bson.ObjectId) (*m.TaskData, error) {
	data.ID = bson.NewObjectId()
	data.DraftID = draftID
	err := r.Col(ctx).Insert(data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (r *dataRepository) Create(ctx *app.ReqCtx, draftID bson.ObjectId, d *upload.Data) (bson.ObjectId, error) {
	numColumns := d.Size()
	columns := make([]*m.TaskDataColumn, numColumns)
	for i := 0; i < numColumns; i++ {
		columns[i] = &m.TaskDataColumn{
			Name:     d.Names[i],
			Variable: "",
			Skipped:  false,
		}
	}

	dataID := bson.NewObjectId()
	taskData := &m.TaskData{
		ID:      dataID,
		DraftID: draftID,
		Columns: columns,
		Values:  d.Data,
	}
	err := r.Col(ctx).Insert(taskData)
	return dataID, err
}

func (r *dataRepository) UpdateColumns(ctx *app.ReqCtx, taskData *m.TaskData, columns []*m.TaskDataColumn) error {
	taskData.Columns = columns
	return r.Col(ctx).Update(bson.M{"_id": taskData.ID, "draftId": taskData.DraftID}, taskData)
}

func (r *dataRepository) Delete(ctx *app.ReqCtx, ID bson.ObjectId) error {
	return r.Col(ctx).RemoveId(ID)
}
