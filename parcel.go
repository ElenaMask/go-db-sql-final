package main

import (
	"database/sql"
	"errors"
	"time"
)

const (
	ParcelStatusRegistered = "registered"
	ParcelStatusSent       = "sent"
	ParcelStatusDelivered  = "delivered"
)

type Parcel struct {
	Number    int
	Client    int
	Status    string
	Address   string
	CreatedAt string
}

type ParcelStore struct {
	DB *sql.DB
}

func (ps *ParcelStore) RegisterParcel(client int, address string) (int, error) {
	result, err := ps.DB.Exec(`
		INSERT INTO parcel (client, status, address, created_at)
		VALUES (?, ?, ?, ?)`,
		client, ParcelStatusRegistered, address, time.Now().Format(time.RFC3339),
	)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil
}

func (ps *ParcelStore) GetParcelsByClient(client int) ([]Parcel, error) {
	rows, err := ps.DB.Query(`SELECT number, client, status, address, created_at FROM parcel WHERE client = ?`, client)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var parcels []Parcel
	for rows.Next() {
		var p Parcel
		if err := rows.Scan(&p.Number, &p.Client, &p.Status, &p.Address, &p.CreatedAt); err != nil {
			return nil, err
		}
		parcels = append(parcels, p)
	}
	return parcels, nil
}

func (ps *ParcelStore) UpdateParcelStatus(number int, status string) error {
	_, err := ps.DB.Exec(`UPDATE parcel SET status = ? WHERE number = ?`, status, number)
	return err
}

func (ps *ParcelStore) UpdateParcelAddress(number int, address string) error {
	var status string
	err := ps.DB.QueryRow(`SELECT status FROM parcel WHERE number = ?`, number).Scan(&status)
	if err != nil {
		return err
	}
	if status != ParcelStatusRegistered {
		return errors.New("address can only be updated for registered parcels")
	}
	_, err = ps.DB.Exec(`UPDATE parcel SET address = ? WHERE number = ?`, address, number)
	return err
}

func (ps *ParcelStore) DeleteParcel(number int) error {
	var status string
	err := ps.DB.QueryRow(`SELECT status FROM parcel WHERE number = ?`, number).Scan(&status)
	if err != nil {
		return err
	}
	if status != ParcelStatusRegistered {
		return errors.New("only registered parcels can be deleted")
	}
	_, err = ps.DB.Exec(`DELETE FROM parcel WHERE number = ?`, number)
	return err
}