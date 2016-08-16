package main

import (
	"github.com/gin-gonic/gin"
	"github.com/olahol/melody"
	"gopkg.in/mgo.v2/bson"
	"log"
	"net/http"
)

func GetUser(c *gin.Context) {
	var user User
	username := c.Params.ByName("id")
	coll := GetColl("mongodb://localhost", "chat", "users")

	err := coll.Find(bson.M{"username": username}).One(&user)
	if err == nil {
		user := &User{
			Username: user.Username,
			Pass:     user.Pass,
		}
		c.JSON(200, user)
	} else {
		c.JSON(404, gin.H{"error": "user not found"})
	}
}

func PostUser(c *gin.Context) {
	var user User
	c.Bind(&user)

	log.Printf("%+v", user)

	coll := GetColl("mongodb://localhost", "chat", "users")
	if user.Username != "" && user.Pass != "" {
		err := coll.Insert(&user)
		if err == nil {
			user := &User{
				Username: user.Username,
				Pass:     user.Pass,
			}
			c.JSON(201, user)
		} else {
			log.Fatal(err)
		}
	} else {
		c.JSON(422, gin.H{"error": "fields are empty"})
	}
}

func main() {
	r := gin.Default()
	m := melody.New()

	v1 := r.Group("api/v1")
	{
		v1.GET("/users/:id", GetUser)
		v1.POST("/users", PostUser)
	}

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
