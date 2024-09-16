package main

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

type Parcel struct {
	Number    int
	Client    int
	Status    string
	Address   string
	CreatedAt string
}
type ParcelStore struct {
	db *sql.DB
}

func NewParcelStore(db *sql.DB) ParcelStore {
	return ParcelStore{db: db}
}

func (s ParcelStore) Add(p Parcel) (int, error) {
	result, err := s.db.Exec(`
		INSERT INTO parcel (client, status, address, created_at)
		VALUES (?, ?, ?, ?)
	`, p.Client, p.Status, p.Address, p.CreatedAt)
	if err != nil {
		return 0, nil
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil
}

func (s ParcelStore) Get(number int) (Parcel, error) {
	row := s.db.QueryRow(`
		SELECT number, client, status, address, created_at
		FROM parcel
		WHERE number = ?
	`, number)

	var p Parcel
	err := row.Scan(&p.Number, &p.Client, &p.Status, &p.Address, &p.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return p, fmt.Errorf("parcel not found")
		}
		return p, err
	}
	return p, nil
}

func (s ParcelStore) GetByClient(client int) ([]Parcel, error) {
	rows, err := s.db.Query(`
		SELECT number, client, status, address, created_at
		FROM parcel
		WHERE client = ?
	`, client)
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

func (s ParcelStore) SetStatus(number int, status string) error {
	_, err := s.db.Exec(`
		UPDATE parcel
		SET status = ?
		WHERE number = ?
	`, status, number)
	return err
}

func (s ParcelStore) SetAddress(number int, address string) error {
	// Проверяем текущий статус посылки
	row := s.db.QueryRow(`
		SELECT status
		FROM parcel
		WHERE number = ?
	`, number)

	var status string
	if err := row.Scan(&status); err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("parcel not found")
		}
		return err
	}

	if status != ParcelStatusRegistered {
		return fmt.Errorf("address can only be updated for registered parcels")
	}

	_, err := s.db.Exec(`
		UPDATE parcel
		SET address = ?
		WHERE number = ?
	`, address, number)
	return err
}

func (s ParcelStore) Delete(number int) error {
	// Проверяем текущий статус посылки
	row := s.db.QueryRow(`
		SELECT status
		FROM parcel
		WHERE number = ?
	`, number)

	var status string
	if err := row.Scan(&status); err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("parcel not found")
		}
		return err
	}

	if status != ParcelStatusRegistered {
		return fmt.Errorf("parcel can only be deleted if status is registered")
	}

	_, err := s.db.Exec(`
		DELETE FROM parcel
		WHERE number = ?
	`, number)
	return err
}
