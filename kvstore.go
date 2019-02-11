package main

import (
	"encoding/json"

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
