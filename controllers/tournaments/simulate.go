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

package tournaments

import (
	"appengine"
	"errors"
	"math/rand"
	"net/http"

	"github.com/santiaago/purple-wing/helpers"
	"github.com/santiaago/purple-wing/helpers/handlers"
	"github.com/santiaago/purple-wing/helpers/log"
	templateshlp "github.com/santiaago/purple-wing/helpers/templates"
	tournamentmdl "github.com/santiaago/purple-wing/models/tournament"
	usermdl "github.com/santiaago/purple-wing/models/user"
)

// Simulate the scores of a phase in a tournament.
func SimulateMatchesJson(w http.ResponseWriter, r *http.Request, u *usermdl.User) error {
	c := appengine.NewContext(r)

	if r.Method == "POST" {
		tournamentId, err := handlers.PermalinkID(r, c, 3)

		if err != nil {
			log.Errorf(c, "Tournament Simulate Matches Handler: error extracting permalink err:%v", err)
			return helpers.BadRequest{errors.New(helpers.ErrorCodeTournamentNotFound)}
		}
		var tournament *tournamentmdl.Tournament
		tournament, err = tournamentmdl.ById(c, tournamentId)
		if err != nil {
			log.Errorf(c, "Tournament Simulate Matches Handler: tournament with id:%v was not found %v", tournamentId, err)
			return helpers.NotFound{errors.New(helpers.ErrorCodeTournamentNotFound)}
		}

		phase := r.FormValue("phase")
		allMatches := tournamentmdl.GetAllMatchesFromTournament(c, tournament)
		phases := tournamentmdl.MatchesGroupByPhase(allMatches)

		mapIdTeams := tournamentmdl.MapOfIdTeams(c, *tournament)
		phaseId := -1
		var results1 []int64
		var results2 []int64
		var matches []*tournamentmdl.Tmatch
		for i, ph := range phases {
			if ph.Name != phase {
				continue
			}
			phaseId = i
			for _, d := range ph.Days {
				for j, m := range d.Matches {
					// simulate match here (call set results)
					r1 := int64(rand.Intn(5))
					r2 := int64(rand.Intn(5))
					results1 = append(results1, r1)
					results2 = append(results2, r2)
					matches = append(matches, &d.Matches[j])
					log.Infof(c, "Tournament Simulate Matches: Match#%v: %v - %v | %v - %v", m.Id, mapIdTeams[m.TeamId1], mapIdTeams[m.TeamId2], r1, r2)
				}
			}
			// phase done we and not break
			break
		}
		if err = tournamentmdl.SetResults(c, matches, results1, results2, tournament); err != nil {
			log.Errorf(c, "Tournament Simulate Matches: unable to set result for matches error: %v", err)
			return helpers.NotFound{errors.New(helpers.ErrorCodeMatchesCannotUpdate)}
		}

		if phaseId >= 0 {
			// only return update phase
			matchesJson := buildMatchesFromTournament(c, *tournament)
			phasesJson := matchesGroupByPhase(matchesJson)

			data := struct {
				Phase PhaseJson
			}{
				phasesJson[phaseId],
			}

			return templateshlp.RenderJson(w, c, data)
		}
		return helpers.InternalServerError{errors.New(helpers.ErrorCodeInternal)}

	}
	return helpers.BadRequest{errors.New(helpers.ErrorCodeNotSupported)}

}
