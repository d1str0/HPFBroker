package main

import (
	"fmt"
	"log"

	"github.com/d1str0/hpfeeds"
	bolt "go.etcd.io/bbolt"
)

const version = "v0.0.1"

var path = "bolt.db"

func main() {
	fmt.Println("///- Starting up HPFBroker")
	fmt.Printf("//- Version %s\n", version)

	db, err := bolt.Open(path, 0666, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Prepare DB and store as IdentityDB
	idb, err := initializeDB(db)
	if err != nil {
		log.Fatal(err)
	}

	b := &hpfeeds.Broker{
		Name: "HPFBroker",
		Port: 10000,
		DB:   idb,
	}
	b.SetDebugLogger(log.Print)
	b.SetInfoLogger(log.Print)
	b.SetErrorLogger(log.Print)

	err = b.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}

type IdentityDB struct {
	DB *bolt.DB
}

var BUCKETS = []string{
	"identities",
}

// Initialize the database and assert certain buckets exist.
func initializeDB(db *bolt.DB) (*IdentityDB, error) {
	err := db.Update(func(tx *bolt.Tx) error {
		for _, b := range BUCKETS {
			_, err := tx.CreateBucketIfNotExists([]byte(b))
			if err != nil {
				return fmt.Errorf("create bucket: %s", err)
			}
		}
		return nil
	})
	return &IdentityDB{DB: db}, err
}

// TODO: Replace this function with something that checks the DB for current identities.
func (db *IdentityDB) Identify(ident string) (*hpfeeds.Identity, error) {
	return &hpfeeds.Identity{
		Ident:       "test",
		Secret:      "test",
		SubChannels: []string{"test"},
		PubChannels: []string{"test"},
	}, nil
}
