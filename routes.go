package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/d1str0/hpfeeds"
	bolt "go.etcd.io/bbolt"
)

func statusHandler(app App) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Welcome to the home page of %s!", Version)
	}
}

//TODO return proer http codes per method
func apiIdentHandler(app App) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ident := r.URL.Path[len("/api/ident/"):]

		// Handle API requests depending on HTTP Method
		switch r.Method {
		case http.MethodGet:
			// Return identity if found
			i, err := GetIdentity(app, ident)
			buf, err := json.Marshal(i)
			if err != nil {
				http.Error(w, err.Error(), 500)
			} else {
				if i.Ident == "" {
					http.Error(w, "Ident not found", 404)
				} else {
					fmt.Fprintf(w, "%s", buf)
				}
			}
		case http.MethodPut:
			// Update user
			var id hpfeeds.Identity
			if r.Body == nil {
				http.Error(w, "Request body required", 400)
				return
			}

			err := json.NewDecoder(r.Body).Decode(&id)
			if err != nil {
				http.Error(w, err.Error(), 400)
				return
			}
			if ident != id.Ident {
				http.Error(w, "Request body not valid for this URI. Ident mismatch.", 400)
				return
			}

			err = SaveIdentity(app, id)
			if err != nil {
				http.Error(w, err.Error(), 500)
			}
		case http.MethodDelete:
			// Delete user
			err := DeleteIdentity(app, ident)
			if err != nil {
				http.Error(w, err.Error(), 500)
				return
			}
		}
	}

}

func routes(app App) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/status", statusHandler(app))
	mux.HandleFunc("/api/ident/", apiIdentHandler(app))
	return mux
}

// Used to identify a user and their identity within hpfeeds broker.
func GetIdentity(app App, ident string) (*hpfeeds.Identity, error) {
	var i hpfeeds.Identity
	err := app.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("identities"))
		v := b.Get([]byte(ident))
		err := json.Unmarshal(v, &i)
		return err
	})
	return &i, err
}

func SaveIdentity(app App, id hpfeeds.Identity) error {
	err := app.DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("identities"))
		buf, err := json.Marshal(id)
		b.Put([]byte(id.Ident), buf)
		return err
	})
	return err
}

func DeleteIdentity(app App, ident string) error {
	err := app.DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("identities"))
		b.Put([]byte(ident), nil)
		return nil
	})
	return err
}
