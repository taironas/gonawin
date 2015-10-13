package gonawintest

import (
	"fmt"
	"math/rand"
	"strings"

	"appengine/aetest"

	mdl "github.com/taironas/gonawin/models"
)

var firstNames = [...]string{"Jon", "Robb", "Cersei", "Daenerys", "Ned", "Tyrion", "Stannis"}
var lastNames = [...]string{"Stark", "Targaryen", "Lannister", "Baratheon"}

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

func CreateWorldCup(c aetest.Context, adminID int64) *mdl.Tournament {
	tournament, _ := mdl.CreateWorldCup(c, adminID)

	return tournament
}

func PredictOnMatch(c aetest.Context, user *mdl.User, matchID, result1, result2 int64) error {
	if predict, err := mdl.CreatePredict(c, user.Id, result1, result2, matchID); err != nil {
		return err
	} else {
		// add p.Id to User predict table.
		if err = user.AddPredictId(c, predict.Id); err != nil {
			return err
		}
	}

	return nil
}

func generateNames() (firstName string, lastName string) {
	firstName = firstNames[rand.Int()%len(firstNames)]
	lastName = lastNames[rand.Int()%len(lastNames)]

	return firstName, lastName
}

func generateUsername(firstName, lastName string) string {
	return fmt.Sprintf("%c%s", strings.ToLower(firstName)[0], strings.ToLower(lastName))
}

func generateEmail(firstName, lastName string) string {
	return fmt.Sprintf("%s.%s@westeros.com", firstName, lastName)
}
