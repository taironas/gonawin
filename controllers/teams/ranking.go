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

	"appengine"

	"github.com/santiaago/gonawin/extract"
	"github.com/santiaago/gonawin/helpers"
	"github.com/santiaago/gonawin/helpers/log"
	templateshlp "github.com/santiaago/gonawin/helpers/templates"

	mdl "github.com/santiaago/gonawin/models"
)

// Ranking handler, use it to retrieve the team ranking.
// Use this handler to get the ranking of a team.
// The ranking is an array of users (members of the team),
//	GET	/j/teams/[0-9]+/ranking/
//
func Ranking(w http.ResponseWriter, r *http.Request, u *mdl.User) error {
	if r.Method != "GET" {
		return &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeNotSupported)}
	}

	c := appengine.NewContext(r)
	desc := "Team Ranking Handler:"
	extract := extract.NewContext(c, desc, r)

	var t *mdl.Team
	var err error
	if t, err = extract.Team(); err != nil {
		return &helpers.InternalServerError{Err: err}
	}

	log.Infof(c, "%s ready to build a user array", desc)
	users := t.RankingByUser(c, 50)

	vm := buildTeamRankingViewModel(users)

	return templateshlp.RenderJson(w, c, vm)
}

type teamRankingViewModel struct {
	Users []mdl.UserJson
}

func buildTeamRankingViewModel(users []*mdl.User) teamRankingViewModel {
	fieldsToKeep := []string{"Id", "Username", "Alias", "Score"}
	u := make([]mdl.UserJson, len(users))
	helpers.TransformFromArrayOfPointers(&users, &u, fieldsToKeep)

	return teamRankingViewModel{u}

}
