package main

import (
	"database/sql"
	"time"
)

func initDB(filepath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", filepath)
	if err != nil {
		return nil, err
	}

	createTableSQL := `CREATE TABLE IF NOT EXISTS files (
        "filename" TEXT,
        "sha" TEXT PRIMARY KEY,
        "date" TEXT,
        "status" TEXT,
        "log" TEXT
    );`

	_, err = db.Exec(createTableSQL)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func fileExistsInDB(db *sql.DB, sha string) (bool, error) {
	query := `SELECT COUNT(*) FROM files WHERE sha = ?`
	var count int
	err := db.QueryRow(query, sha).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func getFilePrevScanDate(db *sql.DB, filePath string) (time.Time, error) {
	query := `SELECT date FROM files WHERE filename = ? ORDER BY date DESC LIMIT 1`
	var date string
	err := db.QueryRow(query, filePath).Scan(&date)
	if err != nil {
		if err == sql.ErrNoRows {
			return time.Time{}, nil // No previous scan date
		}
		return time.Time{}, err
	}
	return time.Parse(time.RFC3339, date)
}

func saveFileInfo(db *sql.DB, filename, sha, status, logMessage string) error {
	insertFileSQL := `INSERT INTO files(filename, sha, date, status, log) VALUES (?, ?, ?, ?, ?)`
	statement, err := db.Prepare(insertFileSQL)
	if err != nil {
		return err
	}
	defer statement.Close()

	_, err = statement.Exec(filename, sha, time.Now().Format(time.RFC3339), status, logMessage)
	if err != nil {
		return err
	}

	return nil
}
