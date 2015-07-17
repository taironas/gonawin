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

package models

import (
	"fmt"
	"strings"

	"appengine"
	"appengine/datastore"

	"github.com/taironas/gonawin/helpers/log"
)

type Tteam struct {
	Id   int64
	Name string
	Iso  string
}

// Get a Tteam entity by id.
func TTeamById(c appengine.Context, teamId int64) (*Tteam, error) {
	var t Tteam
	key := datastore.NewKey(c, "Tteam", "", teamId, nil)

	if err := datastore.Get(c, key, &t); err != nil {
		log.Errorf(c, "team not found : %v", err)
		return nil, err
	}
	return &t, nil
}

// Update tournament team.
// From a phase an old name and a new, update the next phases of the tournament.
func (t *Tournament) UpdateTournamentTeam(c appengine.Context, phaseName, oldName, newName string) error {

	var tb TournamentBuilder
	if tb = GetTournamentBuilder(t); tb == nil {
		return fmt.Errorf("TournamentBuilder not found")
	}

	mapIdTeams := tb.MapOfIdTeams(c, t)
	limits := tb.MapOfPhaseIntervals()

	oldTeamId := int64(0)
	newTeamId := int64(0)
	for k, v := range mapIdTeams {
		if v == oldName {
			oldTeamId = k
		}
		if v == newName {
			newTeamId = k
		}
		if newTeamId > 0 && oldTeamId > 0 {
			break
		}
	}

	// special treatment when old name is prefixed by "TBD"
	// or if the old name was not found in the list of teams.
	if strings.Contains(oldName, "TBD") || (oldTeamId == 0) {
		// get matches of phase
		matches := GetMatchesByPhase(c, t, phaseName)

		for _, m := range matches {
			updateMatch := false
			rule := strings.Split(m.Rule, " ")

			if rule[0] == oldName {
				rule[0] = newName
				m.TeamId1 = newTeamId
				updateMatch = true
			} else if rule[len(rule)-1] == oldName {
				rule[len(rule)-1] = newName
				m.TeamId2 = newTeamId
				updateMatch = true
			}

			if updateMatch {
				m.Rule = fmt.Sprintf("%s", strings.Join(append(rule[:0], rule...), " "))
				if err := UpdateMatch(c, m); err != nil {
					return err
				}
			}
		}

	} else {
		matches2ndPhase := Matches(c, t.Matches2ndStage)
		// low limit, all matches above this limit should be updated.
		low := limits[phaseName][0]
		for _, m := range matches2ndPhase {
			if m.IdNumber < low {
				continue
			}
			updateMatch := false
			// update teams
			if m.TeamId1 == oldTeamId {
				updateMatch = true
				m.TeamId1 = newTeamId
			}
			if m.TeamId2 == oldTeamId {
				updateMatch = true
				m.TeamId2 = newTeamId
			}
			// update rules if needed.
			rule := strings.Split(m.Rule, " ")
			if len(rule) == 2 {
				update := false
				if rule[0] == oldName {
					rule[0] = newName
					update = true
				}
				if rule[1] == oldName {
					rule[1] = newName
					update = true
				}
				if update {
					m.Rule = fmt.Sprintf("%s %s", rule[0], rule[1])
					updateMatch = true
				}
			}
			if updateMatch {
				if err := UpdateMatch(c, m); err != nil {
					return err
				}
			}
		}
	}

	return nil
}
