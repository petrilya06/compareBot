package db

import (
	"database/sql"
	"fmt"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"log"
)

type User struct {
	TgID              int64
	Confirm           bool
	SelectPic         int
	SelectText        int
	CountPhotoCompare int
	CountTextCompare  int
	LastPhotoID       int
	LastMessageID     int
}

type Database struct {
	Conn *sql.DB
}

func NewDatabase() (*Database, error) {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	psqlInfo := "user=postgres password=qwerty123456 host=localhost sslmode=disable"

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return &Database{Conn: db}, nil
}

func (d *Database) CreateTable() error {
	createTableSQL := `CREATE TABLE IF NOT EXISTS users (
		tgID BIGINT PRIMARY KEY,
		confirm BOOLEAN, 
		selectPic INT,
		selectText INT,
		countPhotoCompare INT,
		countTextCompare INT,
		lastPhotoID INT,
		lastMessageID INT
	);`
	_, err := d.Conn.Exec(createTableSQL)
	if err != nil {
		log.Fatal(err)
	}

	return err
}

func (d *Database) InsertUser(user User) error {
	insertSQL := `INSERT INTO users (tgID, confirm, selectPic, selectText, countPhotoCompare, 
                   countTextCompare, lastPhotoID, lastMessageID) 
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	ON CONFLICT (tgID) DO UPDATE SET 
		selectPic = EXCLUDED.selectPic,
	    selectText = EXCLUDED.selectText,
		countPhotoCompare = EXCLUDED.countPhotoCompare,
	    countTextCompare = EXCLUDED.countTextCompare,
		lastPhotoID = EXCLUDED.lastPhotoID,
		lastMessageID = EXCLUDED.lastMessageID;`

	_, err := d.Conn.Exec(insertSQL, user.TgID, user.Confirm, user.SelectPic, user.SelectText,
		user.CountPhotoCompare, user.CountTextCompare, user.LastPhotoID, user.LastMessageID)
	if err != nil {
		log.Fatal(err)
	}

	return err
}

func (d *Database) UpdateUser(user User) error {
	updateSQL := `
		UPDATE users 
		SET confirm = $1, selectPic = $2, selectText = $3, countPhotoCompare = $4, 
		    countTextCompare = $5, lastPhotoID = $6, lastMessageID = $7
		WHERE tgID = $8`
	_, err := d.Conn.Exec(updateSQL, &user.Confirm, &user.SelectPic, &user.SelectText, &user.CountPhotoCompare,
		&user.CountTextCompare, &user.LastPhotoID, &user.LastMessageID, &user.TgID)
	if err != nil {
		return fmt.Errorf("error in update: %w", err)
	}

	return nil
}

//
//func (d *Database) GetDataUser(tgID int64) (User, error) {
//	var user User
//
//	query := `SELECT * FROM users WHERE tgID = $1;`
//	row := d.Conn.QueryRow(query, tgID)
//
//	err := row.Scan(&user.TgID, &user.Confirm, &user.SelectPic, &user.SelectText,
//		&user.CountPhotoCompare, &user.CountTextCompare, &user.LastPhotoID, &user.LastMessageID)
//	if err != nil {
//		if err == sql.ErrNoRows {
//			return nil, fmt.Errorf("no user found with tgID: %d", tgID)
//		}
//		return nil, fmt.Errorf("error retrieving user: %w", err)
//	}
//
//	return &user, nil
//}

func (d *Database) CloseDatabase() {
	if err := d.Conn.Close(); err != nil {
		fmt.Printf("error closing database: %v\n", err)
	}
}
