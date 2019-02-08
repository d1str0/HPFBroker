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

func apiIdentHandler(app App) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ident := r.URL.Path[len("/api/ident/"):]
		i, err := GetIdentity(app, ident)

		fmt.Fprintf(w, "This is where you GET an asdfdent.\n%#v\n%#v)", i, err)
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
