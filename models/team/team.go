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

package team

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"
	
	"appengine"
	"appengine/datastore"
	
	"github.com/santiaago/purple-wing/helpers"
	teamrelmdl "github.com/santiaago/purple-wing/models/teamrel"
	teaminvidmdl "github.com/santiaago/purple-wing/models/teamInvertedIndex"
)

type Team struct {
	Id int64
	KeyName string
	Name string
	AdminId int64
	Private bool
	Created time.Time
}


type TeamCounter struct {
	Count int64
}

func Create(r *http.Request, name string, adminId int64, private bool) *Team {
	c := appengine.NewContext(r)
	// create new team
	teamId, _, _ := datastore.AllocateIDs(c, "Team", nil, 1)
	key := datastore.NewKey(c, "Team", "", teamId, nil)

	team := &Team{ teamId, helpers.TrimLower(name), name, adminId, private, time.Now() }

	_, err := datastore.Put(c, key, team)
	if err != nil {
		c.Errorf("Create: %v", err)
	}
	// udpate inverted index
	teaminvidmdl.Add(r, helpers.TrimLower(name), teamId)
	// update team counter
	errIncrement := datastore.RunInTransaction(c, func(c appengine.Context) error {
		var err1 error
		_, err1 = incrementTeamCounter(c, datastore.NewKey(c, "TeamCounter", "singleton", 0, nil))
		return err1
	}, nil)
	if errIncrement != nil {
		c.Errorf("pw: Error incrementing TeamCounter")
	}

	return team
}

func Destroy(r *http.Request, teamId int64) error {
	c := appengine.NewContext(r)
	
	if team, err := ById(r, teamId); err != nil {
		return errors.New(fmt.Sprintf("Cannot find team with teamId=%d", teamId))
	} else {
		key := datastore.NewKey(c, "Team", "", team.Id, nil)
			
		return datastore.Delete(c, key)	
	}
}

func Find(r *http.Request, filter string, value interface{}) []*Team {
	q := datastore.NewQuery("Team").Filter(filter + " =", value)
	
	var teams []*Team
	
	if _, err := q.GetAll(appengine.NewContext(r), &teams); err == nil {
		return teams
	}
	
	return nil
}

func ById(r *http.Request, id int64) (*Team, error) {
	c := appengine.NewContext(r)
	
	var t Team
	key := datastore.NewKey(c, "Team", "", id, nil)

	if err := datastore.Get(c, key, &t); err != nil {
		c.Errorf("pw: team not found : %v", err)
		return &t, err
	}
	return &t, nil
}

func KeyById(r *http.Request, id int64) (*datastore.Key) {
	c := appengine.NewContext(r)

	key := datastore.NewKey(c, "Team", "", id, nil)

	return key
}

func Update(r *http.Request, id int64, t *Team) error {
	c := appengine.NewContext(r)
	c.Infof("Team.Update Start")
	k := KeyById(r, id)
	oldTeam := new(Team)
	if err := datastore.Get(c,k, oldTeam);err == nil{
		if _, err = datastore.Put(c, k, t); err != nil {
			return err
		}
		teaminvidmdl.Update(r, oldTeam.Name, t.Name, id)
	}
	c.Infof("Team.Update End")
	return nil
}

func FindAll(r *http.Request) []*Team {
	q := datastore.NewQuery("Team")
	
	var teams []*Team
	
	q.GetAll(appengine.NewContext(r), &teams)
	
	return teams
}

func Joined(r *http.Request, teamId int64, userId int64) bool {
	teamRel := teamrelmdl.FindByTeamIdAndUserId(r, teamId, userId)
	return teamRel != nil
}

func Join(r *http.Request, teamId int64, userId int64) error {
	if teamRel := teamrelmdl.Create(r, teamId, userId); teamRel == nil {
		return errors.New("error during team relationship creation")
	}

	return nil
}

func Leave(r *http.Request, teamId int64, userId int64) error {
	return teamrelmdl.Destroy(r, teamId, userId)
}

func IsTeamAdmin(r *http.Request, teamId int64, userId int64) bool {
	if team, err := ById(r, teamId); err == nil {
		return team.AdminId == userId
	}
	
	return false
}

func incrementTeamCounter(c appengine.Context, key *datastore.Key) (int64, error) {
	var x TeamCounter
	if err := datastore.Get(c, key, &x); err != nil && err != datastore.ErrNoSuchEntity {
		return 0, err
	}
	x.Count++
	if _, err := datastore.Put(c, key, &x); err != nil {
		return 0, err
	}
	return x.Count, nil
}

func decrementTeamCounter(c appengine.Context, key *datastore.Key) (int64, error) {
	var x TeamCounter
	if err := datastore.Get(c, key, &x); err != nil && err != datastore.ErrNoSuchEntity {
		return 0, err
	}
	x.Count--
	if _, err := datastore.Put(c, key, &x); err != nil {
		return 0, err
	}
	return x.Count, nil
}

func GetTeamCounter(c appengine.Context)(int64, error){
	key := datastore.NewKey(c, "TeamCounter", "singleton", 0, nil)
	var x TeamCounter
	if err := datastore.Get(c, key, &x); err != nil && err != datastore.ErrNoSuchEntity {
		return 0, err
	}
	return x.Count, nil
}

func GetWordFrequencyForTeam(r *http.Request, id int64, word string)int64{

	if teams := Find(r, "Id", id); teams != nil{
		return helpers.CountTerm(strings.Split(teams[0].KeyName, " "),word)
	}
	return 0
}
