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
	"appengine"
	"appengine/datastore"

	"github.com/santiaago/purple-wing/helpers/log"
)

// Score entity, a placeholder for progression of the score of a user in a tournament.
type Score struct {
	Id           int64
	UserId       int64
	TournamentId int64
	Scores       []int64
}

// The Json version
type ScoreJson struct {
	Id           *int64   `json:",omitempty"`
	UserId       *int64   `json:",omitempty"`
	TournamentId *int64   `json:",omitempty"`
	Scores       *[]int64 `json:",omitempty"`
}

// create a Score entity.
func CreateScore(c appengine.Context, userId int64, tournamentId int64) (*Score, error) {
	sId, _, err := datastore.AllocateIDs(c, "Score", nil, 1)
	if err != nil {
		return nil, err
	}
	key := datastore.NewKey(c, "Score", "", sId, nil)
	scores := make([]int64, 0)
	s := &Score{sId, userId, tournamentId, scores}
	if _, err = datastore.Put(c, key, s); err != nil {
		return nil, err
	}
	return s, nil
}

// Add accuracy to array of accuracies in Accuracy entity
func (s *Score) Add(c appengine.Context, score int64) error {
	s.Scores = append(s.Scores, score)
	return s.Update(c)
}

// Update a team given an id and a team pointer.
func (s *Score) Update(c appengine.Context) error {
	k := ScoreKeyById(c, s.Id)
	oldScore := new(Score)
	if err := datastore.Get(c, k, oldScore); err == nil {
		if _, err = datastore.Put(c, k, s); err != nil {
			log.Infof(c, "Score.Update: error at Put, %v", err)
			return err
		}
	}
	return nil
}

// Get a score key given an id
func ScoreKeyById(c appengine.Context, id int64) *datastore.Key {
	key := datastore.NewKey(c, "Score", "", id, nil)
	return key
}

func ScoreByUserTournament(c appengine.Context, userId interface{}, tournamentId interface{}) []*Score {

	q := datastore.NewQuery("Score").
		Filter("UserId"+" =", userId).
		Filter("TournamentId"+" =", tournamentId)

	var scores []*Score

	if _, err := q.GetAll(c, &scores); err == nil {
		return scores
	} else {
		log.Errorf(c, "ScoreByUserTournament: error occurred during GetAll: %v", err)
		return nil
	}
}

// Get a team given an id.
func ScoreById(c appengine.Context, id int64) (*Score, error) {

	var s Score
	key := datastore.NewKey(c, "Score", "", id, nil)

	if err := datastore.Get(c, key, &s); err != nil {
		log.Errorf(c, " ScoreById: Score not found : %v", err)
		return &s, err
	}
	return &s, nil
}
