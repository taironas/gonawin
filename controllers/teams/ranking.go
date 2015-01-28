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
	"net/http"
	"strconv"

	"appengine"

	"github.com/taironas/route"

	"github.com/santiaago/gonawin/helpers"
	"github.com/santiaago/gonawin/helpers/log"
	templateshlp "github.com/santiaago/gonawin/helpers/templates"

	mdl "github.com/santiaago/gonawin/models"
)

// Team ranking handler:
// Use this handler to get the ranking of a team.
// The ranking is an array of users (members of the team),
//	GET	/j/teams/[0-9]+/ranking/
//
func Ranking(w http.ResponseWriter, r *http.Request, u *mdl.User) error {
	c := appengine.NewContext(r)
	desc := "Team Ranking Handler:"

	if r.Method == "GET" {
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

		var t *mdl.Team
		t, err = mdl.TeamById(c, teamId)
		if err != nil {
			log.Errorf(c, "%s team with id:%v was not found %v", desc, teamId, err)
			return &helpers.NotFound{Err: errors.New(helpers.ErrorCodeTeamNotFound)}
		}

		log.Infof(c, "%s ready to build a user array", desc)
		users := t.RankingByUser(c, 50)

		fieldsToKeep := []string{"Id", "Username", "Score"}
		usersJson := make([]mdl.UserJson, len(users))
		helpers.TransformFromArrayOfPointers(&users, &usersJson, fieldsToKeep)

		data := struct {
			Users []mdl.UserJson
		}{
			usersJson,
		}

		return templateshlp.RenderJson(w, c, data)
	}
	return &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeNotSupported)}
}
