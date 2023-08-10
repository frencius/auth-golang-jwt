package repository

import (
	"context"
)

func (r *Repository) StoreRegistration(ctx context.Context, data *Registration) error {
	query := `
	INSERT INTO "user" (id, full_name, phone_number, password)
	VALUES ($1, $2, $3, $4);
`
	_, err := r.Db.Exec(query, data.ID, data.FullName, data.PhoneNumber, data.Password)
	if err != nil {
		return err
	}

	return nil
}
