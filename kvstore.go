package main

import (
	"encoding/json"

	"github.com/d1str0/hpfeeds"
	bolt "go.etcd.io/bbolt"
)

// Users for the web app
var UserBucket = []byte("users")

// Identities for hpfeeds
var IDBucket = []byte("identities")

var BUCKETS = []string{
	string(UserBucket),
	string(IDBucket),
}

type BoltStore struct {
	db *bolt.DB
}

// Close calls the underlying Bolt db Close()
func (bs BoltStore) Close() {
	bs.db.Close()
}

// GetAllIdentities returns a list of all hpfeeds Identity objects stored in the
// db.
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

// GetIdentity takes an ident and returns their whole identity object.
func (bs BoltStore) GetIdentity(ident string) (*hpfeeds.Identity, error) {
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

// SaveIdentity persists an hpfeeds.Identity in BoltStore.
func (bs BoltStore) SaveIdentity(id hpfeeds.Identity) error {
	err := bs.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(IDBucket)
		buf, err := json.Marshal(id)
		b.Put([]byte(id.Ident), buf)
		return err
	})
	return err
}

// DeleteIdentity removes any saved Identity object matching the ident.
func (bs BoltStore) DeleteIdentity(ident string) error {
	err := bs.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(IDBucket)
		b.Delete([]byte(ident))
		return nil
	})
	return err
}

// DeleteAllIdeneties deletes the Bolt bucket holding identities and recreates
// it, essentially deleting all objects.
func (bs BoltStore) DeleteAllIdentities() error {
	err := bs.db.Update(func(tx *bolt.Tx) error {
		err := tx.DeleteBucket(IDBucket)
		if err != nil {
			return err
		}
		_, err = tx.CreateBucket(IDBucket)
		if err != nil {
			return err
		}

		return nil
	})
	return err
}
