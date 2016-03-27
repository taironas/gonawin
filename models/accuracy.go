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

// Accuracy is a placeholder for progression of the accuracy of a team in a tournament.
// Teams should have a global accuracy as well as an accuracy for each tournament they participate in.
// Teams should be able to see the evolution of their accuracy for each tournament.
//
// The Team accuracy of a specific tournament is computed as follows:
//        (sum(scores of match for each team member) + previous accuracy) / (number of matches played by the team)
//
// If some participants arrive later to the tournament, previous accuracies count as 0, and this does not impact previous teams accuracy.
type Accuracy struct {
	ID           int64
	TeamID       int64
	TournamentID int64
	Accuracies   []float64
}

// AccuracyOverall represents the accuracy for a tournament and its progression.
//
type AccuracyOverall struct {
	ID           int64
	TournamentID int64
	Accuracy     float64       // overall accuracy
	Progression  []Progression // progression of accuracies of team in tournament. (right now the last 5 accuracy logs)
}

// Progression holds the progression of an accuracy
//
type Progression struct {
	Value float64
}

// AccuracyJSON is the JSON representation of the Accuracy entity.
//
type AccuracyJSON struct {
	ID           *int64     `json:"Id,omitempty"`
	TeamID       *int64     `json:"TeamId,omitempty"`
	TournamentID *int64     `json:"TournamentId,omitempty"`
	Accuracies   *[]float64 `json:",omitempty"`
}

// CreateAccuracy creates an Accuracy entity.
//
func CreateAccuracy(c appengine.Context, teamID int64, tournamentID int64, oldmatches int) (*Accuracy, error) {
	accuracyID, _, err := datastore.AllocateIDs(c, "Accuracy", nil, 1)
	if err != nil {
		return nil, err
	}
	key := datastore.NewKey(c, "Accuracy", "", accuracyID, nil)
	accuracies := make([]float64, oldmatches)
	a := &Accuracy{accuracyID, teamID, tournamentID, accuracies}
	if _, err = datastore.Put(c, key, a); err != nil {
		return nil, err
	}
	return a, nil
}

// Add accuracy to array of accuracies in Accuracy entity.
//
func (a *Accuracy) Add(c appengine.Context, acc float64) (float64, error) {
	// add acc with previous acc / # item + 1
	sum := sumFloat64(&a.Accuracies)
	newAcc := float64(sum+acc) / float64(len(a.Accuracies)+1)
	a.Accuracies = append(a.Accuracies, newAcc)
	return newAcc, a.Update(c)
}

// Update a team given an id and a team pointer.
//
func (a *Accuracy) Update(c appengine.Context) error {
	k := AccuracyKeyByID(c, a.ID)
	oldAcc := new(Accuracy)
	if err := datastore.Get(c, k, oldAcc); err == nil {
		if _, err = datastore.Put(c, k, a); err != nil {
			return err
		}
	}
	return nil
}

func sumFloat64(a *[]float64) (sum float64) {
	for _, v := range *a {
		sum += v
	}
	return
}

func sumInt64(a *[]int64) (sum int64) {
	for _, v := range *a {
		sum += v
	}
	return
}

// AccuracyKeyByID gets an accuracy key given an id.
//
func AccuracyKeyByID(c appengine.Context, id int64) *datastore.Key {
	key := datastore.NewKey(c, "Accuracy", "", id, nil)
	return key
}

// AccuracyByID gets a team given an id.
//
func AccuracyByID(c appengine.Context, id int64) (*Accuracy, error) {

	var a Accuracy
	key := datastore.NewKey(c, "Accuracy", "", id, nil)

	if err := datastore.Get(c, key, &a); err != nil {
		log.Errorf(c, " AccuracyByID: accuracy not found : %v", err)
		return &a, err
	}
	return &a, nil
}
