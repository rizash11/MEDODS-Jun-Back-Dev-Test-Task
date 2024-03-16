package main

import (
	"fmt"
	"log"
	"time"

	"github.com/golang-jwt/jwt"
)

func main() {
	jwt, err := generateAccess("new user", []byte("secret"))
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(jwt)
}

func generateAccess(id string, secret []byte) (string, error) {
	claims := jwt.StandardClaims{
		Id:        id, // Id of a user that the access is being granted to
		IssuedAt:  time.Now().Unix(),
		ExpiresAt: time.Now().Add(time.Minute * 60).Unix(), // Expires 60 minutes after it was issued
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS512, claims) // signing method: HMAC SHA512
	return accessToken.SignedString(secret)
}

func generateRefresh(id string, secret []byte) (string, error) {
	claims := jwt.StandardClaims{
		Id:        id, // Id of a user that the refresh token is given to
		IssuedAt:  time.Now().Unix(),
		ExpiresAt: time.Now().Add(time.Hour * 24 * 7).Unix(), // Expires a week after it was issued
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS512, claims) // signing method: HMAC SHA512
	return accessToken.SignedString(secret)
}
