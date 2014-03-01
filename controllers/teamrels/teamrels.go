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

package teamrels

import (
	"errors"
	"net/http"

	"appengine"

	"github.com/santiaago/purple-wing/helpers"
	"github.com/santiaago/purple-wing/helpers/handlers"
	"github.com/santiaago/purple-wing/helpers/log"
	templateshlp "github.com/santiaago/purple-wing/helpers/templates"

	//teammdl "github.com/santiaago/purple-wing/models/team"
	//mdl "github.com/santiaago/purple-wing/models/user"
	mdl "github.com/santiaago/purple-wing/models"

	activitymdl "github.com/santiaago/purple-wing/models/activity"
)

// json create handler for team relations
func CreateJson(w http.ResponseWriter, r *http.Request, u *mdl.User) error {
	c := appengine.NewContext(r)

	if r.Method == "POST" {
		// get team id
		teamId, err := handlers.PermalinkID(r, c, 4)
		if err != nil {
			log.Errorf(c, "Teamrels Create Handler: error when extracting permalink id: %v", err)
			return &helpers.BadRequest{errors.New(helpers.ErrorCodeTeamNotFound)}
		}
		var team *mdl.Team
		if team, err = mdl.TeamById(c, teamId); err != nil {
			log.Errorf(c, "Teamrels Create Handler: team not found: %v", err)
			return &helpers.NotFound{errors.New(helpers.ErrorCodeTeamNotFound)}
		}

		if err := team.Join(c, u); err != nil {
			log.Errorf(c, "Teamrels Create Handler: error on Join team: %v", err)
			return &helpers.InternalServerError{errors.New(helpers.ErrorCodeInternal)}
		}

		var tJson mdl.TeamJson
		fieldsToKeep := []string{"Id", "Name", "AdminId", "Private"}
		helpers.InitPointerStructure(team, &tJson, fieldsToKeep)

		// publish new activity
		actor := activitymdl.ActivityEntity{ID: u.Id, Type: "user", DisplayName: u.Username}
		object := activitymdl.ActivityEntity{ID: team.Id, Type: "team", DisplayName: team.Name}
		target := activitymdl.ActivityEntity{}
		activitymdl.Publish(c, "team", "joined team", actor, object, target, u.Id)

		return templateshlp.RenderJson(w, c, tJson)
	}
	return &helpers.BadRequest{errors.New(helpers.ErrorCodeNotSupported)}
}

// json destroy handler for team relations
func DestroyJson(w http.ResponseWriter, r *http.Request, u *mdl.User) error {
	c := appengine.NewContext(r)

	if r.Method == "POST" {
		// get team id
		teamId, err := handlers.PermalinkID(r, c, 4)
		if err != nil {
			log.Errorf(c, "Teamrels Destroy Handler: error when extracting permalink id: %v", err)
			return &helpers.BadRequest{errors.New(helpers.ErrorCodeTeamNotFound)}
		}

		if mdl.IsTeamAdmin(c, teamId, u.Id) {
			log.Errorf(c, "Teamrels Destroy Handler: Team administrator cannot leave the team")
			return &helpers.Forbidden{errors.New(helpers.ErrorCodeTeamAdminCannotLeave)}
		}

		var team *mdl.Team
		if team, err = mdl.TeamById(c, teamId); err != nil {
			log.Errorf(c, "Teamrels Destroy Handler: team not found: %v", err)
			return &helpers.NotFound{errors.New(helpers.ErrorCodeTeamNotFound)}
		}
		if err := team.Leave(c, u); err != nil {
			log.Errorf(c, "Teamrels Destroy Handler: error on Leave team: %v", err)
			return &helpers.InternalServerError{errors.New(helpers.ErrorCodeInternal)}
		}

		var tJson mdl.TeamJson
		helpers.CopyToPointerStructure(team, &tJson)
		fieldsToKeep := []string{"Id", "Name", "AdminId", "Private"}
		helpers.KeepFields(&tJson, fieldsToKeep)

		// publish new activity
		actor := activitymdl.ActivityEntity{ID: u.Id, Type: "user", DisplayName: u.Username}
		object := activitymdl.ActivityEntity{ID: team.Id, Type: "team", DisplayName: team.Name}
		target := activitymdl.ActivityEntity{}
		activitymdl.Publish(c, "team", "left team", actor, object, target, u.Id)

		return templateshlp.RenderJson(w, c, tJson)
	}
	return &helpers.BadRequest{errors.New(helpers.ErrorCodeNotSupported)}
}
