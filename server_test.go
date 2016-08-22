package main_test

import (
	jwt "github.com/dgrijalva/jwt-go"
	chat "github.com/ttymed/chat-server"
	"reflect"
	"testing"
	"time"
)

func TestCreateToken(t *testing.T) {
	type CustomClaims struct {
		Authorized bool `json:"auth"`
		jwt.StandardClaims
	}

	claims := CustomClaims{
		true,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute * 5).Unix(),
			Issuer:    "test",
		},
	}
	actualTokenString := chat.CreateToken()
	testToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedTest, _ := testToken.SignedString([]byte("secret"))

	if reflect.TypeOf(signedTest) != reflect.TypeOf(actualTokenString) {
		t.Errorf("Test failed, CreateToken not creating token")
	}
}

func TestParseToken(t *testing.T) {
	tokenString := chat.CreateToken()
	isTokenValid := chat.ParseToken(tokenString)
	if isTokenValid != true {
		t.Errorf("Test failed, the token was not valid")
	}
}

func TestGetUser(t *testing.T) {
	var expected chat.User
	actual, _ := chat.GetUser("test", "test", "test")

	if reflect.TypeOf(expected) != reflect.TypeOf(actual) {
		t.Errorf("Test failed, GetUser not returning user")
	}
}

func TestCreateUser(t *testing.T) {
	chat.CreateUser(chat.User{"test1", "test1"}, "test", "test")
	actual, _ := chat.GetUser("test1", "test", "test")
	expected := chat.User{"test1", "test1"}

	if actual.Username != expected.Username || actual.Pass != expected.Pass {
		t.Errorf("Test failed, user not being created")
	}
}
