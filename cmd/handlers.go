package main

import (
	"html/template"
	"log"
	"net/http"
)

// home handler for the home page
func home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		notFound(w)
		return
	}

	ts, err := template.ParseFiles("./ui/html/index.html")
	if err != nil {
		serverError(w, err)
		return
	}

	err = ts.Execute(w, nil)
	if err != nil {
		serverError(w, err)
	}
}

// notFound handler for custom 404 page
func notFound(w http.ResponseWriter) {
	ts, err := template.ParseFiles("./ui/html/404.html")
	if err != nil {
		serverError(w, err)
		return
	}

	w.WriteHeader(http.StatusNotFound)
	err = ts.Execute(w, nil)
	if err != nil {
		serverError(w, err)
	}
}

// serverError handler for custom 500 page
func serverError(w http.ResponseWriter, err error) {
	log.Print(err.Error())
	ts, err := template.ParseFiles("./ui/html/500.html")
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusInternalServerError)
	err = ts.Execute(w, nil)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
