/*
 * Copyright (c) 2014 Santiago Arias | Remy Jourde
 *
 * Permission to use, copy, modify, and distribute this software for any
 * purpose with or without fee is hereby granted, provided that the above
 * copyright notice and this permission notice appear in all copies.
 *
 * THE SOFTWARE IS PROVIDED "AS IS" AND THE AUTHOR DISCLAIMS ALL WARRANTIES
 * WITH REGARD TO THIS SOFTWARE INCLUDING ALL IMPLIED WARRANTIES OF
 * MERCHANTABILITY AND FITNESS. IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR
 * ANY SPECIAL, DIRECT, INDIRECT, OR CONSEQUENTIAL DAMAGES OR ANY DAMAGES
 * WHATSOEVER RESULTING FROM LOSS OF USE, DATA OR PROFITS, WHETHER IN AN
 * ACTION OF CONTRACT, NEGLIGENCE OR OTHER TORTIOUS ACTION, ARISING OUT OF
 * OR IN CONNECTION WITH THE USE OR PERFORMANCE OF THIS SOFTWARE.
 */

package tournaments

import (
	"errors"
	"fmt"
	"net/http"

	"appengine"

	"github.com/santiaago/gonawin/extract"
	"github.com/santiaago/gonawin/helpers"
	"github.com/santiaago/gonawin/helpers/log"
	templateshlp "github.com/santiaago/gonawin/helpers/templates"

	mdl "github.com/santiaago/gonawin/models"
)

// Join handler for tournament.
func Join(w http.ResponseWriter, r *http.Request, u *mdl.User) error {
	if r.Method != "POST" {
		return &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeNotSupported)}
	}

	c := appengine.NewContext(r)
	desc := "Tournament Join Handler:"
	extract := extract.NewContext(c, desc, r)

	var err error
	var tournament *mdl.Tournament

	if tournament, err = extract.Tournament(); err != nil {
		return err
	}

	if err := tournament.Join(c, u); err != nil {
		log.Errorf(c, "%s error on Join tournament: %v", desc, err)
		return &helpers.InternalServerError{Err: errors.New(helpers.ErrorCodeInternal)}
	}

	var tJson mdl.TournamentJson
	fieldsToKeep := []string{"Id", "Name"}
	helpers.InitPointerStructure(tournament, &tJson, fieldsToKeep)

	// publish new activity
	if updatedUser, err := mdl.UserById(c, u.Id); err != nil {
		log.Errorf(c, "User not found %v", u.Id)
	} else {
		updatedUser.Publish(c, "tournament", "joined tournament", tournament.Entity(), mdl.ActivityEntity{})
	}

	msg := fmt.Sprintf("You joined tournament %s.", tournament.Name)
	data := struct {
		MessageInfo string `json:",omitempty"`
		Tournament  mdl.TournamentJson
	}{
		msg,
		tJson,
	}

	return templateshlp.RenderJson(w, c, data)
}

// Leave handler for tournament relationships
func Leave(w http.ResponseWriter, r *http.Request, u *mdl.User) error {
	if r.Method != "POST" {
		return &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeNotSupported)}
	}

	c := appengine.NewContext(r)
	desc := "Tournament Leave Handler:"
	extract := extract.NewContext(c, desc, r)

	var err error
	var tournament *mdl.Tournament

	if tournament, err = extract.Tournament(); err != nil {
		return err
	}

	if err := tournament.Leave(c, u); err != nil {
		log.Errorf(c, "%s error on Leave team: %v", desc, err)
		return &helpers.InternalServerError{Err: errors.New(helpers.ErrorCodeInternal)}
	}

	// return the left tournament
	var tJson mdl.TournamentJson
	fieldsToKeep := []string{"Id", "Name"}
	helpers.InitPointerStructure(tournament, &tJson, fieldsToKeep)

	// publish new activity
	if updatedUser, err := mdl.UserById(c, u.Id); err != nil {
		log.Errorf(c, "User not found %v", u.Id)
	} else {
		updatedUser.Publish(c, "tournament", "left tournament", tournament.Entity(), mdl.ActivityEntity{})
	}

	msg := fmt.Sprintf("You left tournament %s.", tournament.Name)
	data := struct {
		MessageInfo string `json:",omitempty"`
		Tournament  mdl.TournamentJson
	}{
		msg,
		tJson,
	}
	return templateshlp.RenderJson(w, c, data)
}

// Join as Team handler for tournament teams realtionship.
func JoinAsTeam(w http.ResponseWriter, r *http.Request, u *mdl.User) error {
	if r.Method != "POST" {
		return &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeNotSupported)}
	}

	c := appengine.NewContext(r)
	desc := "Tournament Join as a Team Handler:"
	extract := extract.NewContext(c, desc, r)

	var err error
	var tournament *mdl.Tournament

	if tournament, err = extract.Tournament(); err != nil {
		return err
	}

	var teamId int64
	if teamId, err = extract.TeamId(); err != nil {
		return err
	}

	var team *mdl.Team
	if team, err = mdl.TeamById(c, teamId); err != nil {
		log.Errorf(c, "%s team not found: %v", desc, err)
		return &helpers.NotFound{Err: errors.New(helpers.ErrorCodeTeamNotFound)}
	}

	if err = tournament.TeamJoin(c, team); err != nil {
		log.Errorf(c, "%s error when trying to join team: %v", desc, err)
		return &helpers.InternalServerError{Err: errors.New(helpers.ErrorCodeInternal)}
	}

	var tJson mdl.TournamentJson
	fieldsToKeep := []string{"Id", "Name"}
	helpers.InitPointerStructure(tournament, &tJson, fieldsToKeep)

	// publish new activity
	var updatedTeam *mdl.Team
	if updatedTeam, err = mdl.TeamById(c, teamId); err != nil {
		log.Errorf(c, "%s team not found: %v", desc, err)
		return &helpers.NotFound{Err: errors.New(helpers.ErrorCodeTeamNotFound)}
	}
	updatedTeam.Publish(c, "tournament", "joined tournament", tournament.Entity(), mdl.ActivityEntity{})

	msg := fmt.Sprintf("Team %s joined tournament %s.", team.Name, tournament.Name)
	data := struct {
		MessageInfo string `json:",omitempty"`
		Tournament  mdl.TournamentJson
	}{
		msg,
		tJson,
	}

	return templateshlp.RenderJson(w, c, data)
}

// JSON Leave as Team handler for tournament teams realtionship.
func LeaveAsTeam(w http.ResponseWriter, r *http.Request, u *mdl.User) error {
	if r.Method != "POST" {
		return &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeNotSupported)}
	}

	c := appengine.NewContext(r)
	desc := "Tournament Leave as a Team Handler:"
	extract := extract.NewContext(c, desc, r)

	var err error
	var tournament *mdl.Tournament

	if tournament, err = extract.Tournament(); err != nil {
		return err
	}

	var teamId int64
	if teamId, err = extract.TeamId(); err != nil {
		return err
	}

	var team *mdl.Team
	if team, err = mdl.TeamById(c, teamId); err != nil {
		log.Errorf(c, "team not found: %v", desc, err)
		return &helpers.NotFound{Err: errors.New(helpers.ErrorCodeTeamNotFound)}
	}
	// leave team
	if err = tournament.TeamLeave(c, team); err != nil {
		log.Errorf(c, "%s error when trying to leave team: %v", desc, err)
		return &helpers.InternalServerError{Err: errors.New(helpers.ErrorCodeInternal)}
	}
	// return the left tournament

	var tJson mdl.TournamentJson
	fieldsToKeep := []string{"Id", "Name"}
	helpers.InitPointerStructure(tournament, &tJson, fieldsToKeep)

	// publish new activity
	var updatedTeam *mdl.Team
	if updatedTeam, err = mdl.TeamById(c, teamId); err != nil {
		log.Errorf(c, "%s team not found: %v", desc, err)
		return &helpers.NotFound{Err: errors.New(helpers.ErrorCodeTeamNotFound)}
	} else {
		updatedTeam.Publish(c, "tournament", "left tournament", tournament.Entity(), mdl.ActivityEntity{})
	}

	msg := fmt.Sprintf("Team %s left tournament %s.", team.Name, tournament.Name)
	data := struct {
		MessageInfo string `json:",omitempty"`
		Tournament  mdl.TournamentJson
	}{
		msg,
		tJson,
	}

	return templateshlp.RenderJson(w, c, data)
}
