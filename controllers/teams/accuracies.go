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

// Accuracies handler, use it to get the accuracies of a team.
//
// Use this handler to get the accuracies of a team.
//	GET	/j/teams/:teamId/accuracies	retrieves all the tournament accuracies of a team with the given id.
//
// The response is an array of accurracies for the specified team group by tournament with the last 5 progressions.
//
func Accuracies(w http.ResponseWriter, r *http.Request, u *mdl.User) error {

	if r.Method != "GET" {
		return &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeNotSupported)}
	}

	c := appengine.NewContext(r)
	desc := "Team Accuracies Handler:"
	extract := extract.NewContext(c, desc, r)

	var t *mdl.Team
	var err error
	t, err = extract.Team()
	if err != nil {
		return err
	}

	log.Infof(c, "%s ready to build a acc array", desc)
	accs := t.AccuraciesGroupByTournament(c, 5)

	data := struct {
		Accuracies *[]mdl.AccuracyOverall
	}{
		accs,
	}

	return templateshlp.RenderJson(w, c, data)
}

// AccuracyByTournament handler, use it to get the team accuracies by tournament:
//
// Use this handler to get the accuracies of a team for a specific tournament.
//	GET	/j/teams/:teamId/accuracies/:tournamentId	retrieves accuracies of a team with the given id for the specified tournament.
//
// The response is an array of accurracies for the specified team team group by tournament with all it's progressions.
//
func AccuracyByTournament(w http.ResponseWriter, r *http.Request, u *mdl.User) error {
	if r.Method != "GET" {
		return &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeNotSupported)}
	}

	c := appengine.NewContext(r)
	desc := "Team Accuracies by tournament Handler:"
	extract := extract.NewContext(c, desc, r)

	var t *mdl.Team
	var err error
	t, err = extract.Team()
	if err != nil {
		return err
	}

	var tour *mdl.Tournament
	tour, err = extract.Tournament()
	if err != nil {
		return err
	}

	log.Infof(c, "%s ready to build a acc array", desc)
	acc := t.AccuracyByTournament(c, tour)

	data := struct {
		Accuracy *mdl.AccuracyOverall
	}{
		acc,
	}
	return templateshlp.RenderJson(w, c, data)
}
