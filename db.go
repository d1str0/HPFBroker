package hpfbroker

import (
	"encoding/json"
	"fmt"

	"github.com/boltdb/bolt"
	"github.com/d1str0/hpfeeds"
)

// Users for the web app
var UserBucket = []byte("users")

// Identities for hpfeeds
var IDBucket = []byte("identities")

var BUCKETS = []string{
	string(UserBucket),
	string(IDBucket),
}

type DB struct {
	*bolt.DB
}

// Open will open a bolt.DB and return a local DB
func OpenDB(path string) (*DB, error) {
	// Open up the boltDB file
	db, err := bolt.Open(path, 0666, nil)
	if err != nil {
		return nil, err
	}
	err = updateSchema(db)

	return &DB{db}, err
}

// updateSchema makes sure all necessary buckets exist and if any do not, they
// are created.
func updateSchema(db *bolt.DB) error {
	err := db.Update(func(tx *bolt.Tx) error {
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

// Identify will return an identity object if one is found.
func (db *DB) Identify(ident string) (*hpfeeds.Identity, error) {
	return db.GetIdentity(ident)
}

// GetAllIdentities returns a list of all hpfeeds Identity objects stored in the
// db.
func (db *DB) GetAllIdentities() ([]*hpfeeds.Identity, error) {
	var idents []*hpfeeds.Identity
	err := db.View(func(tx *bolt.Tx) error {
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
func (db *DB) GetIdentity(ident string) (*hpfeeds.Identity, error) {
	var i *hpfeeds.Identity
	err := db.View(func(tx *bolt.Tx) error {
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
func (db *DB) SaveIdentity(id *hpfeeds.Identity) error {
	err := db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(IDBucket)
		buf, err := json.Marshal(id)
		b.Put([]byte(id.Ident), buf)
		return err
	})
	return err
}

// DeleteIdentity removes any saved Identity object matching the ident.
func (db *DB) DeleteIdentity(ident string) error {
	err := db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(IDBucket)
		b.Delete([]byte(ident))
		return nil
	})
	return err
}

// DeleteAllIdeneties deletes the Bolt bucket holding identities and recreates
// it, essentially deleting all objects.
func (db *DB) DeleteAllIdentities() error {
	err := db.Update(func(tx *bolt.Tx) error {
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
func (db *DB) GetAllUsers() ([]*User, error) {
	var users []*User
	err := db.View(func(tx *bolt.Tx) error {
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
func (db *DB) GetUser(name string) (*User, error) {
	var u *User
	err := db.View(func(tx *bolt.Tx) error {
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
func (db *DB) SaveUser(u *User) error {
	err := db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(UserBucket)
		buf, err := json.Marshal(u)
		b.Put([]byte(u.Name), buf)
		return err
	})
	return err
}

// DeleteUser removes any saved User object matching the username
func (db *DB) DeleteUser(name string) error {
	err := db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(UserBucket)
		b.Delete([]byte(name))
		return nil
	})
	return err
}

// DeleteAllUsers deletes the Bolt bucket holding users and recreates
// it, essentially deleting all objects.
func (db *DB) DeleteAllUsers() error {
	err := db.Update(func(tx *bolt.Tx) error {
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
