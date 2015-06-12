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
	"errors"
	"fmt"
	"time"

	"appengine"
	"appengine/datastore"

	"github.com/santiaago/gonawin/helpers/log"
)

// A Predict entity is defined by the result of a Match: Result1 and Result2 a match id and a user id.
type Predict struct {
	Id      int64     // predict id
	UserId  int64     // user id, a prediction is binded to a single user.
	Result1 int64     // result of first team
	Result2 int64     // result of second team
	MatchId int64     // match id in tournament
	Created time.Time // date of creation
}

// Create a Predict entity given a name, a user id, a result and a match id admin id and a private mode.
func CreatePredict(c appengine.Context, userId, result1, result2, matchId int64) (*Predict, error) {

	pId, _, err := datastore.AllocateIDs(c, "Predict", nil, 1)
	if err != nil {
		return nil, err
	}
	key := datastore.NewKey(c, "Predict", "", pId, nil)
	p := &Predict{pId, userId, result1, result2, matchId, time.Now()}
	if _, err = datastore.Put(c, key, p); err != nil {
		return nil, err
	}
	return p, nil
}

// Destroy a Predict entity.
func (p *Predict) Destroy(c appengine.Context) error {

	if _, err := PredictById(c, p.Id); err != nil {
		return errors.New(fmt.Sprintf("Cannot find predict with Id=%d", p.Id))
	} else {
		key := datastore.NewKey(c, "Predict", "", p.Id, nil)

		return datastore.Delete(c, key)
	}
}

// DestroyPredicts destroys a list of predicts.
func DestroyPredicts(c appengine.Context, predictIds []int64) error {
	var keys []*datastore.Key
	for _, id := range predictIds {
		keys = append(keys, datastore.NewKey(c, "Predict", "", id, nil))
	}

	return datastore.DeleteMulti(c, keys)
}

// Search for all Predict entities with respect to a filter and a value.
func FindPredicts(c appengine.Context, filter string, value interface{}) []*Predict {

	q := datastore.NewQuery("Predict").Filter(filter+" =", value)

	var predicts []*Predict

	if _, err := q.GetAll(c, &predicts); err == nil {
		return predicts
	} else {
		log.Errorf(c, " Predict.Find, error occurred during GetAll: %v", err)
		return nil
	}
}

// Search for a Predict entity given a userId and a matchId.
// The pair (user id , match id) should be unique. So if the query returns more than one entity we return 'nil' and write in the error log.
func FindPredictByUserMatch(c appengine.Context, userId, matchId int64) *Predict {
	desc := "Predict.FindPredictByUserMatch:"
	q := datastore.NewQuery("Predict").
		Filter("UserId"+" =", userId).
		Filter("MatchId"+" =", matchId)

	var predicts []*Predict

	if _, err := q.GetAll(c, &predicts); err == nil {
		if len(predicts) == 1 {
			return predicts[0]
		} else if len(predicts) == 0 {
			log.Infof(c, "%s no predicts found.", desc)
			return nil
		} else {
			log.Errorf(c, "%s too many predicts found. pair matchId, UserId should be unique.", desc)
			return nil
		}
	} else {
		log.Errorf(c, "%s an error occurred during GetAll: %v", err)
		return nil
	}
}

// Get a Predict given an id.
func PredictById(c appengine.Context, id int64) (*Predict, error) {

	var p Predict
	key := datastore.NewKey(c, "Predict", "", id, nil)

	if err := datastore.Get(c, key, &p); err != nil {
		log.Errorf(c, "predict not found : %v", err)
		return &p, err
	}
	return &p, nil
}

// Get a Predict key given an id.
func PredictKeyById(c appengine.Context, id int64) *datastore.Key {

	key := datastore.NewKey(c, "Predict", "", id, nil)

	return key
}

// Update a Predict entity.
func (p *Predict) Update(c appengine.Context) error {
	k := PredictKeyById(c, p.Id)
	old := new(Predict)
	if err := datastore.Get(c, k, old); err == nil {
		if _, err = datastore.Put(c, k, p); err != nil {
			return err
		}
	}
	return nil
}

// Get all Predicts in datastore.
func FindAllPredicts(c appengine.Context) []*Predict {
	q := datastore.NewQuery("Predict")

	var predicts []*Predict

	if _, err := q.GetAll(c, &predicts); err != nil {
		log.Errorf(c, " Predict.FindAll, error occurred during GetAll call: %v", err)
	}
	return predicts
}

// PredictsByIds returns an array of pointers to Predict entities with respect to an array of ids.
//
func PredictsByIds(c appengine.Context, ids []int64) ([]*Predict, error) {

	predicts := make([]Predict, len(ids))
	keys := PredictKeysByIds(c, ids)

	var wrongIndexes []int

	if err := datastore.GetMulti(c, keys, predicts); err != nil {
		if me, ok := err.(appengine.MultiError); ok {
			for i, merr := range me {
				if merr == datastore.ErrNoSuchEntity {
					log.Errorf(c, "PredictsByIds, missing key: %v %v", err, keys[i].IntID())
					wrongIndexes = append(wrongIndexes, i)
				}
			}
		} else {
			return nil, err
		}
	}

	var existingPredicts []*Predict
	for i := range predicts {
		if !contains(wrongIndexes, i) {
			log.Infof(c, "PredictsByIds %v", predicts[i])
			existingPredicts = append(existingPredicts, &predicts[i])
		}
	}
	return existingPredicts, nil
}

// PredictKeysByIds returns an array of keys with respect to a given array of ids.
//
func PredictKeysByIds(c appengine.Context, ids []int64) []*datastore.Key {
	keys := make([]*datastore.Key, len(ids))
	for i, id := range ids {
		keys[i] = PredictKeyById(c, id)
	}
	return keys
}

type Predicts []*Predict

func (a Predicts) ContainsMatchId(id int64) (bool, int) {
	for i, e := range a {
		if e.MatchId == id {
			return true, i
		}
	}
	return false, -1
}
