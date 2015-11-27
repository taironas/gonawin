package gonawintest

import (
	"errors"
	"fmt"
	"time"

	mdl "github.com/taironas/gonawin/models"
)

// TestTournament represents a testing tournament
type TestTournament struct {
	Name        string
	Description string
	Start       time.Time
	End         time.Time
	AdminID     int64
	UserIDs     []int64
}

// CheckTournament checks that a given tournament is equivalent to a test tournament
func CheckTournament(got *mdl.Tournament, want TestTournament) error {
	var s string
	if got.Name != want.Name {
		s = fmt.Sprintf("want Name == %s, got %s", want.Name, got.Name)
	} else if got.Description != want.Description {
		s = fmt.Sprintf("want Description == %s, got %s", want.Description, got.Description)
	} else if got.Start.Sub(want.Start).Hours() > 0 {
		s = fmt.Sprintf("want Start == %v, got %v", want.Start, got.Start)
	} else if got.End.Sub(want.End).Hours() > 0 {
		s = fmt.Sprintf("want End == %v, got %v", want.End, got.End)
	} else if got.AdminIds[0] != want.AdminID {
		s = fmt.Sprintf("want AdminId == %d, got %d", want.AdminID, got.AdminIds[0])
	} else if len(got.UserIds) != len(want.UserIDs) {
		s = fmt.Sprintf("want UserIds count == %d, got %d", len(want.UserIDs), len(got.UserIds))
	} else {
		return nil
	}

	return errors.New(s)
}
