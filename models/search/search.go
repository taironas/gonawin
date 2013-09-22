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

package search

import (
	"net/http"
	"strings"
	"strconv"

	"appengine"
	"appengine/datastore"

	"github.com/santiaago/purple-wing/helpers"
)

type TeamInvertedIndex struct {
	Id int64
	KeyName string
	TeamIds []byte
}

type TournamentInvertedIndex struct {
	Id int64
	KeyName string
	TournamentIds []byte
}

func CreateTeamInvertedIndex(r *http.Request, word string, teamIds string) *TeamInvertedIndex {
	c := appengine.NewContext(r)
	c.Infof("pw: CreateTeamInvertedIndex")
	id, _, _ := datastore.AllocateIDs(c, "TeamInvertedIndex", nil, 1)
	key := datastore.NewKey(c, "TeamInvertedIndex", "", id, nil)
	
	byteIds := []byte(teamIds)
	t := &TeamInvertedIndex{ id, word, byteIds }

	_, err := datastore.Put(c, key, t)
	if err != nil {
		c.Errorf("pw: CreateTeamInvertedIndex: %v", err)
	}

	return t
}


func CreateTournamentInvertedIndex(r *http.Request, name string, tournamentIds string) *TournamentInvertedIndex {
	c := appengine.NewContext(r)
	
	id, _, _ := datastore.AllocateIDs(c, "TournamentInvertedIndex", nil, 1)
	key := datastore.NewKey(c, "TournamentInvertedIndex", "", id, nil)
	
	byteIds := []byte(tournamentIds)
	t := &TournamentInvertedIndex{ id, helpers.TrimLower(name), byteIds }

	_, err := datastore.Put(c, key, t)
	if err != nil {
		c.Errorf("Create: %v", err)
	}

	return t
}

// AddToTeamInvertedIndex
// Split name by words.
// For each word check if it exist in the Team Inverted Index. 
// if not create a line with the word as key and team id as value.
func AddToTeamInvertedIndex(r *http.Request, name string, id int64){

	c := appengine.NewContext(r)
	c.Infof("pw: AddToTeamInvertedIndex start")
	words := strings.Split(name, " ")
	for _, w:= range words{
		c.Infof("AddToTeamInvertedIndex: Word: %v", w)

		if inv_id := FindTeam(r, "KeyName", w); inv_id == nil{
			c.Infof("create inv id as word does not exist in table")
			CreateTeamInvertedIndex(r, w, strconv.FormatInt(id, 10))
		} else{
			// update row with new info
			c.Infof("update row with new info")
			c.Infof("current info: keyname: %v", inv_id.KeyName)
			c.Infof("current info: teamIDs: %v", string(inv_id.TeamIds))
			k := KeyByIdTeam(r, inv_id.Id)

			if newIds := mergeIds(inv_id.TeamIds, id);len(newIds) > 0{
				c.Infof("current info: new team ids: %v", newIds)
				inv_id.TeamIds  = []byte(newIds)
				if _, err := datastore.Put(c, k, inv_id);err != nil{
					c.Errorf("pw: AddToTeamInvertedIndex error on update")
				}
			}
		}
	}
	c.Infof("pw: AddToTeamInvertedIndex end")
	
}

// AddToTournamentInvertedIndex
// Split name by words.
// For each word check if it exist in the Tournament Inverted Index. 
// if not create a line with the word as key and team id as value.
func AddTournamentInvertedIndex(r *http.Request, name string, id int64){

	c := appengine.NewContext(r)
	c.Infof("pw: AddToTournamentInvertedIndex start")
	words := strings.Split(name, " ")
	for _, w:= range words{
		c.Infof("AddToTournamentInvertedIndex: Word: %v", w)

		if inv_id := FindTournament(r, "KeyName", w); inv_id == nil{
			c.Infof("create inv id as word does not exist in table")
			CreateTournamentInvertedIndex(r, w, strconv.FormatInt(id, 10))
		} else{
			// update row with new info
			c.Infof("update row with new info")
			c.Infof("current info: keyname: %v", inv_id.KeyName)
			c.Infof("current info: teamIDs: %v", string(inv_id.TournamentIds))
			k := KeyByIdTournament(r, inv_id.Id)

			if newIds := mergeIds(inv_id.TournamentIds, id);len(newIds) > 0{
				c.Infof("current info: new team ids: %v", newIds)
				inv_id.TournamentIds  = []byte(newIds)
				if _, err := datastore.Put(c, k, inv_id);err != nil{
					c.Errorf("pw: AddToTournamentInvertedIndex error on update")
				}
			}
		}
	}
	c.Infof("pw: AddToTournamentInvertedIndex end")	
}

func UpdateToTeamInvertedIndex(r *http.Request, oldname string, newname string, id int64){
	c := appengine.NewContext(r)
	c.Infof("pw: UpdateToTeamInvertedIndex start")

	// if word in old and new do nothing
	old_w := strings.Split(oldname," ")
	new_w := strings.Split(newname, " ")

	// remove id  from words in old name that are not present in new name
	for _, wo := range old_w{
		innew := false
		for _, wn := range new_w{
			if wo == wn{
				innew = true
			}
		}
		if !innew{
			c.Infof("remove: %v",wo)
			//remove it
		}
	}

	// add all id words in new name
	for _, wn := range new_w{
		inold := false
		for _, wo := range old_w{
			if wo == wn{
				inold = true
			}
		}
		if !inold{
			// add it
			c.Infof("add: %v", wn)
		}
	}

	c.Infof("pw: AddToTeamInvertedIndex end")	

}

func UpdateTournamentInvertedIndex(r *http.Request, oldname string, newname string, id int64){
	
}

func FindTeam(r *http.Request, filter string, value interface{}) *TeamInvertedIndex {
	q := datastore.NewQuery("TeamInvertedIndex").Filter(filter + " =", value).Limit(1)
	
	var t[]*TeamInvertedIndex
	
	if _, err := q.GetAll(appengine.NewContext(r), &t); err == nil && len(t) > 0 {
		return t[0]
	}
	
	return nil
}


func FindTournament(r *http.Request, filter string, value interface{}) *TournamentInvertedIndex {
	q := datastore.NewQuery("TournamentInvertedIndex").Filter(filter + " =", value).Limit(1)
	
	var t[]*TournamentInvertedIndex
	
	if _, err := q.GetAll(appengine.NewContext(r), &t); err == nil && len(t) > 0 {
		return t[0]
	}
	
	return nil
}

func KeyByIdTeam(r *http.Request, id int64) (*datastore.Key) {
	c := appengine.NewContext(r)

	key := datastore.NewKey(c, "TeamInvertedIndex", "", id, nil)

	return key
}

func KeyByIdTournament(r *http.Request, id int64) (*datastore.Key) {
	c := appengine.NewContext(r)

	key := datastore.NewKey(c, "TournamentInvertedIndex", "", id, nil)

	return key
}


// merge ids in slice of byte with id if it is not already there
// if id is already in the slice return empty string
func mergeIds(teamIds []byte, id int64) string{
	
	strTeamIds := string(teamIds)
	strIds := strings.Split(strTeamIds, " ")
	strId := strconv.FormatInt(id, 10)
	for _, i := range strIds{
		if i == strId{
			return ""
		}
	}
	return strTeamIds + " " + strId
}
