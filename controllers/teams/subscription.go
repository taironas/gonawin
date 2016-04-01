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

package teams

import (
	"errors"
	"fmt"
	"net/http"

	"appengine"

	"github.com/taironas/gonawin/extract"
	"github.com/taironas/gonawin/helpers"
	"github.com/taironas/gonawin/helpers/log"
	templateshlp "github.com/taironas/gonawin/helpers/templates"

	mdl "github.com/taironas/gonawin/models"
)

// Join handler, use it to subscribe a user to a team.
// New user activity will be published
//
//	POST	/j/teams/join/[0-9]+/			Make a user join a team with the given id.
//
// Response: a JSON formatted team.
//
func Join(w http.ResponseWriter, r *http.Request, u *mdl.User) error {

	if r.Method != "POST" {
		return &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeNotSupported)}
	}

	c := appengine.NewContext(r)
	desc := "Team Join Handler"
	extract := extract.NewContext(c, desc, r)

	var team *mdl.Team
	var err error
	if team, err = extract.Team(); err != nil {
		return &helpers.InternalServerError{Err: err}
	}

	if team.Private {
		log.Errorf(c, "%s  Private team cannot be joined without consent.", desc)
		return &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeTeamPrivateJoinForbiden)}
	}

	if err = team.Join(c, u); err != nil {
		log.Errorf(c, "%s  error on Join team: %v", desc, err)
		return &helpers.InternalServerError{Err: errors.New(helpers.ErrorCodeInternal)}
	}

	// publish new activity
	if updatedUser, err := mdl.UserById(c, u.Id); err != nil {
		log.Errorf(c, "%s  User not found %v", desc, u.Id)
	} else {
		updatedUser.Publish(c, "team", "joined team", team.Entity(), mdl.ActivityEntity{})
	}

	vm := buildTeamJoinViewModel(team)
	return templateshlp.RenderJSON(w, c, vm)
}

// TeamJoinViewModel is the view model for Team Join handler.
//
type TeamJoinViewModel struct {
	MessageInfo string `json:",omitempty"`
	Team        mdl.TeamJSON
}

func buildTeamJoinViewModel(team *mdl.Team) TeamJoinViewModel {
	var t mdl.TeamJSON
	fieldsToKeep := []string{"Id", "Name", "AdminIds", "Private"}
	helpers.InitPointerStructure(team, &t, fieldsToKeep)

	msg := fmt.Sprintf("You joined team %s.", team.Name)
	return TeamJoinViewModel{msg, t}
}

// Leave handler, use it to make a user leave a team.
// Use this handler to leave a team.
//	POST	/j/teams/leave/[0-9]+/?			Make a user leave a team with the given id.
//
func Leave(w http.ResponseWriter, r *http.Request, u *mdl.User) error {
	if r.Method != "POST" {
		return &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeNotSupported)}
	}

	c := appengine.NewContext(r)
	desc := "Team Leave Handler:"
	extract := extract.NewContext(c, desc, r)

	var team *mdl.Team
	var err error
	team, err = extract.Team()
	if err != nil {
		return err
	}

	if mdl.IsTeamAdmin(c, team.ID, u.Id) {
		log.Errorf(c, "%s Team administrator cannot leave the team", desc)
		return &helpers.Forbidden{Err: errors.New(helpers.ErrorCodeTeamAdminCannotLeave)}
	}

	if err := team.Leave(c, u); err != nil {
		log.Errorf(c, "%s error on Leave team: %v", desc, err)
		return &helpers.InternalServerError{Err: errors.New(helpers.ErrorCodeInternal)}
	}

	var tJSON mdl.TeamJSON
	helpers.CopyToPointerStructure(team, &tJSON)
	fieldsToKeep := []string{"Id", "Name", "AdminIds", "Private"}
	helpers.KeepFields(&tJSON, fieldsToKeep)

	// publish new activity
	if updatedUser, err := mdl.UserById(c, u.Id); err != nil {
		log.Errorf(c, "User not found %v", u.Id)
	} else {
		updatedUser.Publish(c, "team", "left team", team.Entity(), mdl.ActivityEntity{})
	}

	msg := fmt.Sprintf("You left team %s.", team.Name)
	data := struct {
		MessageInfo string `json:",omitempty"`
		Team        mdl.TeamJSON
	}{
		msg,
		tJSON,
	}

	return templateshlp.RenderJSON(w, c, data)
}
