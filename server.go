package main

import (
	"errors"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/olahol/melody"
	db "github.com/ttymed/mwrap"
	"gopkg.in/mgo.v2/bson"
	"log"
	"net/http"
	"time"
)

var mySigningKey = []byte("secret")

func CreateToken() string {
	type CustomClaims struct {
		Authorized bool `json:"auth"`
		jwt.StandardClaims
	}

	/* set token claims */
	claims := CustomClaims{
		true,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute * 5).Unix(),
			Issuer:    "test",
		},
	}

	/* create token */
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	/* sign the token with a secret */
	signedToken, _ := token.SignedString(mySigningKey)

	return signedToken
}

func ParseToken(tokenStr string) bool {
	type CustomClaims struct {
		Authorized bool `json:"auth"`
		jwt.StandardClaims
	}

	token, err := jwt.ParseWithClaims(tokenStr, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return mySigningKey, nil
	})

	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		log.Printf("%v %v", claims.Authorized, claims.StandardClaims.ExpiresAt)
		return true
	} else {
		log.Println(err)
		return false
	}
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

func HandleSignup(con *gin.Context) {
	var user User
	con.Bind(&user)

	_, err := GetUser(user.Username, "chat", "users")
	if err == nil {
		con.JSON(409, gin.H{"error": "That username already exists"})
	} else {
		if user.Username != "" && user.Pass != "" {
			err := CreateUser(user, "chat", "users")
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

func GetUser(username string, dbName string, collName string) (User, error) {
	var user User
	coll, _ := db.GetColl("mongodb://localhost", dbName, collName)

	err := coll.Find(bson.M{"username": username}).One(&user)
	if err != nil {
		return User{}, errors.New("The User Does Not Exist")
	}
	return user, nil
}

func CreateUser(user User, dbName string, collName string) error {
	coll, _ := db.GetColl("mongodb://localhost", dbName, collName)
	err := coll.Insert(&user)
	if err != nil {
		return errors.New("the user could not be created.  Check CreateUser function.")
	}
	return nil
}

func LoginUser(con *gin.Context) {
	var userReq User
	con.Bind(&userReq)
	userInfo, err := GetUser(userReq.Username, "chat", "users")
	passMatch := (userInfo.Pass == userReq.Pass)
	if err == nil && passMatch == true {
		token := CreateToken()
		con.JSON(200, token)
	} else if err == nil && passMatch == false {
		con.JSON(400, gin.H{"error": "That Password is incorrect"})
	} else if err != nil {
		con.JSON(400, gin.H{"error": "You could not be logged in.  Please check your credentials and try again."})
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
		v1.POST("/users/signup", HandleSignup)
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
