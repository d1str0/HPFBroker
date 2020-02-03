package main

import (
	hpf "github.com/d1str0/HPFBroker"
	api "github.com/d1str0/HPFBroker/api"
	auth "github.com/d1str0/HPFBroker/auth"

	"encoding/base64"
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
	Addr          string
	SigningSecret string // For JWT Signing
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

	secret, err := base64.StdEncoding.DecodeString(hc.SigningSecret)
	if err != nil {
		log.Fatal(err)
	}

	jwt := &auth.JWTSecret{}
	jwt.SetSecret(secret)

	r := auth.InitRBAC()
	sc := &api.ServerContext{Version: Version, JWTSecret: jwt, DB: db, RBAC: r}

	// Run http server concurrently
	// Load routes for the server
	mux := api.NewMux(sc)

	s := http.Server{
		Addr:    hc.Addr,
		Handler: mux,
	}

	go func() {
		log.Fatal(s.ListenAndServe())
	}()

	// Start hpfeeds broker server
	log.Fatal(broker.ListenAndServe())
}
