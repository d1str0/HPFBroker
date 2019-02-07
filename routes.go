package main

import (
	"fmt"
	"net/http"
)

func statusHandler(app App) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Welcome to the home page of %s!", Version)
	}
}

func routes(app App) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/status", statusHandler(app))
	//mux.HandleFunc("/api", apiHandler(app))
	return mux
}
