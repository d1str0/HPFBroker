package main

import (
	"fmt"
	"log"

	"github.com/d1str0/hpfeeds"
)

const version = "v0.0.1"

func main() {
	fmt.Println("///- Starting up HPFBroker")
	fmt.Printf("//- Version %s\n", version)

	db := &IdentityDB{}

	b := &hpfeeds.Broker{
		Name: "HPFBroker",
		Port: 10000,
		DB:   db,
	}
	b.SetDebugLogger(log.Print)
	b.SetInfoLogger(log.Print)
	b.SetErrorLogger(log.Print)

	err := b.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}

type IdentityDB struct {
}

func (db *IdentityDB) Identify(ident string) (*hpfeeds.Identity, error) {
	return &hpfeeds.Identity{
		Ident:       "test",
		Secret:      "test",
		SubChannels: []string{"test"},
		PubChannels: []string{"test"},
	}, nil
}
