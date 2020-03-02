package app

import (
	"github.com/globalsign/mgo"
)

type ReqCtx struct {
	db *mgo.Database
}

func NewReqCtx(db *mgo.Database) *ReqCtx {
	return &ReqCtx{db: db}
}

func (c *ReqCtx) Db() *mgo.Database {
	return c.db
}
