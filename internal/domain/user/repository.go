package user

import (
	"database/sql"
)

// DataBase represents the database where the queries for users will be sent to.
type DataBase struct {
	dbase sql.DB
}

// InsertUser inserts the user into users table.
func (db *DataBase) InsertUser(u *User) (err error) {
	const insertUserQuery = `INSERT INTO converter."user" (username, password) VALUES ($1, $2)`
	_, err = db.dbase.Exec(insertUserQuery, u.Username, u.Password)
	if err != nil {
		return err
	}

	return nil
}

// GetPassword retrieves the database hashed password of a user.
func (db *DataBase) GetPassword(u *User) (string, error) {
	var password string
	const getPasswordByUsername = `SELECT password FROM converter."user" WHERE username=$1;`
	err := db.dbase.QueryRow(getPasswordByUsername, u.Username).Scan(&password)
	return password, err
}
