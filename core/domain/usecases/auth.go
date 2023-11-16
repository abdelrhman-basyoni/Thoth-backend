package domain

import (
	"errors"

	domain "github.com/abdelrhman-basyoni/thoth-backend/core/domain/repositories"
	repos "github.com/abdelrhman-basyoni/thoth-backend/core/implementation/repositories"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthUseCases struct {
	userRepo domain.UserRepository
}

func NewAuthUseCases(db *gorm.DB) *AuthUseCases {
	repo := repos.NewUserRepoSql(db)
	return &AuthUseCases{userRepo: repo}
}

func (auc *AuthUseCases) Login(email, password string) (string, error) {

	candidate := auc.userRepo.GetUserByEmail(email)

	if candidate == nil {
		return "", errors.New("invalid User Credentials")
	}

	err := candidate.ValidatePassword(password)

	if err != nil {
		return "", errors.New("invalid User Credentials")
	}
	return candidate.SignToken()
}

func (auc *AuthUseCases) Create(name, email, password, role string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	err = auc.userRepo.CreateUser(name, email, string(hashedPassword), role)

	if err != nil {
		return errors.New("failed to create user")
	}
	return nil
}
