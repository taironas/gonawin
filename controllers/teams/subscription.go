/*
 * Copyright (c) 2013 Santiago Arias | Remy Jourde | Carlos Bernal
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

	"github.com/santiaago/purple-wing/helpers"
	"github.com/santiaago/purple-wing/helpers/handlers"
	"github.com/santiaago/purple-wing/helpers/log"
	templateshlp "github.com/santiaago/purple-wing/helpers/templates"

	mdl "github.com/santiaago/purple-wing/models"
)

// Join handler will make a user join a team.
// New user activity will be pushlished
//	POST	/j/teams/join/[0-9]+/			Make a user join a team with the given id.
// Reponse: a JSON formatted team.
func Join(w http.ResponseWriter, r *http.Request, u *mdl.User) error {
	c := appengine.NewContext(r)

	if r.Method == "POST" {
		// get team id
		teamId, err := handlers.PermalinkID(r, c, 4)
		if err != nil {
			log.Errorf(c, "Team Join Handler: error when extracting permalink id: %v", err)
			return &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeTeamNotFound)}
		}
		var team *mdl.Team
		if team, err = mdl.TeamById(c, teamId); err != nil {
			log.Errorf(c, "Team Join Handler: team not found: %v", err)
			return &helpers.NotFound{Err: errors.New(helpers.ErrorCodeTeamNotFound)}
		}

		if err := team.Join(c, u); err != nil {
			log.Errorf(c, "Team Join Handler: error on Join team: %v", err)
			return &helpers.InternalServerError{Err: errors.New(helpers.ErrorCodeInternal)}
		}

		var tJson mdl.TeamJson
		fieldsToKeep := []string{"Id", "Name", "AdminIds", "Private"}
		helpers.InitPointerStructure(team, &tJson, fieldsToKeep)

		// publish new activity
		object := mdl.ActivityEntity{Id: team.Id, Type: "team", DisplayName: team.Name}
		target := mdl.ActivityEntity{}
		u.Publish(c, "team", "joined team", object, target)

		msg := fmt.Sprintf("You joined team %s.", team.Name)
		data := struct {
			MessageInfo string `json:",omitempty"`
			Team        mdl.TeamJson
		}{
			msg,
			tJson,
		}

		return templateshlp.RenderJson(w, c, data)
	}
	return &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeNotSupported)}
}

// json destroy handler for team relations
func Leave(w http.ResponseWriter, r *http.Request, u *mdl.User) error {
	c := appengine.NewContext(r)

	if r.Method == "POST" {
		// get team id
		teamId, err := handlers.PermalinkID(r, c, 4)
		if err != nil {
			log.Errorf(c, "Team Leave Handler: error when extracting permalink id: %v", err)
			return &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeTeamNotFound)}
		}

		if mdl.IsTeamAdmin(c, teamId, u.Id) {
			log.Errorf(c, "Team Leave Handler: Team administrator cannot leave the team")
			return &helpers.Forbidden{Err: errors.New(helpers.ErrorCodeTeamAdminCannotLeave)}
		}

		var team *mdl.Team
		if team, err = mdl.TeamById(c, teamId); err != nil {
			log.Errorf(c, "Team Leave Handler: team not found: %v", err)
			return &helpers.NotFound{Err: errors.New(helpers.ErrorCodeTeamNotFound)}
		}
		if err := team.Leave(c, u); err != nil {
			log.Errorf(c, "Team Leave Handler: error on Leave team: %v", err)
			return &helpers.InternalServerError{Err: errors.New(helpers.ErrorCodeInternal)}
		}

		var tJson mdl.TeamJson
		helpers.CopyToPointerStructure(team, &tJson)
		fieldsToKeep := []string{"Id", "Name", "AdminIds", "Private"}
		helpers.KeepFields(&tJson, fieldsToKeep)

		// publish new activity
		object := mdl.ActivityEntity{Id: team.Id, Type: "team", DisplayName: team.Name}
		target := mdl.ActivityEntity{}
		u.Publish(c, "team", "left team", object, target)

		msg := fmt.Sprintf("You left team %s.", team.Name)
		data := struct {
			MessageInfo string `json:",omitempty"`
			Team        mdl.TeamJson
		}{
			msg,
			tJson,
		}

		return templateshlp.RenderJson(w, c, data)
	}
	return &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeNotSupported)}
}
