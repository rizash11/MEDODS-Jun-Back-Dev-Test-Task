package main

import (
	"context"
	"encoding/base64"
	"time"

	"github.com/golang-jwt/jwt"
)

type refreshToken struct {
	Id      string
	Refresh string
}

type tokens struct {
	Access  string
	Refresh string
}

func (app *Application) generateTokens(guid string) (*tokens, error) {
	var newTokens tokens
	var err error
	var hash []byte
	var enc string

	// generating an access token
	newTokens.Access, err = app.generateAccess(guid)
	if err != nil {
		return nil, err
	}

	// generating a refresh token
	newTokens.Refresh, err = app.generateRefresh(guid)
	if err != nil {
		return nil, err
	}

	// bcrypt hashing the refresh token
	hash, err = HashPassword(newTokens.Refresh)
	if err != nil {
		return nil, err
	}
	// base64 encoding the refresh token
	enc = base64.StdEncoding.EncodeToString(hash)

	// connecting to the database and creating a new row for the refresh token
	collection := app.Client.Database("OAuth").Collection("refresh_tokens")
	newDbRow := refreshToken{
		Id:      guid,
		Refresh: enc,
	}

	// inserting the refresh token
	_, err = collection.InsertOne(context.TODO(), newDbRow)
	if err != nil {
		return nil, err
	}

	return &newTokens, nil
}

func (app *Application) generateAccess(id string) (string, error) {
	claims := jwt.StandardClaims{
		Id:        id, // Id of a user that the access is being granted to
		IssuedAt:  time.Now().Unix(),
		ExpiresAt: time.Now().Add(time.Minute * 60).Unix(), // Expires 60 minutes after it was issued
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS512, claims) // signing method: HMAC SHA512
	return accessToken.SignedString(app.SecretKey)
}

func (app *Application) generateRefresh(id string) (string, error) {
	claims := jwt.StandardClaims{
		Id:        id, // Id of a user that the refresh token is given to
		IssuedAt:  time.Now().Unix(),
		ExpiresAt: time.Now().Add(time.Hour * 24 * 7).Unix(), // Expires a week after it was issued
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS512, claims) // signing method: HMAC SHA512
	return accessToken.SignedString(app.SecretKey)
}
