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

type BoltStore struct {
	db *bolt.DB
}

func (bs BoltStore) Close() {
	bs.db.Close()
}

func (bs BoltStore) GetKeys() ([]string, error) {
	var keys []string
	err := bs.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(IDBucket)

		c := b.Cursor()

		for k, _ := c.First(); k != nil; k, _ = c.Next() {
			keys = append(keys, string(k))
		}

		return nil
	})
	return keys, err
}

func (bs BoltStore) GetAllIdentities() ([]*hpfeeds.Identity, error) {
	var idents []*hpfeeds.Identity
	err := bs.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(IDBucket)

		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			i := &hpfeeds.Identity{}
			err := json.Unmarshal(v, &i)
			if err != nil {
				return err
			}

			idents = append(idents, i)
		}

		return nil
	})
	return idents, err
}

// Used to identify a user and their identity within hpfeeds broker.
func GetIdentity(bs BoltStore, ident string) (*hpfeeds.Identity, error) {
	var i *hpfeeds.Identity
	err := bs.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(IDBucket)
		v := b.Get([]byte(ident))
		if v == nil {
			return nil
		}
		i = &hpfeeds.Identity{}
		err := json.Unmarshal(v, &i)
		return err
	})
	return i, err
}

func SaveIdentity(bs BoltStore, id hpfeeds.Identity) error {
	err := bs.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(IDBucket)
		buf, err := json.Marshal(id)
		b.Put([]byte(id.Ident), buf)
		return err
	})
	return err
}

func DeleteIdentity(bs BoltStore, ident string) error {
	err := bs.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(IDBucket)
		b.Delete([]byte(ident))
		return nil
	})
	return err
}
