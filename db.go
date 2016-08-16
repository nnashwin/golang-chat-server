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

// func GetUser(dbName string, collName string, userName string) User {
// 	coll := getColl("mongodb://localhost", dbName, collName)

// 	result := User{}
// 	err := c.Find(bson.M{"id": userName}).One(&result)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	return result
// }
