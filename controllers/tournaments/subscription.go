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
	"time"

	"appengine"

	"github.com/taironas/gonawin/extract"
	"github.com/taironas/gonawin/helpers"
	"github.com/taironas/gonawin/helpers/log"
	templateshlp "github.com/taironas/gonawin/helpers/templates"

	mdl "github.com/taironas/gonawin/models"
)

// Join handler let the user join a tournament.
//
//	POST	/j/tournaments/join/:tournamentId	let a user join a tournament.
//
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

	if time.Now().After(tournament.End) {
		return &helpers.Forbidden{Err: errors.New("Tournament has ended, you cannot join an old tournament")}
	}

	if err = tournament.Join(c, u); err != nil {
		log.Errorf(c, "%s error on Join tournament: %v", desc, err)
		return &helpers.InternalServerError{Err: errors.New(helpers.ErrorCodeInternal)}
	}

	var tJson mdl.TournamentJSON
	fieldsToKeep := []string{"Id", "Name"}
	helpers.InitPointerStructure(tournament, &tJson, fieldsToKeep)

	var updatedUser *mdl.User
	if updatedUser, err = mdl.UserById(c, u.Id); err != nil {
		log.Errorf(c, "User not found %v", u.Id)
	} else {
		updatedUser.Publish(c, "tournament", "joined tournament", tournament.Entity(), mdl.ActivityEntity{})
	}

	data := struct {
		MessageInfo string `json:",omitempty"`
		Tournament  mdl.TournamentJSON
	}{
		fmt.Sprintf("You joined tournament %s.", tournament.Name),
		tJson,
	}

	return templateshlp.RenderJson(w, c, data)
}

// JoinAsTeam makes all members of a team join the tournament.
//
//	POST	j/tournaments/joinasteam/:tournamentId/:teamId	let a team join a tournament.
//
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

	if time.Now().After(tournament.End) {
		return &helpers.Forbidden{Err: errors.New("Tournament has ended, a team cannot join an old tournament")}
	}

	var teamId int64
	if teamId, err = extract.TeamID(); err != nil {
		return err
	}

	var team *mdl.Team
	if team, err = mdl.TeamByID(c, teamId); err != nil {
		log.Errorf(c, "%s team not found: %v", desc, err)
		return &helpers.NotFound{Err: errors.New(helpers.ErrorCodeTeamNotFound)}
	}

	if err = tournament.TeamJoin(c, team); err != nil {
		log.Errorf(c, "%s error when trying to join team: %v", desc, err)
		return &helpers.InternalServerError{Err: errors.New(helpers.ErrorCodeInternal)}
	}

	var tJson mdl.TournamentJSON
	fieldsToKeep := []string{"Id", "Name"}
	helpers.InitPointerStructure(tournament, &tJson, fieldsToKeep)

	// publish new activity
	var updatedTeam *mdl.Team
	if updatedTeam, err = mdl.TeamByID(c, teamId); err != nil {
		log.Errorf(c, "%s team not found: %v", desc, err)
		return &helpers.NotFound{Err: errors.New(helpers.ErrorCodeTeamNotFound)}
	}
	updatedTeam.Publish(c, "tournament", "joined tournament", tournament.Entity(), mdl.ActivityEntity{})

	msg := fmt.Sprintf("Team %s joined tournament %s.", team.Name, tournament.Name)
	data := struct {
		MessageInfo string `json:",omitempty"`
		Tournament  mdl.TournamentJSON
	}{
		msg,
		tJson,
	}

	return templateshlp.RenderJson(w, c, data)
}

// LeaveAsTeam makes the team leave the tournament.
//
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
	if teamId, err = extract.TeamID(); err != nil {
		return err
	}

	var team *mdl.Team
	if team, err = mdl.TeamByID(c, teamId); err != nil {
		log.Errorf(c, "team not found: %v", desc, err)
		return &helpers.NotFound{Err: errors.New(helpers.ErrorCodeTeamNotFound)}
	}
	// leave team
	if err = tournament.TeamLeave(c, team); err != nil {
		log.Errorf(c, "%s error when trying to leave team: %v", desc, err)
		return &helpers.InternalServerError{Err: errors.New(helpers.ErrorCodeInternal)}
	}
	// return the left tournament

	var tJson mdl.TournamentJSON
	fieldsToKeep := []string{"Id", "Name"}
	helpers.InitPointerStructure(tournament, &tJson, fieldsToKeep)

	// publish new activity
	var updatedTeam *mdl.Team
	if updatedTeam, err = mdl.TeamByID(c, teamId); err != nil {
		log.Errorf(c, "%s team not found: %v", desc, err)
		return &helpers.NotFound{Err: errors.New(helpers.ErrorCodeTeamNotFound)}
	} else {
		updatedTeam.Publish(c, "tournament", "left tournament", tournament.Entity(), mdl.ActivityEntity{})
	}

	msg := fmt.Sprintf("Team %s left tournament %s.", team.Name, tournament.Name)
	data := struct {
		MessageInfo string `json:",omitempty"`
		Tournament  mdl.TournamentJSON
	}{
		msg,
		tJson,
	}

	return templateshlp.RenderJson(w, c, data)
}
