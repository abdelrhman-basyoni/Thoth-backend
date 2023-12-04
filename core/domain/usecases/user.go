package domain

import (
	"github.com/abdelrhman-basyoni/thoth-backend/core/domain/entities"
	domain "github.com/abdelrhman-basyoni/thoth-backend/core/domain/repositories"
	repos "github.com/abdelrhman-basyoni/thoth-backend/core/implementation/repositories"
	typ "github.com/abdelrhman-basyoni/thoth-backend/types"
	"gorm.io/gorm"
)

type UserUseCases struct {
	userRepo domain.UserRepository
}

func NewUserUseCases(db *gorm.DB) *UserUseCases {
	repo := repos.NewUserRepoSql(db)
	return &UserUseCases{userRepo: repo}
}

func (uuc *UserUseCases) GetAllUsers() (*typ.PaginatedEntities[entities.GetAllUsersRes], error) {
	return uuc.userRepo.GetAllUsers()
}
