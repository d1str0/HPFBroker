package main

import (
	"fmt"
	"net/http"
)

func rootHandler(app App) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Welcome to the home page of %s!", app.Name)
	}
}

func routes(app App) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", rootHandler(app))
	return mux
}
