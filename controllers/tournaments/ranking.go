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
	"net/http"
	"strconv"

	"appengine"

	"github.com/santiaago/gonawin/helpers"
	"github.com/santiaago/gonawin/helpers/handlers"
	"github.com/santiaago/gonawin/helpers/log"
	templateshlp "github.com/santiaago/gonawin/helpers/templates"

	mdl "github.com/santiaago/gonawin/models"
)

// Tournament ranking handler:
// Use this handler to get the ranking of a tournament.
// The ranking is an array of users (members) or teams,
// You can specify the rankby parameter to be "users" or "teams".
//	GET	/j/tournament/[0-9]+/ranking/
//
// The response is an array of users.
//
func Ranking(w http.ResponseWriter, r *http.Request, u *mdl.User) error {
	c := appengine.NewContext(r)
	desc := "Tournament Ranking Handler:"

	if r.Method == "GET" {
		tournamentId, err := handlers.PermalinkID(r, c, 3)
		if err != nil {
			log.Errorf(c, "%s error extracting permalink err:%v", desc, err)
			return &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeTournamentNotFound)}
		}

		var t *mdl.Tournament
		t, err = mdl.TournamentById(c, tournamentId)
		if err != nil {
			log.Errorf(c, "%s tournament with id:%v was not found %v", desc, tournamentId, err)
			return &helpers.NotFound{Err: errors.New(helpers.ErrorCodeTournamentNotFound)}
		}

		rankby := r.FormValue("rankby")
		// if wrong data, we set rankby to "users"
		if rankby != "teams" && rankby != "users" {
			rankby = "users"
		}

		strlimit := r.FormValue("limit")
		limit := 10
		if len(strlimit) > 0 {
			if n, err := strconv.ParseInt(strlimit, 10, 64); err != nil {
				log.Infof(c, "%s, unable to parse %v, error:%v", strlimit, err)
			} else {
				if n > 0 {
					limit = int(n)
				}
			}
		}

		if rankby == "users" {
			log.Infof(c, "%s ready to build a user array", desc)
			users := t.RankingByUser(c, limit)

			fieldsToKeep := []string{"Id", "Username", "Score"}
			usersJson := make([]mdl.UserJson, len(users))
			helpers.TransformFromArrayOfPointers(&users, &usersJson, fieldsToKeep)

			data := struct {
				Users []mdl.UserJson
			}{
				usersJson,
			}

			return templateshlp.RenderJson(w, c, data)

		} else if rankby == "teams" {
			log.Infof(c, "%s ready to build team array", desc)
			teams := t.RankingByTeam(c, limit)

			fieldsToKeep := []string{"Id", "Name", "Accuracy"}
			teamsJson := make([]mdl.TeamJson, len(teams))
			helpers.TransformFromArrayOfPointers(&teams, &teamsJson, fieldsToKeep)

			data := struct {
				Teams []mdl.TeamJson
			}{
				teamsJson,
			}
			return templateshlp.RenderJson(w, c, data)
		}
	}
	return &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeNotSupported)}
}
