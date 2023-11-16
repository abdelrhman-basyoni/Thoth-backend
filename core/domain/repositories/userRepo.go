package domain

import "github.com/abdelrhman-basyoni/thoth-backend/core/domain/entities"

type UserRepository interface {
	GetUserByEmail(email string) *entities.User
	CreateUser(name, email, password, role string) error
}
