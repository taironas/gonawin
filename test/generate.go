package gonawintest

import (
	"fmt"
	"math/rand"
	"strings"
)

var firstNames = [...]string{"Jon", "Robb", "Cersei", "Daenerys", "Ned", "Tyrion", "Stannis"}
var lastNames = [...]string{"Stark", "Targaryen", "Lannister", "Baratheon"}

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
