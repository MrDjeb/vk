package database

import (
	"database/sql"
	"errors"
	"log"
	"os"

	"github.com/mattn/go-sqlite3"
)

var (
	ErrDuplicate    = errors.New("record already exists")
	ErrNotExists    = errors.New("row not exists")
	ErrUpdateFailed = errors.New("update failed")
	ErrDeleteFailed = errors.New("delete failed")
	ErrCountDate    = errors.New("count of date in table scorer is bad")
)

type DBPhotos struct{ DB *sql.DB }

func (r *DBPhotos) Migrate() error {

	query := `
    CREATE TABLE IF NOT EXISTS photos(
        photo_id INTEGER
    );`
	logDB.Println(query)

	_, err := r.DB.Exec(query)
	return err
}

func (r *DBPhotos) Insert(ph Photo) error {
	logDB.Println("INSERT INTO photos(photo_id) values(?)",
		ph.PhotoID)

	_, err := r.DB.Exec("INSERT INTO photos(photo_id) values(?)",
		ph.PhotoID)

	if err != nil {
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) {
			if errors.Is(sqliteErr.ExtendedCode, sqlite3.ErrConstraintUnique) {
				return ErrDuplicate
			}
		}
		return err
	}

	return nil
}

func (r *DBPhotos) Delete(photoID int) error {
	logDB.Println("DELETE FROM photos WHERE photo_id = ?", photoID)
	_, err := r.DB.Exec("DELETE FROM photos WHERE photo_id = ?", photoID)
	return err
}

func (r *DBPhotos) Read(photoID int) ([]Photo, error) {
	logDB.Println("SELECT * FROM photos WHERE photo_id = ?;", photoID)

	rows, err := r.DB.Query("SELECT * FROM photos WHERE photo_id = ?;", photoID)
	if err != nil {
		return []Photo{}, err
	}
	defer rows.Close()
	var photos []Photo
	for rows.Next() {
		var p Photo
		if err := rows.Scan(&p.PhotoID); err != nil {
			return []Photo{}, err
		}
		photos = append(photos, p)
	}
	return photos, nil
}

func (r *DBPhotos) ReadAll() ([]Photo, error) {
	logDB.Println("SELECT * FROM photos")

	rows, err := r.DB.Query("SELECT * FROM photos")
	if err != nil {
		return []Photo{}, err
	}
	defer rows.Close()
	var photos []Photo
	for rows.Next() {
		var p Photo
		if err := rows.Scan(&p.PhotoID); err != nil {
			return []Photo{}, err
		}
		photos = append(photos, p)
	}
	return photos, nil
}

var logDB *log.Logger

func Init() (Tables, error) {
	logDB = log.New(os.Stderr, "[SQLITE] ", log.LstdFlags|log.Lmsgprefix)

	err := os.MkdirAll("./sqlite", os.ModePerm)
	if err != nil {
		return Tables{}, err
	}
	db, err := sql.Open("sqlite3", "./sqlite/stage.db")
	if err != nil {
		return Tables{}, err
	}
	tables := Tables{DBPhotos{DB: db}}
	tables.Photos.Migrate()
	return tables, err
}
