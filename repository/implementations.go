package repository

import (
	"context"
	"fmt"

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

func (r *Repository) GetUser(ctx context.Context, phone_number string) (*User, error) {
	user := &User{}

	query := fmt.Sprintf(`
	SELECT
		id, phone_number, password, full_name
	FROM 
		"user"
	WHERE
		phone_number = $1`)

	err := r.Db.QueryRow(query, phone_number).Scan(&user.ID, &user.PhoneNumber, &user.Password, &user.FullName)
	if err != nil {
		return nil, err
	}

	return user, err
}

func (r *Repository) UpdateLogin(ctx context.Context, userID string) error {
	loginID := uuid.New().String()
	query := fmt.Sprintf(`
	INSERT INTO login (id, user_id, success_counter) VALUES ($1, $2, 1) 
	ON CONFLICT (user_id) 
	DO UPDATE SET success_counter = (SELECT success_counter FROM login WHERE user_id =$2) + 1;
	`)

	_, err := r.Db.Exec(query, loginID, userID)
	if err != nil {
		return err
	}

	return err
}
