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
	"strconv"

	"appengine"

	"github.com/taironas/route"

	"github.com/santiaago/gonawin/helpers"
	"github.com/santiaago/gonawin/helpers/log"
	templateshlp "github.com/santiaago/gonawin/helpers/templates"

	mdl "github.com/santiaago/gonawin/models"
)

// Join handler will make a user join a team.
// New user activity will be pushlished
//	POST	/j/teams/join/[0-9]+/			Make a user join a team with the given id.
// Reponse: a JSON formatted team.
func Join(w http.ResponseWriter, r *http.Request, u *mdl.User) error {
	c := appengine.NewContext(r)
	desc := "Team Join Handler:"

	if r.Method == "POST" {
		// get team id
		strTeamId, err := route.Context.Get(r, "teamId")
		if err != nil {
			log.Errorf(c, "%s error getting team id, err:%v", desc, err)
			return &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeTeamNotFound)}
		}

		var teamId int64
		teamId, err = strconv.ParseInt(strTeamId, 0, 64)
		if err != nil {
			log.Errorf(c, "%s error converting team id from string to int64, err:%v", desc, err)
			return &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeTeamNotFound)}
		}

		var team *mdl.Team
		if team, err = mdl.TeamById(c, teamId); err != nil {
			log.Errorf(c, "%s team with id:%v was not found %v", desc, teamId, err)
			return &helpers.NotFound{Err: errors.New(helpers.ErrorCodeTeamNotFound)}
		}

		if err := team.Join(c, u); err != nil {
			log.Errorf(c, "%s  error on Join team: %v", desc, err)
			return &helpers.InternalServerError{Err: errors.New(helpers.ErrorCodeInternal)}
		}

		var tJson mdl.TeamJson
		fieldsToKeep := []string{"Id", "Name", "AdminIds", "Private"}
		helpers.InitPointerStructure(team, &tJson, fieldsToKeep)

		// publish new activity
		if updatedUser, err := mdl.UserById(c, u.Id); err != nil {
			log.Errorf(c, "%s  User not found %v", desc, u.Id)
		} else {
			updatedUser.Publish(c, "team", "joined team", team.Entity(), mdl.ActivityEntity{})
		}

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

// leave handler for team relations.
// Use this handler to leave a team.
//	POST	/j/teams/leave/[0-9]+/?			Make a user leave a team with the given id.
//
func Leave(w http.ResponseWriter, r *http.Request, u *mdl.User) error {
	desc := "Team Leave Handler:"
	c := appengine.NewContext(r)

	if r.Method == "POST" {
		// get team id
		strTeamId, err := route.Context.Get(r, "teamId")
		if err != nil {
			log.Errorf(c, "%s error getting team id, err:%v", desc, err)
			return &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeTeamNotFound)}
		}

		var teamId int64
		teamId, err = strconv.ParseInt(strTeamId, 0, 64)
		if err != nil {
			log.Errorf(c, "%s error converting team id from string to int64, err:%v", desc, err)
			return &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeTeamNotFound)}
		}

		if mdl.IsTeamAdmin(c, teamId, u.Id) {
			log.Errorf(c, "%s Team administrator cannot leave the team", desc)
			return &helpers.Forbidden{Err: errors.New(helpers.ErrorCodeTeamAdminCannotLeave)}
		}

		var team *mdl.Team
		if team, err = mdl.TeamById(c, teamId); err != nil {
			log.Errorf(c, "%s team not found: %v", desc, err)
			return &helpers.NotFound{Err: errors.New(helpers.ErrorCodeTeamNotFound)}
		}
		if err := team.Leave(c, u); err != nil {
			log.Errorf(c, "%s error on Leave team: %v", desc, err)
			return &helpers.InternalServerError{Err: errors.New(helpers.ErrorCodeInternal)}
		}

		var tJson mdl.TeamJson
		helpers.CopyToPointerStructure(team, &tJson)
		fieldsToKeep := []string{"Id", "Name", "AdminIds", "Private"}
		helpers.KeepFields(&tJson, fieldsToKeep)

		// publish new activity
		if updatedUser, err := mdl.UserById(c, u.Id); err != nil {
			log.Errorf(c, "User not found %v", u.Id)
		} else {
			updatedUser.Publish(c, "team", "left team", team.Entity(), mdl.ActivityEntity{})
		}

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
