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
	"sort"

	"appengine"
	"appengine/datastore"

	"github.com/santiaago/gonawin/helpers"
	"github.com/santiaago/gonawin/helpers/log"
)

type Tgroup struct {
	Id      int64
	Name    string
	Teams   []Tteam
	Matches []Tmatch
	Points  []int64
	GoalsF  []int64
	GoalsA  []int64
}

// Get a Tgroup entity by id.
func GroupById(c appengine.Context, groupId int64) (*Tgroup, error) {
	var g Tgroup
	key := datastore.NewKey(c, "Tgroup", "", groupId, nil)

	if err := datastore.Get(c, key, &g); err != nil {
		log.Errorf(c, "group not found : %v", err)
		return &g, err
	}
	return &g, nil
}

// Get an array of groups entities (Tgroup) from an array of group ids.
func Groups(c appengine.Context, groupIds []int64) []*Tgroup {

	var groups []*Tgroup

	for _, groupId := range groupIds {

		g, err := GroupById(c, groupId)
		if err != nil {
			log.Errorf(c, " Groups, cannot find group with ID=%", groupId)
		} else {
			groups = append(groups, g)
		}
	}
	return groups
}

// Get pointer to a group key given a group id.
func GroupKeyById(c appengine.Context, id int64) *datastore.Key {

	key := datastore.NewKey(c, "Tgroup", "", id, nil)
	return key
}

// Update an array of groups.
func UpdateGroups(c appengine.Context, groups []*Tgroup) error {
	keys := make([]*datastore.Key, len(groups))
	for i, _ := range keys {
		keys[i] = GroupKeyById(c, groups[i].Id)
	}
	if _, err := datastore.PutMulti(c, keys, groups); err != nil {
		return err
	}
	return nil
}

// Update a group.
func UpdateGroup(c appengine.Context, g *Tgroup) error {
	k := GroupKeyById(c, g.Id)
	oldGroup := new(Tgroup)
	if err := datastore.Get(c, k, oldGroup); err == nil {
		if _, err = datastore.Put(c, k, g); err != nil {
			return err
		}
	}
	return nil
}

// Destroy an array of groups.
func DestroyGroups(c appengine.Context, groupIds []int64) error {
	keys := make([]*datastore.Key, len(groupIds))
	for i, _ := range keys {
		keys[i] = GroupKeyById(c, groupIds[i])
	}
	if err := datastore.DeleteMulti(c, keys); err != nil {
		return err
	}
	return nil
}

// Update points in group with result of match.
func UpdatePointsAndGoals(c appengine.Context, g *Tgroup, m *Tmatch, tournament *Tournament) error {
	for i, t := range g.Teams {
		if t.Id == m.TeamId1 {
			if m.Result1 > m.Result2 {
				g.Points[i] += 3
			} else if m.Result1 == m.Result2 {
				g.Points[i] += 1
			}
			g.GoalsF[i] += m.Result1
			g.GoalsA[i] += m.Result2
		} else if t.Id == m.TeamId2 {
			if m.Result2 > m.Result1 {
				g.Points[i] += 3
			} else if m.Result2 == m.Result1 {
				g.Points[i] += 1
			}
			g.GoalsF[i] += m.Result2
			g.GoalsA[i] += m.Result1
		}
	}
	return nil
}

// Check if the match is part of a group phase in the current tournament.
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

// Get team with highest rank in group based on group points, goals for and goals against.
func getFirstTeamInGroup(c appengine.Context, g *Tgroup) (*Tteam, int) {

	points := make([]int64, len(g.Points))
	copy(points, g.Points)

	argmax1, max1 := helpers.ArgMaxInt64(points)
	points[argmax1] = -1
	argmax2, max2 := helpers.ArgMaxInt64(points)
	if max1 == max2 { // points equal try by difference
		diff := make([]int64, len(points))
		for i, _ := range points {
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

// Get team with second highest rank in group based on points, goals for and against.
func getSecondTeamInGroup(c appengine.Context, g *Tgroup, indexOfFirst int) (*Tteam, int) {

	points := make([]int64, len(g.Points))
	copy(points, g.Points)

	points[indexOfFirst] = -1

	argmax1, max1 := helpers.ArgMaxInt64(points)
	points[argmax1] = -1
	argmax2, max2 := helpers.ArgMaxInt64(points)
	if max1 == max2 { // points equal try by difference
		diff := make([]int64, len(points))
		for i, _ := range points {
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

func getNthTeamInGroup(c appengine.Context, g *Tgroup, n int64) (*Tteam, int) {

	var stats []stat
	for i, _ := range g.Teams {
		s := stat{i, g.Points[i], g.GoalsF[i] - g.GoalsA[i], g.GoalsF[i], g.GoalsA[i]}
		stats = append(stats, s)
	}

	sort.Sort(ByStats(stats))
	id := stats[n].id
	return &g.Teams[id], id
}

func getMaxTeam(groups []*Tgroup, t1 Tteam, t2 Tteam) *Tteam {
	// get groups of each team, compare stats.
	var s1, s2 *stat
	for i, g := range groups {
		for _, t := range g.Teams {
			if t.Id == t1.Id {
				s1 = &stat{i, g.Points[i], g.GoalsF[i] - g.GoalsA[i], g.GoalsF[i], g.GoalsA[i]}
			} else if t.Id == t2.Id {
				s2 = &stat{i, g.Points[i], g.GoalsF[i] - g.GoalsA[i], g.GoalsF[i], g.GoalsA[i]}
			}
		}
	}

	if s1 != nil && s2 != nil {
		stats := ByStats{*s1, *s2}
		if stats.Less(0, 1) {
			return &t2
		}
		return &t1
	}
	return nil
}

type stat struct {
	id     int
	point  int64
	diff   int64
	goalsF int64
	goalsA int64
}

type ByStats []stat

func (a ByStats) Len() int      { return len(a) }
func (a ByStats) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ByStats) Less(i, j int) bool {

	v1 := a[i]
	v2 := a[j]

	if v1.point < v2.point {
		return true
	} else if v1.point > v2.point {
		return false
	}

	// points are equal try by difference
	if v1.diff < v2.diff {
		return true
	} else if v1.diff > v2.diff {
		return false
	}

	// difference are equal try by greates number of goals scored
	if v1.goalsF < v2.goalsF {
		return true
	} else if v1.goalsF > v2.goalsF {
		return false
	}

	return false
}
