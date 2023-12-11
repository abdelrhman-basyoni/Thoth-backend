package entities

import (
	"time"

	"github.com/abdelrhman-basyoni/thoth-backend/utils"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       uint   `json:"id"` // defining id as string so it can work with any database not just sql types
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

func (u *User) ValidatePassword(candidatePassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(candidatePassword))
}

func (u *User) SignToken() (string, error) {
	secretKey := utils.ReadEnv("SECRET_KEY")
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)

	claims["id"] = u.ID
	claims["role"] = u.Role

	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

	return token.SignedString([]byte(secretKey))

}

type GetAllUsersRes struct {
	ID        string `json:"id"` // defining id as string so it can work with any database not just sql types
	Name      string `json:"name"`
	Email     string `json:"email"`
	Role      string `json:"role"`
	CreatedAt string `json:"created_at"`
}
