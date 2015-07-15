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

	"appengine"

	"github.com/taironas/gonawin/extract"
	"github.com/taironas/gonawin/helpers"
	"github.com/taironas/gonawin/helpers/log"
	templateshlp "github.com/taironas/gonawin/helpers/templates"

	mdl "github.com/taironas/gonawin/models"
)

type teamsByPhase struct {
	Name  string
	Teams []teamJson
}

// A teamJson is a variable to hold of team information.
type teamJson struct {
	Name string
	Iso  string
}

// Tournament teams handler:
// Use this handler to get the teams of a tournament.
// Returns an array of teams,
// You can specify the groupby parameter to be "phase".
//	GET	/j/tournament/[0-9]+/teams/
//
// The response is an array of teams group by phase.
//
func Teams(w http.ResponseWriter, r *http.Request, u *mdl.User) error {
	if r.Method != "GET" {
		return &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeNotSupported)}
	}

	c := appengine.NewContext(r)
	desc := "Tournament Teams Handler:"
	extract := extract.NewContext(c, desc, r)

	var err error
	var t *mdl.Tournament

	if t, err = extract.Tournament(); err != nil {
		return err
	}

	groupby := r.FormValue("groupby")
	// if wrong data, we set rankby to "phase"
	if groupby != "phase" {
		groupby = "phase"
	} else if groupby == "phase" {
		log.Infof(c, "%s ready to build a team array", desc)
		matchesJson := buildMatchesFromTournament(c, t, u)
		teamsByPhases := teamsGroupByPhase(t, matchesJson)

		data := struct {
			Phases []teamsByPhase
		}{
			teamsByPhases,
		}

		return templateshlp.RenderJson(w, c, data)
	}
	return &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeNotSupported)}
}

// From an array of Matches, create an array of Phases where the teams are grouped in.
// We use the Phases intervals and the IdNumber of each match to do this operation.
func teamsGroupByPhase(t *mdl.Tournament, matches []MatchJson) []teamsByPhase {
	var tb mdl.TournamentBuilder
	if tb = mdl.GetTournamentBuilder(t); tb == nil {
		return []teamsByPhase{}
	}

	limits := tb.MapOfPhaseIntervals()
	phaseNames := tb.ArrayOfPhases()

	phases := make([]teamsByPhase, len(limits))
	for i, _ := range phases {
		phases[i].Name = phaseNames[i]
		low := limits[phases[i].Name][0]
		high := limits[phases[i].Name][1]

		var filteredMatches []MatchJson
		for _, v := range matches {
			if v.IdNumber >= low && v.IdNumber <= high {
				filteredMatches = append(filteredMatches, v)
			}
		}
		teams := make([]teamJson, 0)
		for _, m := range filteredMatches {
			if !teamContains(teams, m.Team1) {
				t := teamJson{Name: m.Team1, Iso: m.Iso1}
				teams = append(teams, t)
			}
			if !teamContains(teams, m.Team2) {
				t := teamJson{Name: m.Team2, Iso: m.Iso2}
				teams = append(teams, t)
			}
		}
		phases[i].Teams = teams
	}
	return phases
}

func teamContains(teams []teamJson, name string) bool {
	for _, t := range teams {
		if t.Name == name {
			return true
		}
	}
	return false
}

// Update team handler. replaces a team for the second phase of the tournament.
func UpdateTeam(w http.ResponseWriter, r *http.Request, u *mdl.User) error {
	if r.Method != "POST" {
		return &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeNotSupported)}
	}

	c := appengine.NewContext(r)
	desc := "Tournament Update Team handler:"
	extract := extract.NewContext(c, desc, r)

	var err error
	var t *mdl.Tournament

	if t, err = extract.Tournament(); err != nil {
		return err
	}

	phaseName := r.FormValue("phase")
	// if wrong data exit
	if len(phaseName) == 0 {
		log.Errorf(c, "%s phase name is missing.", desc)
		return &helpers.NotFound{Err: errors.New(helpers.ErrorCodeNotSupported)}
	}

	oldName := r.FormValue("old")
	// if wrong data exit
	if len(oldName) == 0 {
		log.Errorf(c, "%s old name is missing.", desc)
		return &helpers.NotFound{Err: errors.New(helpers.ErrorCodeNotSupported)}
	}

	newName := r.FormValue("new")
	// if wrong data exit
	if len(newName) == 0 {
		log.Errorf(c, "%s new name is missing.", desc)
		return &helpers.NotFound{Err: errors.New(helpers.ErrorCodeNotSupported)}
	}
	if err := t.UpdateTournamentTeam(c, phaseName, oldName, newName); err != nil {
		log.Errorf(c, "%s something when wrong while updating a team in the tournament %v. %v", desc, t.Id, err)
		return &helpers.NotFound{Err: errors.New(helpers.ErrorCodeNotSupported)}
	}

	matchesJson := buildMatchesFromTournament(c, t, u)
	teamsByPhases := teamsGroupByPhase(t, matchesJson)

	data := struct {
		Phases []teamsByPhase
	}{
		teamsByPhases,
	}
	return templateshlp.RenderJson(w, c, data)
}
