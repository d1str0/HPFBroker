package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/d1str0/hpfeeds"
	bolt "go.etcd.io/bbolt"
)

const version = "v0.0.1"

var path string

// To be passed to various http handlers.
type App struct {
	DB   *bolt.DB
	Name string
}

func main() {
	fmt.Println("///- Starting up HPFBroker")
	fmt.Printf("//- Version %s\n", version)

	// Grab any command line arguments
	flag.StringVar(&path, "db", "bolt.db", "File path for the BoltDB store file.")
	flag.Parse()

	// Open up the boltDB file
	db, err := bolt.Open(path, 0666, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// KVStore for use with hpfeeds broker
	kv := KVStore{DB: db}
	// App for use with http handlers
	app := App{DB: db, Name: "test"}

	// Prepare DB to ensure we have the appropriate buckets ready
	err = initializeDB(kv)
	if err != nil {
		log.Fatal(err)
	}

	// Configure hpfeeds broker server
	broker := &hpfeeds.Broker{
		Name: "HPFBroker",
		Port: 10000,
		DB:   kv,
	}
	broker.SetDebugLogger(log.Print)
	broker.SetInfoLogger(log.Print)
	broker.SetErrorLogger(log.Print)

	// Run http server concurrently
	go func() {
		// Load routes for the server
		mux := routes(app)

		s := http.Server{
			Addr:    ":8080",
			Handler: mux,
		}
		log.Fatal(s.ListenAndServe())
	}()

	// Start hpfeeds broker server
	err = broker.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}

type KVStore struct {
	DB *bolt.DB
}

var BUCKETS = []string{
	"identities",
}

// Initialize the database and assert certain buckets exist.
func initializeDB(kv KVStore) error {
	err := kv.DB.Update(func(tx *bolt.Tx) error {
		for _, b := range BUCKETS {
			_, err := tx.CreateBucketIfNotExists([]byte(b))
			if err != nil {
				return fmt.Errorf("create bucket: %s", err)
			}
		}
		return nil
	})
	return err
}

func (kv KVStore) Identify(ident string) (*hpfeeds.Identity, error) {
	var i hpfeeds.Identity
	err := kv.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("identities"))
		v := b.Get([]byte(ident))
		err := json.Unmarshal(v, &i)
		return err
	})
	return &i, err
}
