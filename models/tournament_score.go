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
		var p *Predict
		var err1 error
		if p, err1 = u.PredictFromMatchId(c, m.Id); err1 == nil && p == nil {
			// nothing to update
			continue
		} else if err1 != nil {
			log.Errorf(c, "%s unable to get predict for current user %v: %v", desc, u.Id, err1)
			continue
		}
		score := computeScore(m, p)
		if score > 0 {
			users[i].Score += score
			usersToUpdate = append(usersToUpdate, users[i])
		}
	}
	// update users
	if err := UpdateUsers(c, usersToUpdate); err != nil {
		log.Errorf(c, "%s unable udpate users scores: %v", desc, err)
		return errors.New(helpers.ErrorCodeUsersCannotUpdate)
	}
	return nil
}

// Update the score of the teams members of the tournament.
func (t *Tournament) UpdateTeamsScore(c appengine.Context, m *Tmatch) error {
	//desc := "Update Teams score:"
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
