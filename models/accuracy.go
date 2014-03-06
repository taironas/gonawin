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
	"appengine/datastore"

	"github.com/santiaago/purple-wing/helpers/log"
)

type Accuracy struct {
	ID           int64
	TeamID       int64
	TournamentID int64
	Accuracies   []float64
}

type AccuracyJson struct {
	ID           *int64     `json:",omitempty"`
	TeamID       *int64     `json:",omitempty"`
	TournamentID *int64     `json:",omitempty"`
	Accuracies   *[]float64 `json:",omitempty"`
}

func CreateAccuracy(c appengine.Context, teamID int64, tournamentID int64) (*Accuracy, error) {
	var a Accuracy
	a.TeamID = teamID
	a.TournamentID = tournamentID
	a.Accuracies = make([]float64, 0)

	return a.create(c)
}

// creates an activity entity,
func (a *Accuracy) create(c appengine.Context) (*Accuracy, error) {
	// create new accuracy
	id, _, err := datastore.AllocateIDs(c, "Accuracy", nil, 1)
	if err != nil {
		log.Errorf(c, "model/accuracy, create: %v", err)
		return nil, errors.New("model/accuracy, unable to allocate an identifier for Accuracy")
	}

	key := datastore.NewKey(c, "Accuracy", "", id, nil)

	_, err = datastore.Put(c, key, a)
	if err != nil {
		log.Errorf(c, "model/accuracy, create: %v", err)
		return nil, errors.New("model/accuracy, unable to put Accuracy in Datastore")
	}

	return a, nil
}

func (a *Accuracy) Add(c appengine.Context, acc float64) error {
	// add acc with previous acc / # item + 1
	log.Infof(c, "acc: %v", acc)
	log.Infof(c, "accs: %v", a.Accuracies)
	sum := sum(&a.Accuracies)
	log.Infof(c, "sums ready : %v", sum)
	newAcc := float64(sum+acc) / float64(len(a.Accuracies)+1)
	log.Infof(c, "accuracy ready : %v", newAcc)
	a.Accuracies = append(a.Accuracies, newAcc)
	log.Infof(c, "append ready  : %v", a.Accuracies)
	return a.Update(c)
}

// Update a team given an id and a team pointer.
func (a *Accuracy) Update(c appengine.Context) error {

	k := AccuracyKeyById(c, a.ID)
	oldAcc := new(Accuracy)
	if err := datastore.Get(c, k, oldAcc); err == nil {
		if _, err = datastore.Put(c, k, a); err != nil {
			return err
		}
	}
	return nil
}

func sum(a *[]float64) (sum float64) {
	for _, v := range *a {
		sum += v
	}
	return
}

// get an accuracy key given an id
func AccuracyKeyById(c appengine.Context, id int64) *datastore.Key {

	key := datastore.NewKey(c, "Accuracy", "", id, nil)
	return key
}

func AccuracyByTeamTournament(c appengine.Context, teamID interface{}, tournamentID interface{}) []*Accuracy {

	q := datastore.NewQuery("Accuracy").
		Filter("TeamID"+" =", teamID).
		Filter("TournamentID"+" =", tournamentID)

	var accs []*Accuracy

	if _, err := q.GetAll(c, &accs); err == nil {
		return accs
	} else {
		log.Errorf(c, " Team.Find, error occurred during GetAll: %v", err)
		return nil
	}
}

// Get a team given an id.
func AccuracyById(c appengine.Context, id int64) (*Accuracy, error) {

	var a Accuracy
	key := datastore.NewKey(c, "Accuracy", "", id, nil)

	if err := datastore.Get(c, key, &a); err != nil {
		log.Errorf(c, " accuracy not found : %v", err)
		return &a, err
	}
	return &a, nil
}
