package repositories

import (
	"database/sql"
	"log"

	"github.com/akinolaemmanuel49/notify-api/models"
	"github.com/akinolaemmanuel49/notify-api/utils"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) CreateUser(user *models.User) error {
	query := `
	INSERT INTO users(
		first_name,
		last_name,
		email,
		password_hash,
		created_at,
		updated_at
	) VALUES (($1), ($2), ($3), ($4), ($5), ($6))`

	hashedPassword, err := utils.GenerateHashPassword(user.PasswordHash)

	if err != nil {
		log.Panicln("An error occured while hashing password:", err)
	}

	user.PasswordHash = hashedPassword

	_, err = r.db.Exec(query, user.FirstName, user.LastName, user.Email, user.PasswordHash, user.CreatedAt, user.UpdatedAt)
	if err != nil {
		log.Panicln("Error inserting user:", err)
	}
	return err
}
