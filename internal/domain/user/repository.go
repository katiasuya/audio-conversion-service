package user

import (
	"database/sql"
	"log"

	"github.com/katiasuya/audio-conversion-service/pkg/hash"
)

// DataBase represents the database where the queries for users will be sent to.
type DataBase struct {
	dbase sql.DB
}

// InsertUser inserts the user into users table.
func (db *DataBase) InsertUser(u *User) (err error) {
	u.Password, err = hash.HashPassword(u.Password)
	if err != nil {
		return err
	}

	_, err = db.dbase.Exec(`INSERT INTO converter."user" (username, password) VALUES ($1, $2)`, u.Username, u.Password)
	if err != nil {
		return err
	}

	u.Password = ""
	return nil
}

// ComparePasswords compares entered password with a database hashed password of a user.
func (db *DataBase) ComparePasswords(u *User) bool {
	var password string
	row := db.dbase.QueryRow(`SELECT password FROM converter."user" WHERE username=$1;`, u.Username)
	switch err := row.Scan(&password); err {
	case sql.ErrNoRows:
		return false
	case nil:
		return hash.CheckPasswordHash(u.Password, password)
	default:
		log.Fatal(err)
	}

	return false
}
