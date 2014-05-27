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
	"fmt"
	"time"

	"appengine"
	"appengine/datastore"

	"github.com/santiaago/gonawin/helpers/log"
)

// A Price entity is defined by a description of the price that the winner gets for a specific tournament.
type Price struct {
	Id             int64     // price id
	TeamId         int64     // team id, a price is binded to a single team.
	TournamentId   int64     // tournament id, a price is binded to a single team.
	TournamentName string    // tournament name.
	Description    string    // the description of the price
	Created        time.Time // date of creation
}

// Create a Price entity given a description, a team id and a tournament id.
func CreatePrice(c appengine.Context, teamId, tournamentId int64, tournamentName string, description string) (*Price, error) {

	pId, _, err := datastore.AllocateIDs(c, "Price", nil, 1)
	if err != nil {
		return nil, err
	}
	key := datastore.NewKey(c, "Price", "", pId, nil)
	p := &Price{pId, teamId, tournamentId, tournamentName, description, time.Now()}
	if _, err = datastore.Put(c, key, p); err != nil {
		return nil, err
	}
	return p, nil
}

// Destroy a Price entity.
func (p *Price) Destroy(c appengine.Context) error {

	if _, err := PriceById(c, p.Id); err != nil {
		return errors.New(fmt.Sprintf("Cannot find price with Id=%d", p.Id))
	} else {
		key := datastore.NewKey(c, "Price", "", p.Id, nil)

		return datastore.Delete(c, key)
	}
}

// Search for a Predict entity given a userId and a matchId.
// The pair (user id , match id) should be unique. So if the query returns more than one entity we return 'nil' and write in the error log.
func FindPricesByTeam(c appengine.Context, teamId int64) []*Price {
	desc := "Price.FindPriceByTeam:"
	q := datastore.NewQuery("Price").
		Filter("TeamId"+" =", teamId)

	var prices []*Price

	if _, err := q.GetAll(c, &prices); err == nil {
		if len(prices) == 0 {
			log.Infof(c, "%s no prices found.", desc)
			return nil
		} else {
			return prices
		}
	} else {
		log.Errorf(c, "%s an error occurred during GetAll: %v", err)
		return nil
	}
}

// Get a Price given an id.
func PriceById(c appengine.Context, id int64) (*Price, error) {

	var p Price
	key := datastore.NewKey(c, "Price", "", id, nil)

	if err := datastore.Get(c, key, &p); err != nil {
		log.Errorf(c, "Price not found : %v", err)
		return &p, err
	}
	return &p, nil
}

// Get a Price key given an id.
func PriceKeyById(c appengine.Context, id int64) *datastore.Key {

	key := datastore.NewKey(c, "Price", "", id, nil)
	return key
}

// Update a Predict entity.
func (p *Price) Update(c appengine.Context) error {
	k := PriceKeyById(c, p.Id)
	old := new(Price)
	if err := datastore.Get(c, k, old); err == nil {
		if _, err = datastore.Put(c, k, p); err != nil {
			return err
		}
	}
	return nil
}

// Get an array of pointers to Price entities with respect to an array of ids.
func PricesByIds(c appengine.Context, ids []int64) []*Price {

	var prices []*Price
	for _, id := range ids {
		if p, err := PriceById(c, id); err == nil {
			prices = append(prices, p)
		} else {
			log.Errorf(c, " Prices.ByIds, error occurred during ByIds call: %v", err)
		}
	}
	return prices
}
