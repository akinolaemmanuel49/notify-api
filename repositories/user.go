package repositories

import (
	"database/sql"
	"errors"
	"log"
	"strconv"
	"time"

	"github.com/akinolaemmanuel49/notify-api/models"
	"github.com/akinolaemmanuel49/notify-api/utils"
	"github.com/lib/pq"
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
func (r *UserRepository) CreateUser(userInput *models.UserInput, password string) error {
	currentTime := time.Now().UTC().Format(time.RFC3339)
	hashedPassword, err := utils.GenerateHashPassword(password)
	if err != nil {
		log.Println("An error occured while hashing password:", err)
	}

	user := models.User{
		FirstName:    userInput.FirstName,
		LastName:     userInput.LastName,
		Email:        userInput.Email,
		PasswordHash: hashedPassword,
		CreatedAt:    currentTime,
		UpdatedAt:    currentTime,
	}
	query := `
	INSERT INTO users(
		first_name,
		last_name,
		email,
		password_hash,
		created_at,
		updated_at
	) VALUES (($1), ($2), ($3), ($4), ($5), ($6))`

	_, err = r.db.Exec(query, user.FirstName, user.LastName, user.Email, user.PasswordHash, user.CreatedAt, user.UpdatedAt)
	if err != nil {
		if pgErr, ok := err.(*pq.Error); ok {
			if pgErr.Code == "23505" {
				return utils.ErrDuplicateKey
			}
		}
		log.Println("Error inserting user:", err)
	}
	return err
}

func (r *UserRepository) GetUserByID(id int64) (*models.UserProfile, error) {
	query := `
	SELECT 
	id, first_name, last_name, email, created_at, updated_at
	FROM users
	WHERE id = ($1)`

	result := r.db.QueryRow(query, id)

	var userProfile models.UserProfile
	err := result.Scan(&userProfile.ID, &userProfile.FirstName, &userProfile.LastName, &userProfile.Email, &userProfile.CreatedAt, &userProfile.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Println("Error retrieving user:", err)
			return nil, utils.ErrNotFound
		}
		log.Println("Error retrieving user:", err)
		return nil, err
	}
	return &userProfile, nil
}

func (r *UserRepository) GetAllUsers(page, pageSize int) ([]*models.UserProfile, error) {
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * pageSize
	query := `
	SELECT 
	id, first_name, last_name, email, created_at, updated_at
	FROM users
	LIMIT $1
	OFFSET $2`
	results, err := r.db.Query(query, pageSize, offset)
	if err != nil {
		log.Println("Error retrieving users:", err)
		return nil, err
	}
	defer results.Close()

	userProfiles := []*models.UserProfile{}
	for results.Next() {
		var userProfile models.UserProfile
		err := results.Scan(&userProfile.ID, &userProfile.FirstName, &userProfile.LastName, &userProfile.Email, &userProfile.CreatedAt, &userProfile.UpdatedAt)
		if err != nil {
			log.Println("Error scanning user row:", err)
			return nil, err
		}
		userProfiles = append(userProfiles, &userProfile)
	}
	if err := results.Err(); err != nil {
		log.Println("Error iterating over user rows:", err)
		return nil, err
	}
	return userProfiles, nil
}

func (r *UserRepository) UpdateUserByID(id int64, fields map[string]interface{}) error {
	_, err := r.GetUserByID(id)
	if errors.Is(err, utils.ErrNotFound) {
		log.Println("Error updating user:", err)
		return utils.ErrNotFound
	}

	query := "UPDATE users SET "

	var params []interface{}
	i := 1
	for key, value := range fields {
		if key == "password" || key == "created_at" || key == "updated_at" {
			continue
		}

		if i > 1 {
			query += ", "
		}
		query += key + " = $" + strconv.Itoa(i)
		params = append(params, value)
		i++
	}

	query += ", updated_at = $" + strconv.Itoa(i)
	query += " WHERE id = $" + strconv.Itoa(i+1)

	updatedAt := time.Now().UTC().Format(time.RFC3339)
	params = append(params, updatedAt, id)

	_, err = r.db.Exec(query, params...)
	if err != nil {
		log.Println("Error updating user: ", err)
	}
	return err
}

func (r *UserRepository) DeleteUserByID(id int64) error {
	_, err := r.GetUserByID(id)
	if errors.Is(err, utils.ErrNotFound) {
		return utils.ErrNotFound
	}

	query := `
	DELETE FROM users WHERE id = ($1)`

	_, err = r.db.Exec(query, id)

	if err != nil {
		log.Println("Error deleting user: ", err)
	}
	return err
}
