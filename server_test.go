package chat_test

import (
	jwt "github.com/dgrijalva/jwt-go"
	chat "github.com/ttymed/chat-server"
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

	if reflect.TypeOf(testToken) != reflect.TypeOf(actualToken) {
		t.Errorf("Test failed, CreateToken not creating token")
	}
}

func TestGetUser(t *testing.T) {
	var expected chat.User
	actual := chat.GetUser("test", "test", "test")

	if reflect.TypeOf(expected) != reflect.TypeOf(actual) {
		t.Errorf("Test failed, GetUser not returning user")
	}
}
