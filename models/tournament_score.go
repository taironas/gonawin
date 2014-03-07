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

package models

import (
	"errors"

	"appengine"

	"github.com/santiaago/purple-wing/helpers"
	"github.com/santiaago/purple-wing/helpers/log"
)

// Update the score of the participants to the tournament.
func (t *Tournament) UpdateUsersScore(c appengine.Context, m *Tmatch) error {
	desc := "Update users score:"

	users := t.Participants(c)
	usersToUpdate := make([]*User, 0)
	for i, u := range users {
		if score, err := u.ScoreForMatch(c, m); err != nil {
			log.Errorf(c, "%s unable udpate user %v score: %v", desc, u.Id, err)
		} else if score > 0 {
			users[i].Score += score
			usersToUpdate = append(usersToUpdate, users[i])
		}
	}
	if err := UpdateUsers(c, usersToUpdate); err != nil {
		log.Errorf(c, "%s unable udpate users scores: %v", desc, err)
		return errors.New(helpers.ErrorCodeUsersCannotUpdate)
	}
	return nil
}

// Update the accuracy of the teams members in a specific tournament.
func (t *Tournament) UpdateTeamsAccuracy(c appengine.Context, m *Tmatch) error {
	desc := "Update Teams score:"
	teams := t.Teams(c)

	teamsToUpdate := make([]*Team, 0)
	for _, team := range teams {
		sumScore := int64(0)
		players := team.Players(c)
		if len(players) == 0 {
			// a team with 0 players? this should never happen, just skip to the next.
			continue
		}
		max := 3 * len(players) // maximum score for team in current match.
		for _, u := range players {
			if score, err := u.ScoreForMatch(c, m); err != nil {
				log.Errorf(c, "%s unable udpate user %v score: %v", desc, u.Id, err)
			} else {
				sumScore += score
			}
		}

		// compute current accuracy, get accuracy entity , add accuracy to entity.
		newAcc := float64(sumScore) / float64(max)
		if acc, _ := team.TournamentAcc(c, t); acc == nil {
			log.Infof(c, "%s create accuracy if not exist", desc)
			if acc1, err := CreateAccuracy(c, team.Id, t.Id); err != nil {
				log.Errorf(c, "%s unable to create accuracy", desc)
				return err
			} else {
				team.AddTournamentAcc(c, acc1.Id, t.Id)
				log.Infof(c, "%s accuracy exists now, lets update it", desc)
				if err := acc1.Add(c, newAcc); err != nil {
					log.Errorf(c, "%s unable to add accuracy of team %v, ", desc, team.Id, err)
				}
			}
		} else {
			log.Infof(c, "%s accuracy entity exists, lets update it", desc)
			if err := acc.Add(c, newAcc); err != nil {
				log.Errorf(c, "%s unable to add accuracy of team %v, ", desc, team.Id, err)
			}
		}

		// ToDo: update team overall accuracy.
		// go through all acc of team, sum all last acc, divide by number of acc of team.
	}
	if err := UpdateTeams(c, teamsToUpdate); err != nil {
		log.Errorf(c, "%s unable udpate teams scores: %v", desc, err)
		return errors.New(helpers.ErrorCodeTeamsCannotUpdate)
	}

	return nil
}

// Computes the score to be given with respect to a match and a predict.
func computeScore(m *Tmatch, p *Predict) int64 {

	// exact result
	if (m.Result1 == p.Result1) && (m.Result2 == p.Result2) {
		return int64(3)
	}
	// wining trend
	trendW := (m.Result1 > m.Result2)
	ptrendW := (p.Result1 > p.Result2)
	if trendW == ptrendW == true {
		return int64(1)
	}
	// losign trend
	trendL := (m.Result1 < m.Result2)
	ptrendL := (p.Result1 < p.Result2)
	if trendL == ptrendL == true {
		return int64(1)
	}
	// tied trend
	trendT := (m.Result1 == m.Result2)
	ptrendT := (p.Result1 == p.Result2)
	if trendT == ptrendT == true {
		return int64(1)
	}
	// bad predict
	return int64(0)
}
