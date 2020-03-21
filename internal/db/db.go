package db

import (
	"errors"

	"github.com/globalsign/mgo"
)

func MustConnect(uri, username, password string) *mgo.Session {
	session, err := Connect(uri, username, password)
	if err != nil {
		panic(err)
	}
	return session
}

func Connect(uri, username, password string) (*mgo.Session, error) {
	var err error

	var dialInfo *mgo.DialInfo
	if uri == "" {
		return nil, errors.New("unable to connect with empty uri")
	} else {
		dialInfo, err = mgo.ParseURL(uri)
		if err != nil {
			return nil, err
		}
	}
	dialInfo.Username = username
	dialInfo.Password = password

	session, err := mgo.DialWithInfo(dialInfo)
	if err != nil {
		return nil, err
	}

	err = session.Ping()
	if err != nil {
		return nil, err
	}
	return session, nil
}
