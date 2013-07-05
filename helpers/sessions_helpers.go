package helpers

import (
	"models"
)

func UserAuthorized(u *models.GoogleUser) bool {
	return u.Email == "remy.jourde@gmail.com" || u.Email == "santiago.ariassar@gmail.com"
}

func LoggedIn() bool {
	return models.CurrentUser != nil
}

func Logout() {
	models.CurrentUser = nil
}