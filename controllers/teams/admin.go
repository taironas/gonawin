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

// Team add admin handler:
//
// Use this handler to add a user as admin of current team.
//	GET	/j/teams/:teamId/admin/add/:userId
//
func AddAdmin(w http.ResponseWriter, r *http.Request, u *mdl.User) error {
	c := appengine.NewContext(r)
	desc := "Team add admin Handler:"
	if r.Method == "POST" {
		// get team id and user id
		strTeamId, err1 := route.Context.Get(r, "teamId")
		if err1 != nil {
			log.Errorf(c, "%s error getting team id, err:%v", desc, err1)
			return &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeTeamNotFound)}
		}

		var teamId int64
		teamId, err1 = strconv.ParseInt(strTeamId, 0, 64)
		if err1 != nil {
			log.Errorf(c, "%s error converting team id from string to int64, err:%v", desc, err1)
			return &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeTeamNotFound)}
		}

		strUserId, err2 := route.Context.Get(r, "userId")
		if err2 != nil {
			log.Errorf(c, "%s error getting user id, err:%v", desc, err2)
			return &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeUserNotFound)}
		}

		var userId int64
		userId, err2 = strconv.ParseInt(strUserId, 0, 64)
		if err2 != nil {
			log.Errorf(c, "%s error converting user id from string to int64, err:%v", desc, err2)
			return &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeUserNotFound)}
		}

		var team *mdl.Team
		if team, err1 = mdl.TeamById(c, teamId); err1 != nil {
			log.Errorf(c, "%s team not found: %v", desc, err1)
			return &helpers.NotFound{Err: errors.New(helpers.ErrorCodeTeamNotFound)}
		}

		var newAdmin *mdl.User
		newAdmin, err := mdl.UserById(c, userId)
		log.Infof(c, "%s User: %v", desc, newAdmin)
		if err != nil {
			log.Errorf(c, "%s user not found", desc)
			return &helpers.NotFound{Err: errors.New(helpers.ErrorCodeUserNotFound)}
		}

		if err = team.AddAdmin(c, newAdmin.Id); err != nil {
			log.Errorf(c, "%s error on AddAdmin to team: %v", desc, err)
			return &helpers.InternalServerError{Err: errors.New(helpers.ErrorCodeInternal)}
		}

		var tJson mdl.TeamJson
		fieldsToKeep := []string{"Id", "Name", "AdminIds", "Private"}
		helpers.InitPointerStructure(team, &tJson, fieldsToKeep)

		msg := fmt.Sprintf("You added %s as admin of team %s.", newAdmin.Name, team.Name)
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

	return nil
}

// Team remove admin handler:
//
// Use this handler to remove a user as admin of the current team.
//	GET	/j/teams/:teamId/admin/remove/:userId
//
func RemoveAdmin(w http.ResponseWriter, r *http.Request, u *mdl.User) error {
	c := appengine.NewContext(r)
	desc := "Team remove admin Handler:"

	if r.Method == "POST" {
		// get team id and user id
		strTeamId, err1 := route.Context.Get(r, "teamId")
		if err1 != nil {
			log.Errorf(c, "%s error getting team id, err:%v", desc, err1)
			return &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeTeamNotFound)}
		}

		var teamId int64
		teamId, err1 = strconv.ParseInt(strTeamId, 0, 64)
		if err1 != nil {
			log.Errorf(c, "%s error converting team id from string to int64, err:%v", desc, err1)
			return &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeTeamNotFound)}
		}

		strUserId, err2 := route.Context.Get(r, "userId")
		if err2 != nil {
			log.Errorf(c, "%s error getting user id, err:%v", desc, err2)
			return &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeUserNotFound)}
		}

		var userId int64
		userId, err2 = strconv.ParseInt(strUserId, 0, 64)
		if err2 != nil {
			log.Errorf(c, "%s error converting user id from string to int64, err:%v", desc, err2)
			return &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeUserNotFound)}
		}

		var team *mdl.Team
		if team, err1 = mdl.TeamById(c, teamId); err1 != nil {
			log.Errorf(c, "%s team not found: %v.", desc, err1)
			return &helpers.NotFound{Err: errors.New(helpers.ErrorCodeTeamNotFound)}
		}

		var oldAdmin *mdl.User
		oldAdmin, err := mdl.UserById(c, userId)
		log.Infof(c, "%s User: %v.", desc, oldAdmin)
		if err != nil {
			log.Errorf(c, "%s user not found.", desc)
			return &helpers.NotFound{Err: errors.New(helpers.ErrorCodeUserNotFound)}
		}

		if err = team.RemoveAdmin(c, oldAdmin.Id); err != nil {
			log.Errorf(c, "%s error on RemoveAdmin to team: %v.", desc, err)
			return &helpers.InternalServerError{Err: err}
		}

		var tJson mdl.TeamJson
		fieldsToKeep := []string{"Id", "Name", "AdminIds", "Private"}
		helpers.InitPointerStructure(team, &tJson, fieldsToKeep)

		msg := fmt.Sprintf("You removed %s as admin of team %s.", oldAdmin.Name, team.Name)
		data := struct {
			MessageInfo string `json:",omitempty"`
			Team        mdl.TeamJson
		}{
			msg,
			tJson,
		}
		return templateshlp.RenderJson(w, c, data)
	}
	return nil
}
