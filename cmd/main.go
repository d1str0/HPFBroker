package main

import (
	hpf "github.com/d1str0/HPFBroker"

	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/d1str0/hpfeeds"
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

	// Open our database file.
	db, err := hpf.OpenDB(dbc.Path)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	bc := t.BrokerConfig
	// Configure hpfeeds broker server
	broker := &hpfeeds.Broker{
		Name: bc.Name,
		Port: bc.Port,
		DB:   db,
	}
	broker.SetDebugLogger(log.Print)
	broker.SetInfoLogger(log.Print)
	broker.SetErrorLogger(log.Print)

	hc := t.HttpConfig
	// Run http server concurrently
	go func() {
		// Load routes for the server
		mux := hpf.NewMux(db)

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
