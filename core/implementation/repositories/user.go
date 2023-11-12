package repos

import (
	"github.com/abdelrhman-basyoni/thoth-backend/core/domain/entities"
	"github.com/abdelrhman-basyoni/thoth-backend/core/implementation/models"
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
