package services

import (
	"database/sql"
	"errors"
	"github.com/code-chimp/snippetbox/internal/models"
	"github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
	"strings"
)

// UsersService wraps a sql.DB connection pool.
type UsersService struct {
	DB *sql.DB
}

// Insert adds a new user to the database with the provided name, email, and password.
// It returns an error if the insertion fails or if the email is already in use.
func (m *UsersService) Insert(name, email, password string) error {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	query := `INSERT INTO users
	          (name, email, hashed_password, created)
		      VALUES
		      (?, ?, ?, UTC_TIMESTAMP())`

	_, err = m.DB.Exec(query, name, email, string(hashed))
	if err != nil {
		var dbError *mysql.MySQLError

		if errors.As(err, &dbError) {
			if dbError.Number == 1062 && strings.Contains(dbError.Message, "users_uc_email") {
				return models.ErrDuplicateEmail
			}
		}

		return err
	}

	return nil
}

// Authenticate verifies a user's email and password.
// It returns the user's ID if authentication is successful, or an error if it fails.
func (m *UsersService) Authenticate(email, password string) (int, error) {
	var id int
	var hashedPassword []byte

	err := m.DB.QueryRow(
		`SELECT id, hashed_password FROM users WHERE email = ?`,
		email).Scan(&id, &hashedPassword)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, models.ErrInvalidCredentials
		}

		return 0, err
	}

	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return 0, models.ErrInvalidCredentials
		}

		return 0, err
	}

	return id, nil
}

// Exists checks if a user with the given ID exists in the database.
// It returns true if the user exists, or false otherwise.
func (m *UsersService) Exists(id int) (bool, error) {
	var exists bool

	err := m.DB.QueryRow(
		`SELECT EXISTS(SELECT true FROM users WHERE id = ?)`,
		id).Scan(&exists)

	return exists, err
}
