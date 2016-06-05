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
	"appengine"
	"appengine/datastore"

	"github.com/taironas/gonawin/helpers/log"
)

// Score entity is a placeholder for progression of the score of a user in a tournament.
//
// A User should have a score as well as a score for each tournament it participates in.
// It should be able to access the history of his score in a specific tournament.
//
// The score of a user evolves following the same rules.
//        If the prediction matches perfectly you get a +3
//        If prediction matches the trend you get a +1
//        If the prediction does not match the match result you get +0.
//
type Score struct {
	Id           int64
	UserId       int64
	TournamentId int64
	Scores       []int64
}

// ScoreOverall is a placeholder for the overall score of a user in different tournaments.
//
type ScoreOverall struct {
	Id              int64
	UserId          int64
	TournamentId    int64
	Score           int64
	LastProgression int64
}

// ScoreJSON is the Json version of the Score struct
//
type ScoreJSON struct {
	Id           *int64   `json:",omitempty"`
	UserId       *int64   `json:",omitempty"`
	TournamentId *int64   `json:",omitempty"`
	Scores       *[]int64 `json:",omitempty"`
}

// CreateScore creates a Score entity.
//
func CreateScore(c appengine.Context, userID int64, tournamentId int64) (*Score, error) {
	sID, _, err := datastore.AllocateIDs(c, "Score", nil, 1)
	if err != nil {
		return nil, err
	}
	key := datastore.NewKey(c, "Score", "", sID, nil)
	var scores []int64
	s := &Score{sID, userID, tournamentId, scores}
	if _, err = datastore.Put(c, key, s); err != nil {
		return nil, err
	}
	return s, nil
}

// CreateScores creates a Score entity.
//
func CreateScores(c appengine.Context, userIDs []int64, tournamentId int64) ([]*Score, []*datastore.Key, error) {
	var keys []*datastore.Key
	var scoreEntities []*Score

	for _, id := range userIDs {
		sID, _, err := datastore.AllocateIDs(c, "Score", nil, 1)
		if err != nil {
			return nil, nil, err
		}
		k := datastore.NewKey(c, "Score", "", sID, nil)
		keys = append(keys, k)

		var scores []int64
		s := &Score{sID, id, tournamentId, scores}
		scoreEntities = append(scoreEntities, s)
	}

	return scoreEntities, keys, nil
}

// SaveScores saves an array of scores to the datastore.
//
func SaveScores(c appengine.Context, scores []*Score, keys []*datastore.Key) error {
	if _, err := datastore.PutMulti(c, keys, scores); err != nil {
		return err
	}
	return nil
}

// Add adds score to array of score in Score entity.
//
func (s *Score) Add(c appengine.Context, score int64) error {
	s.Scores = append(s.Scores, score)
	return s.Update(c)
}

// AddScores adds new scores to each score entity and update all scores at the end.
//
func AddScores(c appengine.Context, tournamentScores []*Score, scores []int64) error {
	var scoresToUpdate []*Score
	for i := range tournamentScores {
		if tournamentScores[i] != nil {
			tournamentScores[i].Scores = append(tournamentScores[i].Scores, scores[i])
			scoresToUpdate = append(scoresToUpdate, tournamentScores[i])
		}
	}
	if err := UpdateScores(c, scoresToUpdate); err != nil {
		return err
	}
	return nil

}

// UpdateScores updates an array of scores.
//
func UpdateScores(c appengine.Context, scores []*Score) error {
	keys := make([]*datastore.Key, len(scores))
	for i := range keys {
		keys[i] = ScoreKeyByID(c, scores[i].Id)
	}
	if _, err := datastore.PutMulti(c, keys, scores); err != nil {
		return err
	}
	return nil
}

// Update updates a score given an id and a score pointer.
//
func (s *Score) Update(c appengine.Context) error {
	k := ScoreKeyByID(c, s.Id)
	oldScore := new(Score)
	if err := datastore.Get(c, k, oldScore); err == nil {
		if _, err = datastore.Put(c, k, s); err != nil {
			log.Errorf(c, "Score.Update: error at Put, %v", err)
			return err
		}
	}
	return nil
}

// ScoreKeyByID gets a score key given an id.
//
func ScoreKeyByID(c appengine.Context, id int64) *datastore.Key {
	key := datastore.NewKey(c, "Score", "", id, nil)
	return key
}

// ScoreByUserTournament gets an array of scores for a user, tournament pair.
//
func ScoreByUserTournament(c appengine.Context, userID interface{}, tournamentId interface{}) []*Score {

	q := datastore.NewQuery("Score").
		Filter("UserId"+" =", userID).
		Filter("TournamentId"+" =", tournamentId)

	var scores []*Score

	if _, err := q.GetAll(c, &scores); err != nil {
		log.Errorf(c, "ScoreByUserTournament: error occurred during GetAll: %v", err)
		return nil
	}

	return scores
}

// ScoreByID gets a team given an id.
//
func ScoreByID(c appengine.Context, id int64) (*Score, error) {
	var s Score
	key := datastore.NewKey(c, "Score", "", id, nil)

	if err := datastore.Get(c, key, &s); err != nil {
		log.Errorf(c, "ScoreById: Score not found : %v", err)
		return &s, err
	}
	return &s, nil
}
