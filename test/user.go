package gonawintest

import (
	"appengine/aetest"

	mdl "github.com/taironas/gonawin/models"
)

// CreateUsers creates n users into the datastore
func CreateUsers(c aetest.Context, count int64) []*mdl.User {
	var users []*mdl.User

	var i int64
	for i = 0; i < count; i++ {
		firstName, lastName := generateNames()
		username := generateUsername(firstName, lastName)
		email := generateEmail(firstName, lastName)

		if user, _ := mdl.CreateUser(c, email, username, firstName+" "+lastName, "", false, ""); user != nil {
			users = append(users, user)
		}
	}

	return users
}
