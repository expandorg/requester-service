package mongo

import (
	"os"

	"github.com/globalsign/mgo"
)

type SessionFactory struct {
	mesterSession *mgo.Session
	dbName        string
}

func NewSessionFactory(sess *mgo.Session) *SessionFactory {
	return &SessionFactory{mesterSession: sess, dbName: os.Getenv("DB_NAME")}
}

func (sf *SessionFactory) Get() *mgo.Session {
	return sf.mesterSession.Clone()
}

func (sf *SessionFactory) GetDb(sess *mgo.Session) *mgo.Database {
	return sess.DB(sf.dbName)
}
