package repos

import (
	"sync"

	"github.com/abdelrhman-basyoni/thoth-backend/core/domain/entities"
	"github.com/abdelrhman-basyoni/thoth-backend/core/implementation/models"
	typ "github.com/abdelrhman-basyoni/thoth-backend/types"
	"gorm.io/gorm"
)

type UserRepoSql struct {
	db *gorm.DB
}

func NewUserRepoSql(db *gorm.DB) *UserRepoSql {
	return &UserRepoSql{db: db}
}

func (ur *UserRepoSql) GetUserByEmail(email string) *entities.User {

	var user entities.User
	res := ur.db.First(&user, "email = ?", email)

	// Check if a record was found
	if res.RowsAffected == 0 {
		return nil
	}

	return &user

}

func (ur *UserRepoSql) CreateUser(name, email, password, role string) error {

	res := ur.db.Create(&models.User{Name: name, Email: email, Password: password, Role: role})

	return res.Error

}

func (ur *UserRepoSql) GetAllUsers() (*typ.PaginatedEntities[entities.GetAllUsersRes], error) {
	res := typ.PaginatedEntities[entities.GetAllUsersRes]{}
	var wg sync.WaitGroup
	wg.Add(2)
	var totalErr error

	go func() {
		defer wg.Done()
		if err := ur.db.Model(&models.User{}).Count(&res.Total).Error; err != nil {
			totalErr = err
		}
	}()

	go func() {
		defer wg.Done()
		if err := ur.db.Model(&models.User{}).Find(&res.Entities).Error; err != nil {
			totalErr = err
		}
	}()

	wg.Wait()

	if totalErr != nil {
		return nil, totalErr
	}

	return &res, nil

}
