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
	"math/rand"

	"appengine"
	"appengine/datastore"

	"github.com/taironas/gonawin/helpers"
	"github.com/taironas/gonawin/helpers/log"
)

// Tgroup represents the group of teams of a tournament
//
type Tgroup struct {
	ID      int64
	Name    string
	Teams   []Tteam
	Matches []Tmatch
	Points  []int64
	GoalsF  []int64
	GoalsA  []int64
}

// GroupByID gets a Tgroup entity by id.
//
func GroupByID(c appengine.Context, groupID int64) (*Tgroup, error) {
	var g Tgroup
	key := datastore.NewKey(c, "Tgroup", "", groupID, nil)

	if err := datastore.Get(c, key, &g); err != nil {
		log.Errorf(c, "group not found : %v", err)
		return &g, err
	}
	return &g, nil
}

// Groups gets an array of groups entities (Tgroup) from an array of group ids.
//
func Groups(c appengine.Context, groupIDs []int64) []*Tgroup {

	var groups []*Tgroup

	for _, groupID := range groupIDs {

		g, err := GroupByID(c, groupID)
		if err != nil {
			log.Errorf(c, " Groups, cannot find group with ID=%", groupID)
		} else {
			groups = append(groups, g)
		}
	}
	return groups
}

// GroupKeyByID gets pointer to a group key given a group id.
//
func GroupKeyByID(c appengine.Context, ID int64) *datastore.Key {

	key := datastore.NewKey(c, "Tgroup", "", ID, nil)
	return key
}

// UpdateGroups updates an array of groups.
//
func UpdateGroups(c appengine.Context, groups []*Tgroup) error {
	keys := make([]*datastore.Key, len(groups))
	for i := range keys {
		keys[i] = GroupKeyByID(c, groups[i].ID)
	}
	if _, err := datastore.PutMulti(c, keys, groups); err != nil {
		return err
	}
	return nil
}

// UpdateGroup updates a group.
//
func UpdateGroup(c appengine.Context, g *Tgroup) error {
	k := GroupKeyByID(c, g.ID)
	oldGroup := new(Tgroup)
	if err := datastore.Get(c, k, oldGroup); err == nil {
		if _, err = datastore.Put(c, k, g); err != nil {
			return err
		}
	}
	return nil
}

// DestroyGroups destroys an array of groups.
//
func DestroyGroups(c appengine.Context, groupIDs []int64) error {
	keys := make([]*datastore.Key, len(groupIDs))
	for i := range keys {
		keys[i] = GroupKeyByID(c, groupIDs[i])
	}
	if err := datastore.DeleteMulti(c, keys); err != nil {
		return err
	}
	return nil
}

// UpdatePointsAndGoals update points in group with result of match.
//
func UpdatePointsAndGoals(c appengine.Context, g *Tgroup, m *Tmatch, tournament *Tournament) error {
	for i, t := range g.Teams {
		if t.ID == m.TeamId1 {
			if m.Result1 > m.Result2 {
				g.Points[i] += 3
			} else if m.Result1 == m.Result2 {
				g.Points[i]++
			}
			g.GoalsF[i] += m.Result1
			g.GoalsA[i] += m.Result2
		} else if t.ID == m.TeamId2 {
			if m.Result2 > m.Result1 {
				g.Points[i] += 3
			} else if m.Result2 == m.Result1 {
				g.Points[i]++
			}
			g.GoalsF[i] += m.Result2
			g.GoalsA[i] += m.Result1
		}
	}
	return nil
}

// IsMatchInGroup checks if the match is part of a group phase in the current tournament.
//
func (t *Tournament) IsMatchInGroup(c appengine.Context, m *Tmatch) (bool, *Tgroup) {
	groups := Groups(c, t.GroupIds)
	for i, g := range groups {
		for _, match := range g.Matches {
			if m.Id == match.Id {
				return true, groups[i]
			}
		}
	}
	return false, nil
}

// getFirstTeamInGroup gets team with highest rank in group based on group points, goals for and goals against.
//
func getFirstTeamInGroup(c appengine.Context, g *Tgroup) (*Tteam, int) {

	points := make([]int64, len(g.Points))
	copy(points, g.Points)

	argmax1, max1 := helpers.ArgMaxInt64(points)
	points[argmax1] = -1
	argmax2, max2 := helpers.ArgMaxInt64(points)
	if max1 == max2 { // points equal try by difference
		diff := make([]int64, len(points))
		for i := range points {
			diff[i] = g.GoalsF[i] - g.GoalsA[i]
		}
		if diff[argmax1] > diff[argmax2] {
			return &g.Teams[argmax1], argmax1
		} else if diff[argmax1] < diff[argmax2] {
			return &g.Teams[argmax2], argmax2
		} else { // diff equal try by greatest number of goals scored
			if g.GoalsF[argmax1] > g.GoalsF[argmax2] {
				return &g.Teams[argmax1], argmax1
			} else if g.GoalsF[argmax1] < g.GoalsF[argmax2] {
				return &g.Teams[argmax2], argmax2
			} else { // still equal try at random for now...
				if rand.Intn(2) == 0 {
					return &g.Teams[argmax1], argmax1
				}
				return &g.Teams[argmax2], argmax2
			}
		}
	} else {
		return &g.Teams[argmax1], argmax1
	}
}

// getSecondTeamInGroup gets team with second highest rank in group based on points, goals for and against.
func getSecondTeamInGroup(c appengine.Context, g *Tgroup, indexOfFirst int) (*Tteam, int) {

	points := make([]int64, len(g.Points))
	copy(points, g.Points)

	points[indexOfFirst] = -1

	argmax1, max1 := helpers.ArgMaxInt64(points)
	points[argmax1] = -1
	argmax2, max2 := helpers.ArgMaxInt64(points)
	if max1 == max2 { // points equal try by difference
		diff := make([]int64, len(points))
		for i := range points {
			diff[i] = g.GoalsF[i] - g.GoalsA[i]
		}
		if diff[argmax1] > diff[argmax2] {
			return &g.Teams[argmax1], argmax1
		} else if diff[argmax1] < diff[argmax2] {
			return &g.Teams[argmax2], argmax2
		} else { // diff equal try by greatest number of goals scored
			if g.GoalsF[argmax1] > g.GoalsF[argmax2] {
				return &g.Teams[argmax1], argmax1
			} else if g.GoalsF[argmax1] < g.GoalsF[argmax2] {
				return &g.Teams[argmax2], argmax2
			} else { // still equal try at random for now...
				if rand.Intn(2) == 0 {
					return &g.Teams[argmax1], argmax1
				}
				return &g.Teams[argmax2], argmax2
			}
		}
	} else {
		return &g.Teams[argmax1], argmax1
	}
}
