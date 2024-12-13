package db

import (
	"database/sql"
	"errors"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

type DB struct {
	conn *sql.DB
}

func NewDB() (*DB, error) {
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbHost := os.Getenv("DB_HOST")
	dbName := os.Getenv("DB_NAME")

	if dbUser == "" || dbPass == "" || dbHost == "" || dbName == "" {
		return nil, fmt.Errorf("database environment variables not set. Required: DB_USER, DB_PASS, DB_HOST, DB_NAME")
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true&charset=utf8mb4&loc=Local", dbUser, dbPass, dbHost, dbName)
	conn, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("error opening DB: %w", err)
	}
	if err := conn.Ping(); err != nil {
		return nil, fmt.Errorf("error pinging DB: %w", err)
	}

	return &DB{conn: conn}, nil
}

func (db *DB) FindCheatByHash(fileHash string) (*Cheat, error) {
	var cheatID int
	err := db.conn.QueryRow("SELECT cheat_id FROM cheat_hashes WHERE file_hash = ?", fileHash).Scan(&cheatID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // no match
		}
		return nil, err
	}
	return db.FindCheatByID(cheatID)
}

func (db *DB) FindCheatByID(cheatID int) (*Cheat, error) {
	var c Cheat
	err := db.conn.QueryRow("SELECT cheat_id, cheat_name, cheat_description, date_added FROM cheats WHERE cheat_id = ?", cheatID).
		Scan(&c.CheatID, &c.Name, &c.Desc, &c.DateAdded)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &c, nil
}

func (db *DB) LoadAllSignatures() ([]Signature, error) {
	rows, err := db.conn.Query("SELECT signature_id, cheat_id, signature_pattern, description, date_added FROM signatures")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sigs []Signature
	for rows.Next() {
		var s Signature
		if err := rows.Scan(&s.SignatureID, &s.CheatID, &s.SignaturePattern, &s.Description, &s.DateAdded); err != nil {
			return nil, err
		}
		sigs = append(sigs, s)
	}
	return sigs, nil
}

func (db *DB) InsertCheatHash(cheatID int, hash, desc string) error {
	_, err := db.conn.Exec("INSERT INTO cheat_hashes (cheat_id, file_hash, description) VALUES (?, ?, ?)", cheatID, hash, desc)
	return err
}
