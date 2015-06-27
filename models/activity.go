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
	"time"

	"appengine"
	"appengine/datastore"

	"github.com/santiaago/gonawin/helpers/log"
)

// Activity is an update that shows the activity of the user on gonawin.
//
// An activity can be published as long as a type, a verb and an actor
// has been specified.
type Activity struct {
	Id        int64
	Type      string         // Type of the activity (welcome, team, tournament, match, accuracy, predict, score)
	Verb      string         // Describes the action
	Actor     ActivityEntity // The one who/which performs the action
	Object    ActivityEntity // The one who/which is used to performs the action (can be empty)
	Target    ActivityEntity // The one who/which is affected by the action (can be empty)
	Published time.Time
	CreatorID int64
}

// Activity Entity
type ActivityEntity struct {
	Id          int64
	Type        string
	DisplayName string // Name which will be displayed in the view
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

// Publisher interface
type Publisher interface {
	Publish(c appengine.Context, activityType string, verb string, object ActivityEntity, target ActivityEntity) error
	Entity() ActivityEntity
}

// Returns activities for a specific user.
func FindActivities(c appengine.Context, u *User, count int64, page int64) []*Activity {
	var activities []*Activity

	// loop backward on all of these ids to fetch the activities
	ids := u.ActivityIds
	log.Infof(c, "calculateStartAndEnd(%v, %v, %v)", int64(len(ids)), count, page)
	start, end := calculateStartAndEnd(int64(len(ids)), count, page)

	log.Infof(c, " Activity.FindActivities: start = %d, end = %d", start, end)

	for i := start; i >= end; i-- {
		key := datastore.NewKey(c, "Activity", "", ids[i], nil)

		var activity Activity
		if err := datastore.Get(c, key, &activity); err != nil {
			log.Errorf(c, " Activity.FindActivities: error occurred during Get call id: %v: %v", ids[i], err)
			continue // skip activity if not found..
		}
		activities = append(activities, &activity)
	}

	return activities
}

// Deletes activities for a specific user.
func DestroyActivities(c appengine.Context, activityIds []int64) error {
	var keys []*datastore.Key
	for _, id := range activityIds {
		keys = append(keys, datastore.NewKey(c, "Activity", "", id, nil))
	}

	return datastore.DeleteMulti(c, keys)
}

// Update an array of users.
func SaveActivities(c appengine.Context, activities []*Activity) error {
	var keys []*datastore.Key // := make([]*datastore.Key, len(activities))
	var acts []*Activity
	for i, _ := range activities {
		if activities[i] != nil {
			key := datastore.NewKey(c, "Activity", "", activities[i].Id, nil)
			keys = append(keys, key)
			acts = append(acts, activities[i])
		}
	}
	if _, err := datastore.PutMulti(c, keys, acts); err != nil {
		return err
	}
	return nil
}

// save an activity entity in datastore
// returns the id of the newly saved activity
func (a *Activity) save(c appengine.Context) error {
	// create new activity
	id, _, err1 := datastore.AllocateIDs(c, "Activity", nil, 1)
	if err1 != nil {
		log.Errorf(c, " Activity.save: error occurred during AllocateIDs call: %v", err1)
		return errors.New("Activity.save: unable to allocate an identifier for Activity")
	}
	key := datastore.NewKey(c, "Activity", "", id, nil)
	a.Id = id
	if _, err := datastore.Put(c, key, a); err != nil {
		log.Errorf(c, " Activity.save: error occurred during Put call: %v", err)
		return errors.New("Activity.save: unable to put Activity in Datastore")
	}
	return nil
}

// add new activity id for a specific user
func (a *Activity) AddNewActivityId(c appengine.Context, u *User) error {
	// add new activity id to user activities
	u.ActivityIds = append(u.ActivityIds, a.Id)
	return nil
	// // update user with new activity id
	// return u.Update(c)
}

// Calculates the start and the end position in the activities slice.
// Used for activities pagination.
func calculateStartAndEnd(size, count, page int64) (start, end int64) {
	if size-(count*page) >= 0 {
		start = size - (page-1)*count - 1
		end = start - count + 1
	} else {
		start = count + size - (count * page) - 1
		end = 0
	}

	return start, end
}
