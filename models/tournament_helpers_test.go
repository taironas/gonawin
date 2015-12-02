package models

import (
	"errors"
	"fmt"
)

func checkTournament(got *Tournament, want *Tournament) error {
	var s string
	if got.Name != want.Name {
		s = fmt.Sprintf("want Name == %s, got %s", want.Name, got.Name)
	} else if got.Description != want.Description {
		s = fmt.Sprintf("want Description == %s, got %s", want.Description, got.Description)
	} else if len(got.AdminIds) != len(want.AdminIds) {
		s = fmt.Sprintf("want AdminIds count == %d, got %d", len(want.AdminIds), len(got.AdminIds))
	} else if len(got.GroupIds) != len(want.GroupIds) {
		s = fmt.Sprintf("want GroupIds count == %d, got %d", len(want.GroupIds), len(got.GroupIds))
	} else if len(got.Matches1stStage) != len(want.Matches1stStage) {
		s = fmt.Sprintf("want Matches1stStage count == %d, got %d", len(want.Matches1stStage), len(got.Matches1stStage))
	} else if len(got.Matches2ndStage) != len(want.Matches2ndStage) {
		s = fmt.Sprintf("want Matches2ndStage count == %d, got %d", len(want.Matches2ndStage), len(got.Matches2ndStage))
	} else if len(got.UserIds) != len(want.UserIds) {
		s = fmt.Sprintf("want UserIds count == %d, got %d", len(want.UserIds), len(got.UserIds))
	} else if len(got.TeamIds) != len(want.TeamIds) {
		s = fmt.Sprintf("want TeamIds count == %d, got %d", len(want.TeamIds), len(got.TeamIds))
	} else {
		return nil
	}

	return errors.New(s)
}
