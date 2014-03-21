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
	"appengine"
	"appengine/datastore"

	"github.com/santiaago/purple-wing/helpers/log"
)

type Tteam struct {
	Id   int64
	Name string
	Iso  string
}

// Get a Tteam entity by id.
func TTeamById(c appengine.Context, teamId int64) (*Tteam, error) {
	var t Tteam
	key := datastore.NewKey(c, "Tteam", "", teamId, nil)

	if err := datastore.Get(c, key, &t); err != nil {
		log.Errorf(c, "team not found : %v", err)
		return nil, err
	}
	return &t, nil
}

// From tournament entity build map of teams.
func MapOfIdTeams(c appengine.Context, tournament *Tournament) map[int64]string {

	var mapIdTeams map[int64]string
	mapIdTeams = make(map[int64]string)

	groups := Groups(c, tournament.GroupIds)
	for _, g := range groups {
		for _, t := range g.Teams {
			mapIdTeams[t.Id] = t.Name
		}
	}
	return mapIdTeams
}
