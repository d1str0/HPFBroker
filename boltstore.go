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

// GetAllUsers returns a list of all hpfeeds Identity objects stored in the
// db.
func (bs BoltStore) GetAllUsers() ([]*User, error) {
	var users []*User
	err := bs.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(UserBucket)

		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			u := &User{}
			err := json.Unmarshal(v, &u)
			if err != nil {
				return err
			}

			users = append(users, u)
		}

		return nil
	})
	return users, err
}

// GetUser takes an username and returns their whole user object.
func (bs BoltStore) GetUser(name string) (*User, error) {
	var u *User
	err := bs.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(UserBucket)
		v := b.Get([]byte(name))
		if v == nil {
			return nil
		}
		u = &User{}
		err := json.Unmarshal(v, &u)
		return err
	})
	return u, err
}

// SaveUser persists a User in BoltStore.
func (bs BoltStore) SaveUser(u User) error {
	err := bs.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(UserBucket)
		buf, err := json.Marshal(u)
		b.Put([]byte(u.Name), buf)
		return err
	})
	return err
}

// DeleteUser removes any saved User object matching the username
func (bs BoltStore) DeleteUser(name string) error {
	err := bs.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(UserBucket)
		b.Delete([]byte(name))
		return nil
	})
	return err
}

// DeleteAllUsers deletes the Bolt bucket holding users and recreates
// it, essentially deleting all objects.
func (bs BoltStore) DeleteAllUsers() error {
	err := bs.db.Update(func(tx *bolt.Tx) error {
		err := tx.DeleteBucket(UserBucket)
		if err != nil {
			return err
		}
		_, err = tx.CreateBucket(UserBucket)
		if err != nil {
			return err
		}

		return nil
	})
	return err
}
