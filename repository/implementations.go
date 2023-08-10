package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

func (r *Repository) StoreRegistration(ctx context.Context, data *User) error {
	query := `
	INSERT INTO "user" (id, full_name, phone_number, password)
	VALUES ($1, $2, $3, $4);
`
	_, err := r.Db.Exec(query, data.ID, data.FullName, data.PhoneNumber, data.Password)
	if err != nil {
		return err
	}

	return err
}

func (r *Repository) GetUser(ctx context.Context, phoneNumber string) (*User, error) {
	user := &User{}

	query := `
	SELECT
		id, phone_number, password, full_name
	FROM 
		"user"
	WHERE
		phone_number = $1`

	err := r.Db.QueryRow(query, phoneNumber).Scan(&user.ID, &user.PhoneNumber, &user.Password, &user.FullName)
	if err != nil {
		return nil, err
	}

	return user, err
}

func (r *Repository) GetUserByID(ctx context.Context, userID string) (*User, error) {
	user := &User{}

	query := `
	SELECT
		id, phone_number, password, full_name
	FROM 
		"user"
	WHERE
		id = $1`

	err := r.Db.QueryRow(query, userID).Scan(&user.ID, &user.PhoneNumber, &user.Password, &user.FullName)
	if err != nil {
		return nil, err
	}

	return user, err
}

func (r *Repository) UpdateLogin(ctx context.Context, userID string) error {
	loginID := uuid.New().String()
	query := `
	INSERT INTO login (id, user_id, success_counter) VALUES ($1, $2, 1) 
	ON CONFLICT (user_id) 
	DO UPDATE SET success_counter = (SELECT success_counter FROM login WHERE user_id =$2) + 1;
	`

	_, err := r.Db.Exec(query, loginID, userID)
	if err != nil {
		return err
	}

	return err
}

func (r *Repository) UpdateProfile(ctx context.Context, data *User) error {
	query := `
	UPDATE
		"user"
	SET
		full_name = $2,
		phone_number = $3
	WHERE
		id = $1
	`

	result, err := r.Db.Exec(query, data.ID, data.FullName, data.PhoneNumber)
	if err != nil {
		return err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if affected < 1 {
		err = errors.New("user is not exist")
		return err
	}

	return nil
}
