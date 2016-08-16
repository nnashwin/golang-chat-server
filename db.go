package main

import (
	"gopkg.in/mgo.v2"
	"log"
)

func NewSession(dbAddr string) *mgo.Session {
	session, err := mgo.Dial(dbAddr)
	if err != nil {
		log.Fatal(err)
	}

	return session
}

func GetColl(dbAddr string, dbName string, collName string) *mgo.Collection {
	session := NewSession(dbAddr)

	c := session.DB(dbName).C(collName)

	return c
}
