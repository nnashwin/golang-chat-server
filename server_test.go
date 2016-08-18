package chat_test

import (
	jwt "github.com/dgrijalva/jwt-go"
	chat "github.com/ttymed/chat-server"
	"log"
	"reflect"
	"testing"
)

func TestCreateToken(t *testing.T) {
	type CustomClaims struct {
		Authorized bool `json:"auth"`
		jwt.StandardClaims
	}

	claims := CustomClaims{
		true,
		jwt.StandardClaims{
			ExpiresAt: 15000,
			Issuer:    "test",
		},
	}
	actualToken := chat.CreateToken()
	testToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	log.Println(reflect.TypeOf(testToken) == reflect.TypeOf(actualToken))
}
