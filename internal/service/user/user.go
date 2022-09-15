package user

import (
	"context"

	"github.com/danieeelfr/swd-challenge/internal/config"
	"github.com/danieeelfr/swd-challenge/internal/models"
	repo "github.com/danieeelfr/swd-challenge/internal/repository"
)

// User holds the user implementation
type User struct {
	repository *repo.MySQLRepo
}

// NewUserService returns the user service or error
func NewUserService(cfg *config.Config) (*User, error) {

	r := repo.NewMySQLRepo(cfg.MySQLRepositoryConfig)
	err := r.Connect()
	if err != nil {
		return nil, err
	}

	return &User{repository: r}, nil
}

// GetUser get the user to valida on login
func (u *User) GetUser(ctx context.Context, filter *models.UserFilter) (*models.User, error) {

	user := new(models.User)

	u.repository.DB.Table("users")

	if err := u.repository.DB.
		Where("password = ?", filter.Password).
		Where("user = ?", filter.User).
		Find(&user).Error; err != nil {
		return nil, err
	}

	return user, nil
}
