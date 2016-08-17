package main

import (
	//	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/olahol/melody"
	"gopkg.in/mgo.v2/bson"
	"net/http"
)

func GetUser(con *gin.Context) {
	var user User
	username := con.Params.ByName("username")
	coll := GetColl("mongodb://localhost", "chat", "users")

	err := coll.Find(bson.M{"username": username}).One(&user)
	if err != nil {
		con.JSON(404, "user not found")
	} else {
		response := &User{
			Username: user.Username,
			Pass:     user.Pass,
		}
		con.JSON(200, response)
	}
}

func PostUser(con *gin.Context) {
	var user User
	con.Bind(&user)

	coll := GetColl("mongodb://localhost", "chat", "users")
	err := coll.Find(bson.M{"username": user.Username}).One(&user)
	if err == nil {
		con.JSON(409, gin.H{"error": "That username already exists"})
	} else {

		if user.Username != "" && user.Pass != "" {
			err := coll.Insert(&user)
			if err == nil {
				user := &User{
					Username: user.Username,
					Pass:     user.Pass,
				}

				con.JSON(201, user)

			} else {
				con.JSON(500, gin.H{"error": "server insertion error"})
			}

		} else {
			con.JSON(422, gin.H{"error": "fields are empty"})
		}

	}
}

func main() {
	r := gin.New()

	r.Static("/js", "./js")
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	m := melody.New()

	v1 := r.Group("api/v1")
	{
		v1.GET("/users/:id", GetUser)
		v1.POST("/users", PostUser)
	}

	r.GET("/", func(c *gin.Context) {
		http.ServeFile(c.Writer, c.Request, "index.html")
	})

	r.GET("/signup", func(c *gin.Context) {
		http.ServeFile(c.Writer, c.Request, "signup.html")
	})

	r.GET("/ws", func(c *gin.Context) {
		m.HandleRequest(c.Writer, c.Request)
	})

	m.HandleMessage(func(s *melody.Session, msg []byte) {
		m.Broadcast(msg)
	})

	r.Run(":5000")
}
