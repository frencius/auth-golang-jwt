// This file contains the interfaces for the repository layer.
// The repository layer is responsible for interacting with the database.
// For testing purpose we will generate mock implementations of these
// interfaces using mockgen. See the Makefile for more information.
package repository

import (
	"context"
)

type RepositoryInterface interface {
	StoreRegistration(ctx context.Context, data *User) error
	GetUser(ctx context.Context, phoneNumber string) (*User, error)
	GetUserByID(ctx context.Context, userID string) (*User, error)
	UpdateLogin(ctx context.Context, userID string) error
	UpdateProfile(ctx context.Context, data *User) error
}
