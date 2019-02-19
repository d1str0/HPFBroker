package main

import (
	"encoding/json"

	"github.com/d1str0/hpfeeds"
	bolt "go.etcd.io/bbolt"
)

var IDBucket = []byte("identities")

var BUCKETS = []string{
	string(IDBucket),
}

type KVStore interface {
	Get(key string) (interface{}, error)
	Put(key string, v interface{}) error
}

type BoltStore struct {
	db *bolt.DB
}

func (bs BoltStore) Get(key string) (interface{}, error) {
	var i interface{}
	err := bs.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(IDBucket)
		v := b.Get([]byte(key))
		err := json.Unmarshal(v, &i)
		return err
	})
	return i, err
}

func (bs BoltStore) Put(key string, i interface{}) error {
	err := bs.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(IDBucket)
		buf, err := json.Marshal(i)
		b.Put([]byte(key), buf)
		return err
	})
	return err
}

// Used to identify a user and their identity within hpfeeds broker.
func GetIdentity(bs BoltStore, ident string) (*hpfeeds.Identity, error) {
	var i hpfeeds.Identity
	err := bs.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("identities"))
		v := b.Get([]byte(ident))
		err := json.Unmarshal(v, &i)
		return err
	})
	return &i, err
}

func SaveIdentity(bs BoltStore, id hpfeeds.Identity) error {
	err := bs.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("identities"))
		buf, err := json.Marshal(id)
		b.Put([]byte(id.Ident), buf)
		return err
	})
	return err
}

func DeleteIdentity(bs BoltStore, ident string) error {
	err := bs.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("identities"))
		b.Put([]byte(ident), nil)
		return nil
	})
	return err
}
