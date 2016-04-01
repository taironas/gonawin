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

	"github.com/taironas/gonawin/extract"
	"github.com/taironas/gonawin/helpers"
	templateshlp "github.com/taironas/gonawin/helpers/templates"

	mdl "github.com/taironas/gonawin/models"
)

// Accuracies handler, use it to get the accuracies of a team.
//
//	GET	/j/teams/:teamId/accuracies	retrieves all the tournament accuracies of a team with the given id.
//
// The response is an array of accurracies for the specified team group by tournament with the last 5 progressions.
//
func Accuracies(w http.ResponseWriter, r *http.Request, u *mdl.User) error {

	if r.Method != "GET" {
		return &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeNotSupported)}
	}

	c := appengine.NewContext(r)
	extract := extract.NewContext(c, "Team Accuracies Handler:", r)

	var t *mdl.Team
	var err error
	if t, err = extract.Team(); err != nil {
		return err
	}

	vm := buildTeamAccuraciesViewModel(c, t)

	return templateshlp.RenderJSON(w, c, vm)
}

type teamAccuraciesViewModel struct {
	Accuracies *[]mdl.AccuracyOverall
}

func buildTeamAccuraciesViewModel(c appengine.Context, t *mdl.Team) teamAccuraciesViewModel {
	accs := t.AccuraciesGroupByTournament(c, 5)
	return teamAccuraciesViewModel{accs}
}

// AccuracyByTournament handler, use it to get the team accuracies by tournament:
//
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
	if t, err = extract.Team(); err != nil {
		return err
	}

	var tour *mdl.Tournament
	if tour, err = extract.Tournament(); err != nil {
		return err
	}

	vm := buildAccuracyByTournamentViewModel(c, t, tour)
	return templateshlp.RenderJSON(w, c, vm)
}

type accuracyByTournamentViewModel struct {
	Accuracy *mdl.AccuracyOverall
}

func buildAccuracyByTournamentViewModel(c appengine.Context, t *mdl.Team, tour *mdl.Tournament) accuracyByTournamentViewModel {
	acc := t.AccuracyByTournament(c, tour)
	return accuracyByTournamentViewModel{acc}
}
