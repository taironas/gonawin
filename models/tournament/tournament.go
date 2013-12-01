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
	"strings"
	"time"

	"appengine"
	"appengine/datastore"

	"github.com/santiaago/purple-wing/helpers"
	"github.com/santiaago/purple-wing/helpers/log"
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


func Create(c appengine.Context, name string, description string, start time.Time, end time.Time, adminId int64 ) (*Tournament, error) {
	// create new tournament
	tournamentID, _, err := datastore.AllocateIDs(c, "Tournament", nil, 1)
	if err != nil {
		return nil, err
	}
	
	key := datastore.NewKey(c, "Tournament", "", tournamentID, nil)

	tournament := &Tournament{ tournamentID, helpers.TrimLower(name), name, description, start, end, adminId, time.Now() }

	_, err = datastore.Put(c, key, tournament)
	if err != nil {
		return nil, err
	}

	tournamentinvidmdl.Add(c, helpers.TrimLower(name),tournamentID)
	
	return tournament, nil
}

func Destroy(c appengine.Context, tournamentId int64) error {

	if tournament, err := ById(c, tournamentId); err != nil {
		return errors.New(fmt.Sprintf("Cannot find tournament with tournamentId=%d", tournamentId))
	} else {
		key := datastore.NewKey(c, "Tournament", "", tournament.Id, nil)

		return datastore.Delete(c, key)  
	}
}

func Find(c appengine.Context, filter string, value interface{}) []*Tournament {

	q := datastore.NewQuery("Tournament").Filter(filter + " =", value)

	var tournaments []*Tournament

	if _, err := q.GetAll(c, &tournaments); err == nil {
		return tournaments
	} else {
		log.Errorf(c, " Tournament.Find, error occurred during GetAll: %v", err)
		return nil
	}
}

func ById(c appengine.Context, id int64)(*Tournament, error){

	var t Tournament
	key := datastore.NewKey(c, "Tournament", "", id, nil)

	if err := datastore.Get(c, key, &t); err != nil {
		log.Errorf(c, " tournament not found : %v", err)
		return &t, err
	}
	return &t, nil
}


func KeyById(c appengine.Context, id int64)(*datastore.Key){

	key := datastore.NewKey(c, "Tournament", "", id, nil)

	return key
}

func Update(c appengine.Context, id int64, t *Tournament) error {
	
	k := KeyById(c, id)
	oldTournament := new(Tournament)
	if err := datastore.Get(c, k, oldTournament); err == nil {
		if _, err = datastore.Put(c, k, t); err != nil {
			return err
		}
		tournamentinvidmdl.Update(c, oldTournament.Name, t.Name, id)
	}
	return nil
}

func FindAll(c appengine.Context) []*Tournament {
	
	q := datastore.NewQuery("Tournament")

	var tournaments []*Tournament

	if _, err := q.GetAll(c, &tournaments); err != nil {
		log.Errorf(c, " Tournament.FindAll, error occurred during GetAll call: %v", err)
	}

	return tournaments
}

// find with respect to array of ids
func ByIds(c appengine.Context, ids []int64) []*Tournament {
	
	var tournaments []*Tournament
	for _, id := range ids {
		if tournament, err := ById(c, id); err == nil {
			tournaments = append(tournaments, tournament)
		} else {
			log.Errorf(c, " Tournament.ByIds, error occurred during ByIds call: %v", err)
		}
	}
	return tournaments
}

func Joined(c appengine.Context, tournamentId int64, userId int64) bool {
	tournamentRel := tournamentrelmdl.FindByTournamentIdAndUserId(c, tournamentId, userId)
	return tournamentRel != nil
}

func Join(c appengine.Context, tournamentId int64, userId int64) error {
	if tournamentRel, err := tournamentrelmdl.Create(c, tournamentId, userId); tournamentRel == nil {
		return errors.New(fmt.Sprintf(" Tournament.Join, error during tournament relationship creation: %v", err))
	}

	return nil
}

func Leave(c appengine.Context, tournamentId int64, userId int64) error {
	return tournamentrelmdl.Destroy(c, tournamentId, userId)
}

func IsTournamentAdmin(c appengine.Context, tournamentId int64, userId int64) bool {
	if tournament, err := ById(c, tournamentId); err == nil {
		return tournament.AdminId == userId
	}
	
	return false
}

// Check if a Team has joined the tournament so update relations
func TeamJoined(c appengine.Context, tournamentId int64, teamId int64) bool {
	tournamentteamRel := tournamentteamrelmdl.FindByTournamentIdAndTeamId(c, tournamentId, teamId)
	return tournamentteamRel != nil
}

// Team joins the Tournament
func TeamJoin(c appengine.Context, tournamentId int64, teamId int64) error {
	if tournamentteamRel, err := tournamentteamrelmdl.Create(c, tournamentId, teamId); tournamentteamRel == nil {
		return errors.New(fmt.Sprintf(" Tournament.TeamJoin, error during tournament team relationship creation: %v", err))
	}

	return nil
}

// Team leaves the Tournament
func TeamLeave(c appengine.Context, tournamentId int64, teamId int64) error {
	return tournamentteamrelmdl.Destroy(c, tournamentId, teamId)
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

func GetWordFrequencyForTournament(c appengine.Context, id int64, word string) int64 {

	if tournaments := Find(c, "Id", id); tournaments != nil{
		return helpers.CountTerm(strings.Split(tournaments[0].KeyName, " "),word)
	}
	return 0
}
