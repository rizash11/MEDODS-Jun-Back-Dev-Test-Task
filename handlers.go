package main

import (
	"html/template"
	"net/http"
	"strings"
)

// Keeps all of the handlers of the app.
func (app *Application) Routes() *http.ServeMux {
	router := http.NewServeMux()
	router.HandleFunc("/", app.home)
	router.HandleFunc("/generate-tokens", app.generateTokensRoute)
	// router.Handle("/frontend/", http.StripPrefix("/frontend", http.FileServer(http.Dir("./frontend/"))))

	return router
}

func (app *Application) home(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.URL.Path != "/":
		app.NotFound(w)
		return

	case r.Method != http.MethodGet:
		app.ClientError(w, http.StatusMethodNotAllowed)
		return
	}

	templ, err := template.ParseFiles("ui/html/home.page.html")
	if err != nil {
		app.ServerError(w, err)
	}

	templ.Execute(w, nil)
	// app.Render(w, r, "cultivation.page.html")
}

func (app *Application) generateTokensRoute(w http.ResponseWriter, r *http.Request) {
	switch {
	case strings.ToLower(r.URL.Path) != "/generate-tokens":
		app.NotFound(w)
		return

	case r.Method != http.MethodPost:
		app.ClientError(w, http.StatusMethodNotAllowed)
		return
	}

	guid := r.FormValue("generate-input")
	if guid == "" {
		app.ClientError(w, http.StatusBadRequest)
	}

	tokens, err := app.generateTokens(guid)
	if err != nil {
		app.ServerError(w, err)
	}

	templ, err := template.ParseFiles("ui/html/home.page.html")
	if err != nil {
		app.ServerError(w, err)
	}
	templ.Execute(w, tokens)
}
