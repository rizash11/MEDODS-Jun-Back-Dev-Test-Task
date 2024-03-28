package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"runtime/debug"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

type Application struct {
	InfoLog  *log.Logger
	ErrorLog *log.Logger
	// TemplateCache map[string]*template.Template
	Client    *mongo.Client
	SecretKey []byte
}

func (app *Application) connectMongoDB(password *string) error {
	// Use the SetServerAPIOptions() method to set the version of the Stable API on the client
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI("mongodb+srv://sajmerdenr:" + *password + "@cluster0.bxnjikl.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0").SetServerAPIOptions(serverAPI)
	// Create a new client and connect to the server
	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		return err
	}

	// Send a ping to confirm a successful connection
	if err := client.Database("admin").RunCommand(context.TODO(), bson.D{{"ping", 1}}).Err(); err != nil {
		return err
	}
	app.InfoLog.Println("Pinged your deployment. You successfully connected to MongoDB!")
	app.Client = client

	return nil
}

func (app *Application) NotFound(w http.ResponseWriter) {
	http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
}

func (app *Application) ClientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (app *Application) ServerError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.ErrorLog.Output(2, trace)

	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

// HashPassword generates a bcrypt hash for the given password.
func HashPassword(password string) ([]byte, error) {
	if len(password) > 72 {
		password = password[len(password)-72:]
	}

	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return bytes, err
}

// VerifyPassword verifies if the given password matches the stored hash.
func VerifyPassword(password, hash string) bool {
	if len(password) > 72 {
		password = password[len(password)-72:]
	}

	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
