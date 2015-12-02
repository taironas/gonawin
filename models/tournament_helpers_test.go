package models

import (
	"errors"
	"fmt"
	"time"
)

type testTournament struct {
	name        string
	description string
	start       time.Time
	end         time.Time
	adminID     int64
	userIDs     []int64
}

// checkTournament checks that a given tournament is equivalent to a test tournament
func checkTournament(got *Tournament, want *testTournament) error {
	var s string
	if got.Name != want.name {
		s = fmt.Sprintf("want Name == %s, got %s", want.name, got.Name)
	} else if got.Description != want.description {
		s = fmt.Sprintf("want Description == %s, got %s", want.description, got.Description)
	} else if got.Start.Sub(want.start).Hours() > 0 {
		s = fmt.Sprintf("want Start == %v, got %v", want.start, got.Start)
	} else if got.End.Sub(want.end).Hours() > 0 {
		s = fmt.Sprintf("want End == %v, got %v", want.end, got.End)
	} else if got.AdminIds[0] != want.adminID {
		s = fmt.Sprintf("want AdminId == %d, got %d", want.adminID, got.AdminIds[0])
	} else if len(got.UserIds) != len(want.userIDs) {
		s = fmt.Sprintf("want UserIds count == %d, got %d", len(want.userIDs), len(got.UserIds))
	} else {
		return nil
	}

	return errors.New(s)
}
