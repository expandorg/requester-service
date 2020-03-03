package db

import (
	"crypto/tls"
	"fmt"
	"net"
	"os"
	"time"

	"github.com/expandorg/requester-service/pkg/svc-kit/cfg/env"
	"github.com/globalsign/mgo"
)

// Connect returns database & session.
func Connect(e env.Env) (*mgo.Session, error) {
	session, err := getSession(e, os.Getenv("DB_HOST"))

	if err != nil {
		return nil, err
	}
	awaitConnect(session)
	session.SetMode(mgo.Monotonic, true)
	return session, nil
}

func getSession(e env.Env, host string) (*mgo.Session, error) {
	if e == env.Compose || e == env.Local {
		return dialLocal(host)
	}
	return dialRemote(host)
}

func dialLocal(host string) (*mgo.Session, error) {
	var credentials string
	cs := fmt.Sprintf("mongodb://%s%s:27017", credentials, host)
	return mgo.Dial(cs)
}

func dialRemote(host string) (*mgo.Session, error) {
	tlsConfig := &tls.Config{}

	dialInfo := &mgo.DialInfo{
		Addrs:    []string{host},
		Database: os.Getenv("AUTH_DB"),
		Username: os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
	}
	dialInfo.DialServer = func(addr *mgo.ServerAddr) (net.Conn, error) {
		conn, err := tls.Dial("tcp", addr.String(), tlsConfig)
		return conn, err
	}
	return mgo.DialWithInfo(dialInfo)
}

func awaitConnect(session *mgo.Session) {
	for {
		err := session.Ping()
		if err == nil {
			break
		}
		fmt.Println("Retrying connection:", err)
		time.Sleep(time.Second)
	}
}
