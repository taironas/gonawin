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

package tournament

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"appengine"
	"appengine/datastore"

	"github.com/santiaago/purple-wing/helpers"
	tournamentinvidmdl "github.com/santiaago/purple-wing/models/tournamentInvertedIndex"
	tournamentrelmdl "github.com/santiaago/purple-wing/models/tournamentrel"
	tournamentteamrelmdl "github.com/santiaago/purple-wing/models/tournamentteamrel"
)

type Tournament struct {
	Id int64
	KeyName string
	Name string
	Description string
	Start time.Time
	End time.Time
	AdminId int64
	Created time.Time
}

type TournamentCounter struct {
	Count int64
}


func Create(r *http.Request, name string, description string, start time.Time, end time.Time, adminId int64 ) *Tournament {
	c := appengine.NewContext(r)
	// create new tournament
	tournamentID, _, _ := datastore.AllocateIDs(c, "Tournament", nil, 1)
	key := datastore.NewKey(c, "Tournament", "", tournamentID, nil)

	tournament := &Tournament{ tournamentID, helpers.TrimLower(name), name, description, start, end, adminId, time.Now() }

	_, err := datastore.Put(c, key, tournament)
	if err != nil {
		c.Errorf("Create: %v", err)
	}

	tournamentinvidmdl.Add(r, helpers.TrimLower(name),tournamentID)
	return tournament
}

func Destroy(r *http.Request, tournamentId int64) error {
	c := appengine.NewContext(r)

	if tournament, err := ById(r, tournamentId); err != nil {
		return errors.New(fmt.Sprintf("Cannot find tournament with tournamentId=%d", tournamentId))
	} else {
		key := datastore.NewKey(c, "Tournament", "", tournament.Id, nil)

		return datastore.Delete(c, key)  
	}
}

func Find(r *http.Request, filter string, value interface{}) []*Tournament {
	q := datastore.NewQuery("Tournament").Filter(filter + " =", value)

	var tournaments []*Tournament

	if _, err := q.GetAll(appengine.NewContext(r), &tournaments); err == nil {
		return tournaments
	}

	return nil
}

func ById(r *http.Request, id int64)(*Tournament, error){
	c := appengine.NewContext(r)

	var t Tournament
	key := datastore.NewKey(c, "Tournament", "", id, nil)

	if err := datastore.Get(c, key, &t); err != nil {
		c.Errorf("pw: tournament not found : %v", err)
		return &t, err
	}
	return &t, nil
}


func KeyById(r *http.Request, id int64)(*datastore.Key){
	c := appengine.NewContext(r)

	key := datastore.NewKey(c, "Tournament", "", id, nil)

	return key
}


func Update(r *http.Request, id int64, t *Tournament) error{
	c := appengine.NewContext(r)
	k := KeyById(r, id)
	oldTournament := new(Tournament)
	if err := datastore.Get(c, k, oldTournament); err == nil{
		if _, err = datastore.Put(c, k, t); err != nil {
			return err
		}
		tournamentinvidmdl.Update(r, oldTournament.Name, t.Name, id)
	}
	return nil
}

func FindAll(r *http.Request) []*Tournament {
	q := datastore.NewQuery("Tournament")

	var tournaments []*Tournament

	q.GetAll(appengine.NewContext(r), &tournaments)

	return tournaments
}

// find with respect to array of ids
func ByIds(r *http.Request, ids []int64) []*Tournament{
		
	var tournaments []*Tournament
	for _, id := range ids{
		if tournament, err := ById(r, id); err == nil{
			tournaments = append(tournaments, tournament)
		}
	}
	return tournaments
}

func Joined(r *http.Request, tournamentId int64, userId int64) bool {
	tournamentRel := tournamentrelmdl.FindByTournamentIdAndUserId(r, tournamentId, userId)
	return tournamentRel != nil
}

func Join(r *http.Request, tournamentId int64, userId int64) error {
	if tournamentRel := tournamentrelmdl.Create(r, tournamentId, userId); tournamentRel == nil {
		return errors.New("error during tournament relationship creation")
	}

	return nil
}

func Leave(r *http.Request, tournamentId int64, userId int64) error {
	return tournamentrelmdl.Destroy(r, tournamentId, userId)
}

func IsTournamentAdmin(r *http.Request, tournamentId int64, userId int64) bool {
	if tournament, err := ById(r, tournamentId); err == nil {
		return tournament.AdminId == userId
	}
	
	return false
}

// Check if a Team has joined the tournament so update relations
func TeamJoined(r *http.Request, tournamentId int64, teamId int64) bool {
	tournamentteamRel := tournamentteamrelmdl.FindByTournamentIdAndTeamId(r, tournamentId, teamId)
	return tournamentteamRel != nil
}

// Team joins the Tournament
func TeamJoin(r *http.Request, tournamentId int64, teamId int64) error {
	if tournamentteamRel := tournamentteamrelmdl.Create(r, tournamentId, teamId); tournamentteamRel == nil {
		return errors.New("error during tournament team relationship creation")
	}

	return nil
}

// Team leaves the Tournament
func TeamLeave(r *http.Request, tournamentId int64, teamId int64) error {
	return tournamentteamrelmdl.Destroy(r, tournamentId, teamId)
}

func incrementTournamentCounter(c appengine.Context, key *datastore.Key) (int64, error) {
	var x TournamentCounter
	if err := datastore.Get(c, key, &x); err != nil && err != datastore.ErrNoSuchEntity {
		return 0, err
	}
	x.Count++
	if _, err := datastore.Put(c, key, &x); err != nil {
		return 0, err
	}
	return x.Count, nil
}

func decrementTournamentCounter(c appengine.Context, key *datastore.Key) (int64, error) {
	var x TournamentCounter
	if err := datastore.Get(c, key, &x); err != nil && err != datastore.ErrNoSuchEntity {
		return 0, err
	}
	x.Count--
	if _, err := datastore.Put(c, key, &x); err != nil {
		return 0, err
	}
	return x.Count, nil
}


func GetTournamentCounter(c appengine.Context)(int64, error){
	key := datastore.NewKey(c, "TournamentCounter", "singleton", 0, nil)
	var x TournamentCounter
	if err := datastore.Get(c, key, &x); err != nil && err != datastore.ErrNoSuchEntity {
		return 0, err
	}
	return x.Count, nil
}

func GetWordFrequencyForTournament(r *http.Request, id int64, word string)int64{

	if tournaments := Find(r, "Id", id); tournaments != nil{
		return helpers.CountTerm(strings.Split(tournaments[0].KeyName, " "),word)
	}
	return 0
}
