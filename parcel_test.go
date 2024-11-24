package main

import (
	"database/sql"
	"testing"
	_ "github.com/mattn/go-sqlite3"
)

func setupTestDB() (*sql.DB, func()) {
	db, _ := sql.Open("sqlite3", ":memory:")
	db.Exec(`CREATE TABLE parcel (
		number INTEGER PRIMARY KEY AUTOINCREMENT,
		client INTEGER,
		status TEXT,
		address TEXT,
		created_at TEXT
	)`)
	return db, func() { db.Close() }
}

func TestRegisterParcel(t *testing.T) {
	db, teardown := setupTestDB()
	defer teardown()

	store := &ParcelStore{DB: db}
	id, err := store.RegisterParcel(1, "123 Main St")
	if err != nil {
		t.Fatalf("failed to register parcel: %v", err)
	}
	if id == 0 {
		t.Fatalf("expected valid parcel ID, got 0")
	}
}

func TestGetParcelsByClient(t *testing.T) {
	db, teardown := setupTestDB()
	defer teardown()

	store := &ParcelStore{DB: db}
	store.RegisterParcel(1, "123 Main St")
	store.RegisterParcel(1, "456 Elm St")

	parcels, err := store.GetParcelsByClient(1)
	if err != nil || len(parcels) != 2 {
		t.Fatalf("failed to fetch parcels: %v", err)
	}
}