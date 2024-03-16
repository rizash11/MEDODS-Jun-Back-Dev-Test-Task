package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/golang-jwt/jwt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	password := flag.String("p", "", "password for MongoDB of \"sajmerdenr\" user, and cluster named \"cluster0\")")
	flag.Parse()

	client, err := connectMongoDB(password)
	if err != nil {
		log.Fatalln(err)
	}
	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			log.Fatalln(err)
		}
	}()

}

func connectMongoDB(password *string) (*mongo.Client, error) {
	// Use the SetServerAPIOptions() method to set the version of the Stable API on the client
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI("mongodb+srv://sajmerdenr:" + *password + "@cluster0.bxnjikl.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0").SetServerAPIOptions(serverAPI)
	// Create a new client and connect to the server
	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		return nil, err
	}

	// Send a ping to confirm a successful connection
	if err := client.Database("admin").RunCommand(context.TODO(), bson.D{{"ping", 1}}).Err(); err != nil {
		return nil, err
	}
	fmt.Println("Pinged your deployment. You successfully connected to MongoDB!")

	return client, nil
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
