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
	"fmt"
	"strconv"
	"strings"
	"time"

	"appengine"

	"github.com/santiaago/gonawin/helpers/log"
)

type Tday struct {
	Date    time.Time
	Matches []Tmatch
}

type Tphase struct {
	Name string
	Days []Tday
}

// ByDate implements sort.Interface for []Tday based on the date field.
type ByDate []Tday

func (a ByDate) Len() int           { return len(a) }
func (a ByDate) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByDate) Less(i, j int) bool { return a[i].Date.Before(a[j].Date) }

// Check if the match m passed as argument is the last match of a phase in a specific tournament.
// it returns a boolean and the index of the phase the match was found
func lastMatchOfPhase(c appengine.Context, m *Tmatch, phases *[]Tphase) (bool, int64) {

	for i, ph := range *phases {
		if n := len(ph.Days); n >= 1 {
			lastDay := ph.Days[n-1]
			if n = len(lastDay.Matches); n >= 1 {
				lastMatch := lastDay.Matches[n-1]
				if lastMatch.IdNumber == m.IdNumber {
					return true, int64(i)
				}
			}
		}
	}
	return false, int64(-1)
}

// UpdateNextPhase updates next phase in tournament.
//
func UpdateNextPhase(c appengine.Context, t *Tournament, currentphase *Tphase, nextphase *Tphase) error {

	// the array of phases that will be update.
	// it is an array as a phase can trigger an update in multiple phases, like semi-finals
	// trigger update of Third place and Finals
	var phases []*Tphase
	phases = append(phases, nextphase)
	// compute ranking of previous phase
	var mapOfTeams map[string]*Tteam
	mapOfTeams = make(map[string]*Tteam)

	if currentphase.Name == cFirstStage {
		// compute ranking of groups
		// get all groups.
		groups := Groups(c, t.GroupIds)
		for _, g := range groups {
			team1, argTeam1 := getFirstTeamInGroup(c, g)
			team2, _ := getSecondTeamInGroup(c, g, argTeam1)
			mapOfTeams["1"+g.Name] = team1
			mapOfTeams["2"+g.Name] = team2
		}
	} else {
		// compute ranking just by match winners
		if currentphase.Name == cFinals || currentphase.Name == cThirdPlace {
			// nothing to do.
			return nil
		}

		currentmatches := GetMatchesByPhase(c, t, currentphase.Name)
		if currentphase.Name != cSemiFinals {
			log.Infof(c, "Not SemiFinals Update Next phase: current %v", currentphase.Name)
			log.Infof(c, "Not SemiFinals Update Next phase: next %v", nextphase.Name)

			for _, m := range currentmatches {
				// ToDo: handle penalties
				if m.Result1 >= m.Result2 {
					team1, _ := TTeamById(c, m.TeamId1)
					mapOfTeams["W"+strconv.Itoa(int(m.IdNumber))] = team1
					log.Infof(c, "Not SemiFinals Update Next phase: rule: W%v teams: %v", strconv.Itoa(int(m.IdNumber)), team1.Name)

				} else if m.Result1 < m.Result2 {
					team2, _ := TTeamById(c, m.TeamId2)
					mapOfTeams["W"+strconv.Itoa(int(m.IdNumber))] = team2
					log.Infof(c, "Not SemiFinals Update Next phase: rule: W%v teams: %v", strconv.Itoa(int(m.IdNumber)), team2.Name)
				}
			}
		} else {
			// append finals phases to array of phases to update.
			var finals Tphase
			finals.Name = cFinals
			phases = append(phases, &finals)

			for _, m := range currentmatches {
				// ToDo: handle penalties
				if m.Result1 >= m.Result2 {
					team1, _ := TTeamById(c, m.TeamId1)
					team2, _ := TTeamById(c, m.TeamId2)
					mapOfTeams["W"+strconv.Itoa(int(m.IdNumber))] = team1
					mapOfTeams["L"+strconv.Itoa(int(m.IdNumber))] = team2
					log.Infof(c, "Update Next phase: rule: W%v teams: %v", strconv.Itoa(int(m.IdNumber)), team1.Name)
					log.Infof(c, "Update Next phase: rule: L%v teams: %v", strconv.Itoa(int(m.IdNumber)), team2.Name)

				} else if m.Result1 < m.Result2 {
					team2, _ := TTeamById(c, m.TeamId2)
					team1, _ := TTeamById(c, m.TeamId1)
					mapOfTeams["W"+strconv.Itoa(int(m.IdNumber))] = team2
					mapOfTeams["L"+strconv.Itoa(int(m.IdNumber))] = team1
					log.Infof(c, "Update Next phase: rule: W%v teams: %v", strconv.Itoa(int(m.IdNumber)), team2.Name)
					log.Infof(c, "Update Next phase: rule: L%v teams: %v", strconv.Itoa(int(m.IdNumber)), team1.Name)
				}
			}

		}
	}

	// update phase (matches) with new teams
	for _, ph := range phases {
		matches := GetMatchesByPhase(c, t, ph.Name)
		for i, m := range matches {
			log.Infof(c, "Update Next phase: current rule: %v", m.Rule)
			rule := strings.Split(m.Rule, " ")
			if len(rule) != 2 {
				continue
			}

			if val, ok := mapOfTeams[rule[0]]; ok {
				log.Infof(c, "Update Next phase: match found: %v", val.Name)
				matches[i].TeamId1 = val.Id
			} else {
				return errors.New(fmt.Sprintf("Cannot parse rule in tournament =%d", t.Id))
			}
			if val, ok := mapOfTeams[rule[1]]; ok {
				log.Infof(c, "Update Next phase: match found: %v", val.Name)
				matches[i].TeamId2 = val.Id
			} else {
				return errors.New(fmt.Sprintf("Cannot parse rule in tournament =%d", t.Id))
			}
			matches[i].Rule = ""
			matches[i].Ready = true
			matches[i].CanPredict = true
		}
		if err := UpdateMatches(c, matches); err != nil {
			log.Errorf(c, "Set Results: unable to set results on matches: %v", err)
			return err
		}
	}

	return nil
}
