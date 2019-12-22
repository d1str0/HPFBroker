package main

import (
	"testing"

	"github.com/d1str0/hpfeeds"
)

func TestKvstore_BoltStore(t *testing.T) {
	bs := getTestDB(t)
	defer bs.Close()

	t.Run("Get Nonexistant", func(t *testing.T) {
		i, err := GetIdentity(bs, "test")
		if err != nil {
			t.Fatal(err)
		}

		if i != nil {
			t.Error("Expected nil identity returned.")
		}
	})

	id1 := hpfeeds.Identity{Ident: "test-ident", Secret: "test-secret", SubChannels: []string{}, PubChannels: []string{}}

	t.Run("Save Identity", func(t *testing.T) {
		err := SaveIdentity(bs, id1)
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("Get Existing", func(t *testing.T) {
		i, err := GetIdentity(bs, "test-ident")
		if err != nil {
			t.Fatal(err)
		}

		if i == nil {
			t.Error("Unexpected nil identity returned.")
		}

		// expected, got
		assertEqualIdentity(t, id1, *i)
	})

	t.Run("Delete", func(t *testing.T) {
		err := DeleteIdentity(bs, "test-ident")
		if err != nil {
			t.Fatal(err)
		}
	})

}

func assertEqualIdentity(t *testing.T, expect hpfeeds.Identity, got hpfeeds.Identity) {
	if expect.Ident != got.Ident {
		t.Errorf("Mismatched Idents:\n\tgot %s \n\twant %s", got.Ident, expect.Ident)
	}
	if expect.Secret != got.Secret {
		t.Errorf("Mismatched Secrets:\n\tgot %s \n\twant %s", got.Secret, expect.Secret)
	}

	if !testEq(expect.SubChannels, got.SubChannels) {
		t.Errorf("Mismatched SubChannels:\n\tgot %v \n\twant %v", got.SubChannels, expect.SubChannels)
	}
	if !testEq(expect.PubChannels, got.PubChannels) {
		t.Errorf("Mismatched PubChannels:\n\tgot %v \n\twant %v", got.PubChannels, expect.PubChannels)
	}
}

func testEq(a, b []string) bool {

	// If one is nil, the other must also be nil.
	if (a == nil) != (b == nil) {
		return false
	}

	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}