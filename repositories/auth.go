package repositories

import (
	"database/sql"
	"errors"
	"log"

	"github.com/akinolaemmanuel49/notify-api/models"
	"github.com/akinolaemmanuel49/notify-api/utils"
)

type AuthRepository struct {
	db *sql.DB
}

func NewAuthRepository(db *sql.DB) *AuthRepository {
	return &AuthRepository{
		db: db,
	}
}

func (r *AuthRepository) AuthenticateUser(credentials *models.AuthCredentials) (bool, int64, error) {
	query := `
	SELECT 
	id, password_hash
	FROM users
	WHERE email = ($1)`

	result := r.db.QueryRow(query, credentials.Email)

	var passwordHash string
	var ID int64
	err := result.Scan(&ID, &passwordHash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, 0, ErrNotFound
		}
		log.Panicln("Error authenticating user:", err)
		return false, 0, err
	}
	if utils.VerifyPassword(passwordHash, credentials.Password) {
		return true, ID, nil
	}
	return false, 0, ErrInvalidCredentials
}
