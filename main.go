package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
)

func main() {
	password := flag.String("p", "", "password for MongoDB of \"sajmerdenr\" user, and cluster named \"cluster0\"")
	flag.Parse()

	app := Application{
		InfoLog:  log.New(os.Stdout, "INFO \t", log.Ldate|log.Ltime),
		ErrorLog: log.New(os.Stderr, "ERROR: \t", log.Ldate|log.Ltime|log.Lshortfile),
	}

	// Determining a secret key for signing tokens
	app.SecretKey = []byte(os.Getenv("SECRET_KEY"))
	if len(app.SecretKey) == 0 {
		app.SecretKey = []byte("super_secret")
		app.InfoLog.Printf("defaulting secret to \"%s\"\n", string(app.SecretKey))
	}

	// Connecting to MongoDB
	err := app.connectMongoDB(password)
	if err != nil {
		app.ErrorLog.Fatalln(err)
	}
	defer func() {
		if err = app.Client.Disconnect(context.TODO()); err != nil {
			app.ErrorLog.Fatalln(err)
		}
	}()

	// Determining a port for HTTP service.
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		app.InfoLog.Printf("defaulting to port %s", port)
	}

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: app.Routes(),
	}

	app.InfoLog.Println("Starting a server at http://localhost:" + port)
	app.ErrorLog.Fatalln(srv.ListenAndServe())

}
