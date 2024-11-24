package main

import (
	"database/sql"
	"fmt"
	"log"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db, err := sql.Open("sqlite3", "tracker.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	store := &ParcelStore{DB: db}

	clientID := 123
	parcelID, err := store.RegisterParcel(clientID, "123 Main St")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Registered parcel with ID %d\n", parcelID)

	parcels, err := store.GetParcelsByClient(clientID)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Parcels for client %d: %+v\n", clientID, parcels)

	err = store.UpdateParcelAddress(parcelID, "456 Elm St")
	if err != nil {
		log.Println("Error updating address:", err)
	} else {
		fmt.Println("Updated address")
	}

	err = store.UpdateParcelStatus(parcelID, ParcelStatusSent)
	if err != nil {
		log.Println("Error updating status:", err)
	} else {
		fmt.Println("Updated status to sent")
	}

	err = store.DeleteParcel(parcelID)
	if err != nil {
		log.Println("Error deleting parcel:", err)
	} else {
		fmt.Println("Deleted parcel")
	}
}