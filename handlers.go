package main

import (
	"html/template"
	"net/http"
)

// Keeps all of the handlers of the app.
func (app *Application) Routes() *http.ServeMux {
	router := http.NewServeMux()
	router.HandleFunc("/", app.home)
	router.HandleFunc("/generateTokens", app.generateTokens)
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

func (app *Application) generateTokens(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.URL.Path != "/":
		app.NotFound(w)
		return

	case r.Method != http.MethodGet:
		app.ClientError(w, http.StatusMethodNotAllowed)
		return
	}

}
