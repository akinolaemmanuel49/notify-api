package repositories

import (
	"database/sql"
	"errors"
	"log"
	"strconv"

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

// CreateUser creates a new instance of UserRepository
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

func (r *UserRepository) GetUserByID(id int64) (*models.User, error) {
	query := `
	SELECT 
	id, first_name, last_name, email, password_hash, created_at, updated_at
	FROM users
	WHERE id = ($1)`

	result := r.db.QueryRow(query, id)

	var user models.User
	err := result.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.PasswordHash, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		log.Panicln("Error retrieving user:", err)
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetAllUsers(page, pageSize int) ([]*models.User, error) {
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * pageSize
	query := `
	SELECT 
	id, first_name, last_name, email, password_hash, created_at, updated_at
	FROM users
	LIMIT $1
	OFFSET $2`
	results, err := r.db.Query(query, pageSize, offset)
	if err != nil {
		log.Panicln("Error retrieving users:", err)
		return nil, err
	}
	defer results.Close()

	users := []*models.User{}
	for results.Next() {
		var user models.User
		err := results.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.PasswordHash, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			log.Panicln("Error scanning user row:", err)
			return nil, err
		}
		users = append(users, &user)
	}
	if err := results.Err(); err != nil {
		log.Panicln("Error iterating over user rows:", err)
		return nil, err
	}
	return users, nil
}

// TODO: Exclude password
func (r *UserRepository) UpdateUserByID(id int64, fields map[string]interface{}) error {
	_, err := r.GetUserByID(id)
	if errors.Is(err, ErrNotFound) {
		return ErrNotFound
	}

	query := "UPDATE users SET "

	var params []interface{}
	i := 1
	for key, value := range fields {
		if key == "password" {
			continue
		}
		
		if i > 1 {
			query += ", "
		}
		query += key + " = $" + strconv.Itoa(i)
		params = append(params, value)
		i++
	}

	query += " WHERE id = $" + strconv.Itoa(i)
	params = append(params, id)

	_, err = r.db.Exec(query, params...)
	if err != nil {
		log.Panicln("Error updating notification: ", err)
	}
	return err
}

func (r *UserRepository) DeleteUserByID(id int64) error {
	_, err := r.GetUserByID(id)
	if errors.Is(err, ErrNotFound) {
		return ErrNotFound
	}

	query := `
	DELETE FROM notifications WHERE id = ($1)`

	_, err = r.db.Exec(query, id)

	if err != nil {
		log.Panicln("Error deleting notification: ", err)
	}
	return err
}
