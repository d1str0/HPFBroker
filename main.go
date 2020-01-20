package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/d1str0/hpfeeds"
	bolt "go.etcd.io/bbolt"
)

const Version = "v0.0.1"

var configFilename string

// Configuration for any BoltDB options
type DBConfig struct {
	Path string
}

// Configuration for the HPFeeds broker srver.
type BrokerConfig struct {
	Name string
	Port int
}

// Configuration for the webserver
type HttpConfig struct {
	Addr string
	//SessionSecret string // For Gorilla sessions
}

type tomlConfig struct {
	DBConfig     `toml:"database"`
	BrokerConfig `toml:"hpfeeds"`
	HttpConfig   `toml:"http"`
}

func main() {
	fmt.Println("///- Starting up HPFBroker")
	fmt.Printf("//- Version %s\n", Version)

	// Grab any command line arguments
	flag.StringVar(&configFilename, "config", "config.toml", "File path for the config file (TOML).")
	flag.Parse()

	// TODO: Rename this var.
	var t tomlConfig

	_, err := toml.DecodeFile(configFilename, &t)
	if err != nil {
		if _, ok := err.(*os.PathError); ok {
			log.Fatal("Config file not found.")
		} else {
			log.Fatal(err.Error())
		}
	}

	dbc := t.DBConfig

	// Open up the boltDB file
	db, err := bolt.Open(dbc.Path, 0666, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// For use with HTTP handlers
	bs := BoltStore{db: db}

	// Prepare DB to ensure we have the appropriate buckets ready
	err = initializeDB(bs)
	if err != nil {
		log.Fatal(err)
	}

	bc := t.BrokerConfig
	// Configure hpfeeds broker server
	broker := &hpfeeds.Broker{
		Name: bc.Name,
		Port: bc.Port,
		DB:   bs,
	}
	broker.SetDebugLogger(log.Print)
	broker.SetInfoLogger(log.Print)
	broker.SetErrorLogger(log.Print)

	hc := t.HttpConfig
	// Run http server concurrently
	go func() {
		// Load routes for the server
		mux := NewMux(bs)

		s := http.Server{
			Addr:    hc.Addr,
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

// Initialize the database and assert certain buckets exist.
func initializeDB(bs BoltStore) error {
	err := bs.db.Update(func(tx *bolt.Tx) error {
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

func (bs BoltStore) Identify(ident string) (*hpfeeds.Identity, error) {
	var i hpfeeds.Identity
	err := bs.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("identities"))
		v := b.Get([]byte(ident))
		err := json.Unmarshal(v, &i)
		return err
	})
	return &i, err
}
