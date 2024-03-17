package main

import (
	"time"

	"github.com/golang-jwt/jwt"
)

func (app *Application) generateAccess(id string, secret []byte) (string, error) {
	claims := jwt.StandardClaims{
		Id:        id, // Id of a user that the access is being granted to
		IssuedAt:  time.Now().Unix(),
		ExpiresAt: time.Now().Add(time.Minute * 60).Unix(), // Expires 60 minutes after it was issued
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS512, claims) // signing method: HMAC SHA512
	return accessToken.SignedString(secret)
}

func (app *Application) generateRefresh(id string, secret []byte) (string, error) {
	claims := jwt.StandardClaims{
		Id:        id, // Id of a user that the refresh token is given to
		IssuedAt:  time.Now().Unix(),
		ExpiresAt: time.Now().Add(time.Hour * 24 * 7).Unix(), // Expires a week after it was issued
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS512, claims) // signing method: HMAC SHA512
	return accessToken.SignedString(secret)
}
