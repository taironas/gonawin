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
	"fmt"
	"math/rand"
	"net/http"

	"appengine"

	"github.com/taironas/gonawin/extract"
	"github.com/taironas/gonawin/helpers"
	"github.com/taironas/gonawin/helpers/log"
	templateshlp "github.com/taironas/gonawin/helpers/templates"
	mdl "github.com/taironas/gonawin/models"
)

// SimulateMatches simulates the scores of a phase in a tournament.
//    POST /j/tournaments/[0-9]+/matches/simulate?phase=:phaseName
//
func SimulateMatches(w http.ResponseWriter, r *http.Request, u *mdl.User) error {
	if r.Method != "POST" {
		return &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeNotSupported)}
	}
	c := appengine.NewContext(r)
	desc := "Tournament Simulate Matches Handler:"
	extract := extract.NewContext(c, desc, r)

	var err error
	var t *mdl.Tournament

	if t, err = extract.Tournament(); err != nil {
		return err
	}

	phase := r.FormValue("phase")
	allMatches := mdl.GetAllMatchesFromTournament(c, t)
	phases := mdl.MatchesGroupByPhase(t, allMatches)

	mapIDTeams := mdl.MapOfIDTeams(c, t)
	phaseID := -1
	var results1 []int64
	var results2 []int64
	var matches []*mdl.Tmatch
	for i, ph := range phases {
		if ph.Name != phase {
			continue
		}
		phaseID = i
		for _, d := range ph.Days {
			for j, m := range d.Matches {
				// simulate match here (call set results)
				r1 := int64(rand.Intn(5))
				r2 := int64(rand.Intn(5))
				results1 = append(results1, r1)
				results2 = append(results2, r2)
				matches = append(matches, &d.Matches[j])
				log.Infof(c, "Tournament Simulate Matches: Match#%v: %v - %v | %v - %v", m.ID, mapIDTeams[m.TeamID1], mapIDTeams[m.TeamID2], r1, r2)
			}
		}
		// phase done we and not break
		break
	}
	if err = mdl.SetResults(c, matches, results1, results2, t); err != nil {
		log.Errorf(c, "Tournament Simulate Matches: unable to set result for matches error: %v", err)
		return &helpers.NotFound{Err: errors.New(helpers.ErrorCodeMatchesCannotUpdate)}
	}

	// publish activities:
	for i, match := range matches {
		object := mdl.ActivityEntity{Id: match.TeamID1, Type: "tteam", DisplayName: mapIDTeams[match.TeamID1]}
		target := mdl.ActivityEntity{Id: match.TeamID2, Type: "tteam", DisplayName: mapIDTeams[match.TeamID2]}
		verb := ""
		if results1[i] > results2[i] {
			verb = fmt.Sprintf("won %d-%d against", results1[i], results2[i])
		} else if results1[i] < results2[i] {
			verb = fmt.Sprintf("lost %d-%d against", results1[i], results2[i])
		} else {
			verb = fmt.Sprintf("tied %d-%d against", results1[i], results2[i])
		}
		t.Publish(c, "match", verb, object, target)
	}

	if phaseID >= 0 {
		// only return update phase
		matchesJSON := buildMatchesFromTournament(c, t, u)
		phasesJSON := matchesGroupByPhase(t, matchesJSON)

		data := struct {
			Phase PhaseJson
		}{
			phasesJSON[phaseID],
		}
		return templateshlp.RenderJson(w, c, data)
	}
	return &helpers.InternalServerError{Err: errors.New(helpers.ErrorCodeInternal)}
}
