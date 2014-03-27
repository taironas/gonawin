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
	"time"

	"appengine"
	"appengine/datastore"

	"github.com/santiaago/purple-wing/helpers/log"
)

type Activity struct {
	Id        int64
	Type      string
	Verb      string
	Actor     ActivityEntity
	Object    ActivityEntity
	Target    ActivityEntity
	Published time.Time
	CreatorID int64
}

type ActivityEntity struct {
	Id          int64
	Type        string
	DisplayName string
}

type ActivityJson struct {
	Id        *int64          `json:",omitempty"`
	Type      *string         `json:",omitempty"`
	Verb      *string         `json:",omitempty"`
	Actor     *ActivityEntity `json:",omitempty"`
	Object    *ActivityEntity `json:",omitempty"`
	Target    *ActivityEntity `json:",omitempty"`
	Published *time.Time      `json:",omitempty"`
	CreatorID *int64          `json:",omitempty"`
}

// creates an activity entity,
func (a *Activity) save(c appengine.Context) error {
	// create new user
	id, _, err := datastore.AllocateIDs(c, "Activity", nil, 1)
	if err != nil {
		log.Errorf(c, "model/activity, create: %v", err)
		return errors.New("model/activity, unable to allocate an identifier for Activity")
	}
	key := datastore.NewKey(c, "Activity", "", id, nil)
	a.Id = id
	_, err = datastore.Put(c, key, a)
	if err != nil {
		log.Errorf(c, "model/activity, create: %v", err)
		return errors.New("model/activity, unable to put Activity in Datastore")
	}
	return nil
}

type Activities []*Activity

// Methods required by sort.Interface.
func (a Activities) Len() int {
	return len(a)
}
func (a Activities) Less(i, j int) bool {
	return a[i].Published.After(a[j].Published)
}
func (a Activities) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}
