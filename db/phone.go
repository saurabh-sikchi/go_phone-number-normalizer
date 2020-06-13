package db

import (
	"database/sql"

	_ "github.com/lib/pq"
)

// Phone represents phone_numbers table record in the database
type Phone struct {
	ID     int
	Number string
}

func Open(driverName, dataSource string) (*DB, error) {
	db, err := sql.Open(driverName, dataSource)
	if err != nil {
		return nil, err
	}
	return &DB{db}, nil
}

type DB struct {
	db *sql.DB
}

func (db *DB) Close() error {
	return db.db.Close()
}

func (db *DB) Seed() error {
	for _, phone := range []string{
		"1234567890",
		"123 456 7891",
		"(123) 456 7892",
		"(123) 456-7893",
		"123-456-7894",
		"123-456-7890",
		"1234567892",
		"(123)456-7892",
	} {
		if _, err := insertPhone(db.db, phone); err != nil {
			return err
		}
	}
	return nil
}

func insertPhone(db *sql.DB, phone string) (int, error) {
	statement := `INSERT INTO phone_numbers(value) VALUES($1) RETURNING id`
	var id int
	err := db.QueryRow(statement, phone).Scan(&id)
	if err != nil {
		return -1, err
	}
	return id, nil
}

func Reset(driverName, dataSource, dbName string) error {
	db, err := sql.Open(driverName, dataSource)
	if err != nil {
		return err
	}
	err = resetDB(db, dbName)
	if err != nil {
		return err
	}
	return db.Close()
}

func Migrate(driverName, dataSource string) error {
	db, err := sql.Open(driverName, dataSource)
	if err != nil {
		return err
	}
	err = createPhoneNumbersTable(db)
	if err != nil {
		return nil
	}
	return db.Close()
}

func createPhoneNumbersTable(db *sql.DB) error {
	statement := `
		CREATE TABLE IF NOT EXISTS phone_numbers (
			id SERIAL,
			value VARCHAR(255)
		)
	`
	_, err := db.Exec(statement)
	return err
}

func resetDB(db *sql.DB, name string) error {
	_, err := db.Exec("DROP DATABASE IF EXISTS " + name)
	if err != nil {
		return nil
	}
	return createDB(db, name)
}

func createDB(db *sql.DB, name string) error {
	_, err := db.Exec("CREATE DATABASE " + name)
	return err
}

func (db *DB) AllPhones() ([]Phone, error) {
	rows, err := db.db.Query("SELECT id, value FROM phone_numbers")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var phones []Phone
	for rows.Next() {
		var p string
		var id int
		if err := rows.Scan(&id, &p); err != nil {
			return nil, err
		}
		phones = append(phones, Phone{id, p})
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return phones, nil
}

func (db *DB) FindPhone(number string) (*Phone, error) {
	var p Phone

	err := db.db.QueryRow("SELECT id, value FROM phone_numbers WHERE value = $1", number).Scan(&p.ID, &p.Number)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	return &p, nil
}

func (db *DB) DeletePhone(id int) error {
	statement := `DELETE FROM phone_numbers WHERE id = $1`
	_, err := db.db.Exec(statement, id)
	return err
}

func (db *DB) UpdatePhone(p *Phone) error {
	statement := `UPDATE phone_numbers SET value = $2 WHERE id = $1`
	_, err := db.db.Exec(statement, p.ID, p.Number)
	return err
}
