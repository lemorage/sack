package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"text/template"
)

func startServer(mux *http.ServeMux, port int) {
	addr := fmt.Sprintf(":%d", port)
	log.Printf("Starting server on %s...\n", addr)
	err := http.ListenAndServe(addr, mux)
	if err != nil {
		log.Fatalf("Error starting server: %s", err)
	}
}

func parseTemplates() *template.Template {
	funcMap := template.FuncMap{
		"add": func(i int) int {
			return i + 1
		},
		"sub": func(i int) int {
			return i - 1
		},
	}
	return template.Must(template.New("base").Funcs(funcMap).ParseGlob("ui/html/*.gohtml"))
}

func generateHTMLFiles(config Config, tmpl *template.Template, layout string) {
	dir := "./ui/html/pages"
	keys := sortedPageKeys(config.Pages)

	for _, key := range keys {
		pageConfig := config.Pages[key]
		pageNumber, _ := extractNumber(key)
		pageFilename := fmt.Sprintf("%s/page%d.gohtml", dir, pageNumber)
		newPage, err := os.Create(pageFilename)
		if err != nil {
			log.Fatalf("Error creating page file for %s: %s", key, err)
		}
		defer newPage.Close()

		err = tmpl.ExecuteTemplate(newPage, "base", struct {
			CurrentPage int
			TotalPages  int
			PageConfig  PageConfig
			Layout      string
		}{
			CurrentPage: pageNumber,
			TotalPages:  len(config.Pages),
			PageConfig:  pageConfig,
			Layout:      layout,
		})
		if err != nil {
			log.Fatalf("Error executing template for page %s: %s", key, err)
		}
		log.Printf("Generated HTML for %s\n", key)
	}
}

func setupHandlers(config Config) *http.ServeMux {
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./ui/static"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	for i := 1; i <= len(config.Pages); i++ {
		pageFilename := fmt.Sprintf("./ui/html/pages/page%d.gohtml", i)
		mux.HandleFunc("/model"+strconv.Itoa(i), func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, pageFilename)
		})
	}

	mux.HandleFunc("/", home)
	return mux
}
