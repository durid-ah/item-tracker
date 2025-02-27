package services

import (
	"database/sql"
	"log"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Id       string `json:"id"`
	Username string `json:"username"`
	Password []byte `json:"password"`
}

type UserService struct {
	Db *sql.DB
}

func (service *UserService) Add(newUser *User) (int64, error) {
	// convert password to hash and store it in the Password field
	bytes, err := bcrypt.GenerateFromPassword(newUser.Password, 14)
	if err != nil {
		log.Fatalln("Failed to hash password")
	}
	newUser.Password = bytes


	stmt, stmtErr := service.Db.Prepare(`INSERT INTO users (username, password) VALUES (?,?)`)
	if stmtErr != nil {
		return 0, stmtErr
	}

	defer stmt.Close()

	res, err := stmt.Exec(newUser.Username, newUser.Password)
	if err != nil {
		return 0, err
	}

	id, _ := res.LastInsertId()
	return id, nil
}

func (service *UserService) GetByUsername(username string) (*User, error) {
	row := service.Db.QueryRow(`
		SELECT
			id,
			username,
			password
		FROM users
		WHERE username = ?`, username)
	
	user := new(User)
	rowErr := row.Scan(&user.Id, &user.Username, &user.Password)
	if rowErr != nil {
		return nil, rowErr
	}

	return user, nil
}