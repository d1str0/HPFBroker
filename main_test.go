package main

import (
	"testing"

	bolt "go.etcd.io/bbolt"
)

const TestDBPath = "test.db"

func getTestDB(t *testing.T) BoltStore {
	// Open up the boltDB file
	db, err := bolt.Open(TestDBPath, 0666, nil)
	if err != nil {
		t.Fatalf("Error opening db: %s", err.Error())
	}

	// For use with HTTP handlers
	bs := BoltStore{db: db}

	// Prepare DB to ensure we have the
	// appropriate buckets ready
	err = initializeDB(bs)
	if err != nil {
		t.Fatalf("Error initializing db: %s", err.Error())
	}
	return bs
}
