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
	//"strings"

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
	TeamIds []byte
}

func CreateTeamInvertedIndex(r *http.Request, name string, teamIds string) *TeamInvertedIndex {
	c := appengine.NewContext(r)
	
	id, _, _ := datastore.AllocateIDs(c, "TeamInvertedIndex", nil, 1)
	key := datastore.NewKey(c, "TeamInvertedIndex", "", id, nil)
	
	byteIds := []byte(teamIds)
	t := &TeamInvertedIndex{ id, helpers.TrimLower(name), byteIds }

	_, err := datastore.Put(c, key, t)
	if err != nil {
		c.Errorf("Create: %v", err)
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

func AddToTeamInvertedIndex(r *http.Request, name string, id int64){
	// create inverted indexes for this teamID
	// to do this split words in name
	// for each word check if it exist in the table if not create a line with
	// word as key and team id as value

	//words := strings.Split(name, " ")
	//for _; w: words;{
	//	
	//}
}

func AddTournamentInvertedIndex(r *http.Request, name string, id int64){
	//words := strings.Split(name, " ")	
	// create inverted indexes for this tournamentID
	// to do this split words in name
	// for each word check if it exist in the table if not create a line with
	// word as key and tournament id as value
}

func UpdateToTeamInvertedIndex(r *http.Request, oldname string, newname string, id int64){
	
}

func UpdateTournamentInvertedIndex(r *http.Request, oldname string, newname string, id int64){
	
}







