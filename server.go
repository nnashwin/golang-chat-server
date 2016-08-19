package chat

import (
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/olahol/melody"
	db "github.com/ttymed/mwrap"
	"gopkg.in/mgo.v2/bson"
	"log"
	"net/http"
	//	"time"
)

var mySigningKey = []byte("secret")

func CreateToken() *jwt.Token {
	type CustomClaims struct {
		Authorized bool `json:"auth"`
		jwt.StandardClaims
	}

	/* set token claims */
	claims := CustomClaims{
		true,
		jwt.StandardClaims{
			ExpiresAt: 15000,
			Issuer:    "test",
		},
	}

	/* create token */
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	/* sign the token with a secret */
	return token
}

// Routes

func HandleGetUser(con *gin.Context) {
	var user User
	username := con.Params.ByName("id")
	coll, _ := db.GetColl("mongodb://localhost", "chat", "users")

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

func GetUser(username string, dbName string, collName string) User {
	var user User
	coll, _ := db.GetColl("mongodb://localhost", dbName, collName)

	err := coll.Find(bson.M{"username": username}).One(&user)
	if err != nil {
		log.Fatal("User not found")
	}
	return user
}

func CreateUser(con *gin.Context) {
	var user User
	con.Bind(&user)

	coll, _ := db.GetColl("mongodb://localhost", "chat", "users")
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

func LoginUser(con *gin.Context) {
	var userReq User
	con.Bind(&userReq)

	userInfo := User{}
	coll, _ := db.GetColl("mongodb://localhost", "chat", "users")
	err := coll.Find(bson.M{"username": userReq.Username}).One(&userInfo)
	log.Printf("%+v", userInfo)
	log.Printf("%+v", userReq)

	passMatch := (userInfo.Pass == userReq.Pass)
	if err == nil && passMatch == true {
		con.JSON(200, userInfo)
	} else {
		con.JSON(400, gin.H{"error": "That Password is incorrect"})
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
		v1.GET("/users/:id", HandleGetUser)
		v1.POST("/users/signup", CreateUser)
		v1.POST("/users/login", LoginUser)
	}

	//r.POST("/get-token", SendToken)

	r.GET("/", func(c *gin.Context) {
		http.ServeFile(c.Writer, c.Request, "index.html")
	})

	r.GET("/login", func(c *gin.Context) {
		http.ServeFile(c.Writer, c.Request, "login.html")
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
