package domain

import (
	"github.com/abdelrhman-basyoni/thoth-backend/core/domain/entities"
	typ "github.com/abdelrhman-basyoni/thoth-backend/types"
)

type UserRepository interface {
	GetUserByEmail(email string) *entities.User
	CreateUser(name, email, password, role string) error
	GetAllUsers() (*typ.PaginatedEntities[entities.GetAllUsersRes], error)
}
