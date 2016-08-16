package main

import (
	"github.com/gin-gonic/gin"
	"github.com/olahol/melody"
	"gopkg.in/mgo.v2/bson"
	"log"
	"net/http"
)

func GetUser(c *gin.Context) {
	username := c.Params.ByName("username")
	coll := GetColl("mongodb://localhost", "chat", "users")

	result := User{}
	err := coll.Find(bson.M{"id": username}).One(&result)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	r := gin.Default()
	m := melody.New()

	v1 := r.Group("api/v1")
	{
		v1.GET("/users/:id", GetUser)
		v1.POST("/users", PostUser)
		v1.PUT("/users/:id", UpdateUser)
		v1.DELETE("/users/:id", DeleteUser)
	}

	r.POST("/users", func(c *gin.Context) {
		coll := GetColl("mongodb://localhost", "chat", "users")

	})
	r.GET("/", func(c *gin.Context) {
		http.ServeFile(c.Writer, c.Request, "index.html")
	})

	r.GET("/ws", func(c *gin.Context) {
		m.HandleRequest(c.Writer, c.Request)
	})

	m.HandleMessage(func(s *melody.Session, msg []byte) {
		m.Broadcast(msg)
	})

	r.Run(":5000")
}
