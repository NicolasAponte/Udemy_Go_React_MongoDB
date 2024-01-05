package bd

import (
	"github.com/naponte/Udemy_Go_React_MongoDB/models"
	"golang.org/x/crypto/bcrypt"
)

func Login(email string, password string) (models.User, bool) {
	user, exists, _ := UserAlreadyExists(email)
	if !exists {
		return user, false
	}

	passBytes := []byte(password)
	passDB := []byte(user.Password)

	err := bcrypt.CompareHashAndPassword(passDB, passBytes)
	if err != nil {
		return user, false
	}

	return user, true
}
