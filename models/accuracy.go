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

// Accuracy entity is a placeholder for progression of the accuracy of a team in a tournament.
// Teams should have a global accuracy as well as an accuracy for each tournament they participate in.
// Teams should be able to see the evolution of their accuracy for each tournament.
//
// The Team accuracy of a specific tournament is computed as follows:
//        (sum(scores of match for each team member) + previous accuracy) / (number of matches played by the team)
//
// If some participants arrive later to the tournament, previous accuracies count as 0, and this does not impact previous teams accuracy.
type Accuracy struct {
	Id           int64
	TeamId       int64
	TournamentId int64
	Accuracies   []float64
}

type AccuracyOverall struct {
	Id           int64
	TournamentId int64
	Accuracy     float64       // overall accuracy
	Progression  []Progression // progression of accuracies of team in tournament. (right now the last 5 accuracy logs)
}

type Progression struct {
	Value float64
}

// The Json version
type AccuracyJson struct {
	Id           *int64     `json:",omitempty"`
	TeamId       *int64     `json:",omitempty"`
	TournamentId *int64     `json:",omitempty"`
	Accuracies   *[]float64 `json:",omitempty"`
}

// Create an Accuracy entity.
func CreateAccuracy(c appengine.Context, teamId int64, tournamentId int64, oldmatches int) (*Accuracy, error) {
	aId, _, err := datastore.AllocateIDs(c, "Accuracy", nil, 1)
	if err != nil {
		return nil, err
	}
	key := datastore.NewKey(c, "Accuracy", "", aId, nil)
	accs := make([]float64, oldmatches)
	a := &Accuracy{aId, teamId, tournamentId, accs}
	if _, err = datastore.Put(c, key, a); err != nil {
		return nil, err
	}
	return a, nil
}

// Add accuracy to array of accuracies in Accuracy entity
func (a *Accuracy) Add(c appengine.Context, acc float64) (float64, error) {
	// add acc with previous acc / # item + 1
	log.Infof(c, "Accuracy add %v", acc)
	sum := sumFloat64(&a.Accuracies)
	log.Infof(c, "Accuracy add sum of accs: %v", a.Accuracies)
	log.Infof(c, "Accuracy add sum of accs: %v", sum)
	newAcc := float64(sum+acc) / float64(len(a.Accuracies)+1)
	log.Infof(c, "Accuracy add new acc: %v", newAcc)
	a.Accuracies = append(a.Accuracies, newAcc)
	log.Infof(c, "Accuracy add append: %v", a.Accuracies)
	return newAcc, a.Update(c)
}

// Update a team given an id and a team pointer.
func (a *Accuracy) Update(c appengine.Context) error {
	k := AccuracyKeyById(c, a.Id)
	oldAcc := new(Accuracy)
	if err := datastore.Get(c, k, oldAcc); err == nil {
		if _, err = datastore.Put(c, k, a); err != nil {
			log.Infof(c, "Accuracy.Update: error at Put, %v", err)
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

// Get an accuracy key given an id.
func AccuracyKeyById(c appengine.Context, id int64) *datastore.Key {
	key := datastore.NewKey(c, "Accuracy", "", id, nil)
	return key
}

// Get a team given an id.
func AccuracyById(c appengine.Context, id int64) (*Accuracy, error) {

	var a Accuracy
	key := datastore.NewKey(c, "Accuracy", "", id, nil)

	if err := datastore.Get(c, key, &a); err != nil {
		log.Errorf(c, " AccuracyById: accuracy not found : %v", err)
		return &a, err
	}
	return &a, nil
}
