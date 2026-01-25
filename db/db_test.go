package db

import (
	"testing"
)

var db, err = OpenDB(":memory:")

func TestInitDb(t *testing.T) {
	if err != nil {
		t.Fatalf("InitDB returned error: %v", err)
	}
	if db == nil {
		t.Fatal("InitDB returned a nil db")
	}
	defer db.Close()
}
