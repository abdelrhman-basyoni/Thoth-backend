// utils/database.go
package config

import (
	"fmt"

	"github.com/abdelrhman-basyoni/thoth-backend/core/implementation/models"
	"github.com/abdelrhman-basyoni/thoth-backend/utils"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDB() *gorm.DB {
	fmt.Println("Connecting to database...")
	uri := utils.ReadEnv("SQL_URI")
	db, err := gorm.Open(postgres.Open(uri), &gorm.Config{})
	// Migrate the User model
	db.AutoMigrate(&models.User{})
	db.AutoMigrate(&models.Blog{})
	db.AutoMigrate(&models.Comment{})

	if err != nil {
		panic("Failed to connect to database")
	}
	fmt.Println("Connected to database")
	return db
}
