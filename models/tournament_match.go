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
	"sort"
	"time"

	"appengine"
	"appengine/datastore"

	"github.com/taironas/gonawin/helpers"
	"github.com/taironas/gonawin/helpers/log"
)

// Tmatch represents a tournament match.
//
type Tmatch struct {
	Id         int64     // datastore match id
	IdNumber   int64     // id of match in tournament
	Date       time.Time // date of match
	TeamId1    int64     // id of 1st team
	TeamId2    int64     // id of 2nd team
	Location   string    // match location
	Rule       string    // we use this field to store a specific match rule.
	Result1    int64     // result of 1st team
	Result2    int64     // result of 2nd team
	Finished   bool      // is match finished
	Ready      bool      // is match ready for predictions.
	CanPredict bool      // can user make a prediction (used to block predictions when match has started).
}

// MatchByID gets a Tmatch entity by id.
//
func MatchByID(c appengine.Context, matchID int64) (*Tmatch, error) {
	var m Tmatch
	key := datastore.NewKey(c, "Tmatch", "", matchID, nil)

	if err := datastore.Get(c, key, &m); err != nil {
		log.Errorf(c, "match not found : %v", err)
		return &m, err
	}
	return &m, nil
}

// Matches returns the corresponding array of matches from an array of ids.
//
func Matches(c appengine.Context, matchIds []int64) []*Tmatch {
	var matches []*Tmatch
	for _, matchID := range matchIds {
		m, err := MatchByID(c, matchID)
		if err != nil {
			log.Errorf(c, " Matches, cannot find match with Id=%", matchID)
		} else {
			matches = append(matches, m)
		}
	}
	return matches
}

// GetMatchByIDNumber gets match entity by iDNumber.
//
func GetMatchByIDNumber(c appengine.Context, tournament Tournament, matchInternalID int64) *Tmatch {
	matches1stStage := Matches(c, tournament.Matches1stStage)
	for _, m := range matches1stStage {
		if m.IdNumber == matchInternalID {
			return m
		}
	}
	matches2ndStage := Matches(c, tournament.Matches2ndStage)
	for _, m := range matches2ndStage {
		if m.IdNumber == matchInternalID {
			return m
		}
	}
	return nil
}

// MatchKeyByID returns a pointer to a match key given a match id.
//
func MatchKeyByID(c appengine.Context, id int64) *datastore.Key {
	key := datastore.NewKey(c, "Tmatch", "", id, nil)
	return key
}

// UpdateMatch updates a match.
//
func UpdateMatch(c appengine.Context, m *Tmatch) error {
	k := MatchKeyByID(c, m.Id)
	oldMatch := new(Tmatch)
	if err := datastore.Get(c, k, oldMatch); err == nil {
		if _, err = datastore.Put(c, k, m); err != nil {
			return err
		}
	}
	return nil
}

// UpdateMatches updates an array of matches.
//
func UpdateMatches(c appengine.Context, matches []*Tmatch) error {
	keys := make([]*datastore.Key, len(matches))
	for i := range keys {
		keys[i] = MatchKeyByID(c, matches[i].Id)
	}
	if _, err := datastore.PutMulti(c, keys, matches); err != nil {
		return err
	}
	return nil
}

// DestroyMatches destroys an array of matches.
//
func DestroyMatches(c appengine.Context, matchIds []int64) error {
	keys := make([]*datastore.Key, len(matchIds))
	for i := range keys {
		keys[i] = MatchKeyByID(c, matchIds[i])
	}
	if err := datastore.DeleteMulti(c, keys); err != nil {
		return err
	}
	return nil
}

// SetResults sets results on an array of matches and triggers a match update and group update.
//
func SetResults(c appengine.Context, matches []*Tmatch, results1 []int64, results2 []int64, t *Tournament) error {
	desc := "Set Results:"
	if len(matches) != len(results1) || len(matches) != len(results2) {
		log.Errorf(c, "%s unable to set result on matches", desc)
		return errors.New(helpers.ErrorCodeMatchesCannotUpdate)
	}

	for i, m := range matches {
		if results1[i] < 0 || results2[i] < 0 {
			log.Errorf(c, "%s unable to set result on match with id: %v, %v", desc, m.Id)
			return errors.New(helpers.ErrorCodeMatchCannotUpdate)
		}
		m.Result1 = results1[i]
		m.Result2 = results2[i]
		m.Finished = true
	}

	// batch match update
	if err := UpdateMatches(c, matches); err != nil {
		log.Errorf(c, "%s unable to set results on matches: %v", desc, err)
		return err
	}
	allMatches := GetAllMatchesFromTournament(c, t)
	phases := MatchesGroupByPhase(t, allMatches)

	for _, m := range matches {
		log.Infof(c, "%s Trigger current match: %v", desc, m.Id)

		if ismatch, g := t.IsMatchInGroup(c, m); ismatch == true {
			if err := UpdatePointsAndGoals(c, g, m, t); err != nil {
				log.Errorf(c, "%s Update Points and Goals: unable to update points and goals for group for match with id:%v error: %v", desc, m.IdNumber, err)
				return errors.New(helpers.ErrorCodeMatchCannotUpdate)
			}
			if err := UpdateGroup(c, g); err != nil {
				log.Errorf(c, "%s unable to update group: %v", desc, err)
				return err
			}
		}
		if isLast, phaseID := lastMatchOfPhase(c, m, &phases); isLast == true {
			log.Infof(c, "%s -------------------------------------------------->", desc)
			log.Infof(c, "%s Trigger update of next phase here: next phase: %v", desc, phaseID+1)
			log.Infof(c, "%s Trigger update of next phase here: next phase: %v", desc, m)
			if int(phaseID+1) < len(phases) {
				UpdateNextPhase(c, t, &phases[phaseID], &phases[phaseID+1])
			}
			log.Infof(c, "%s -------------------------------------------------->", desc)
			// update flag first phase complete.
			if phaseID == 0 {
				t.IsFirstStageComplete = true
				t.Update(c)
			}
		}
	}

	log.Infof(c, "%s points and goals updated", desc)

	return nil
}

// SetResult sets the result of a match entity and triggers a match update in datastore and score updates.
//
func SetResult(c appengine.Context, m *Tmatch, result1 int64, result2 int64, t *Tournament) error {

	desc := "Set Result:"
	if result1 < 0 || result2 < 0 {
		log.Errorf(c, "%s unable to set result on match with id: %v", desc, m.Id)
		return errors.New(helpers.ErrorCodeMatchCannotUpdate)
	}
	m.Result1 = result1
	m.Result2 = result2
	m.Finished = true

	var err error
	if err = UpdateMatch(c, m); err != nil {
		log.Errorf(c, "%s unable to set result on match with id: %v, %v", desc, m.Id, err)
		return err
	}

	// update score for all users.
	if err1 := t.UpdateUsersScore(c, m); err1 != nil {
		log.Errorf(c, "%s unable to update users score on match with id: %v, %v", desc, m.Id, err)
	}
	// update score for all teams.
	if err1 := t.UpdateTeamsAccuracy(c, m); err1 != nil {
		log.Errorf(c, "%s unable to update teams score on match with id: %v, %v", desc, m.Id, err)
	}

	if ismatch, g := t.IsMatchInGroup(c, m); ismatch == true {
		if err := UpdatePointsAndGoals(c, g, m, t); err != nil {
			log.Errorf(c, "%s Update Points and Goals: unable to update points and goals for group for match with id:%v error: %v", desc, m.IdNumber, err)
			return errors.New(helpers.ErrorCodeMatchCannotUpdate)
		}
		UpdateGroup(c, g)
	}

	if t.TwoLegged == false {
		allMatches := GetAllMatchesFromTournament(c, t)
		phases := MatchesGroupByPhase(t, allMatches)
		if isLast, phaseID := lastMatchOfPhase(c, m, &phases); isLast == true {
			log.Infof(c, "%s -------------------------------------------------->", desc)
			log.Infof(c, "%s Trigger update of next phase here: next phase: %v", desc, phaseID+1)
			log.Infof(c, "%s Trigger update of next phase here: next phase: %v", desc, m)
			if int(phaseID+1) < len(phases) {
				UpdateNextPhase(c, t, &phases[phaseID], &phases[phaseID+1])
			}
			log.Infof(c, "%s -------------------------------------------------->", desc)
			// update flag first phase complete.
			if phaseID == 0 {
				t.IsFirstStageComplete = true
				t.Update(c)
			}
		}
	}

	return nil
}

// GetAllMatchesFromTournament gets an array of all matches of a tournament.
//
func GetAllMatchesFromTournament(c appengine.Context, tournament *Tournament) []*Tmatch {

	matches := Matches(c, tournament.Matches1stStage)
	matches2ndPhase := Matches(c, tournament.Matches2ndStage)

	// append 2nd round to first one
	for _, m := range matches2ndPhase {
		matches = append(matches, m)
	}

	return matches
}

// GetMatchesByPhase gets all matches of a specific phase.
//
func GetMatchesByPhase(c appengine.Context, t *Tournament, phaseName string) []*Tmatch {

	var tb TournamentBuilder
	if tb = GetTournamentBuilder(t); tb == nil {
		return []*Tmatch{}
	}
	limits := tb.MapOfPhaseIntervals()

	low := limits[phaseName][0]
	high := limits[phaseName][1]

	matches := GetAllMatchesFromTournament(c, t)

	var filteredMatches []*Tmatch
	for i, v := range matches {
		if v.IdNumber >= low && v.IdNumber <= high {
			filteredMatches = append(filteredMatches, matches[i])
		}
	}
	return filteredMatches
}

// MatchesGroupByPhase gets all matches grouped by phases. Returns an array of phases.
//
func MatchesGroupByPhase(t *Tournament, matches []*Tmatch) []Tphase {

	var tb TournamentBuilder
	if tb = GetTournamentBuilder(t); tb == nil {
		return []Tphase{}
	}

	limits := tb.MapOfPhaseIntervals()
	phaseNames := tb.ArrayOfPhases()

	phases := make([]Tphase, len(limits))
	for i := range phases {
		phases[i].Name = phaseNames[i]
		low := limits[phases[i].Name][0]
		high := limits[phases[i].Name][1]

		var filteredMatches []Tmatch
		for _, v := range matches {
			if v.IdNumber >= low && v.IdNumber <= high {
				filteredMatches = append(filteredMatches, *v)
			}
		}
		phases[i].Days = MatchesGroupByDay(filteredMatches)
	}
	return phases
}

// MatchesGroupByDay gets all matches grouped by days. Returns an array of days.
//
func MatchesGroupByDay(matches []Tmatch) []Tday {

	mapOfDays := make(map[string][]Tmatch)

	const shortForm = "Jan/02/2006"
	for _, m := range matches {
		currentDate := m.Date.Format(shortForm)
		_, ok := mapOfDays[currentDate]
		if ok {
			mapOfDays[currentDate] = append(mapOfDays[currentDate], m)
		} else {
			var arrayMatches []Tmatch
			arrayMatches = append(arrayMatches, m)
			mapOfDays[currentDate] = arrayMatches
		}
	}

	var days []Tday
	days = make([]Tday, len(mapOfDays))
	i := 0
	for key, value := range mapOfDays {
		days[i].Date, _ = time.Parse(shortForm, key)
		days[i].Matches = value
		i++
	}
	sort.Sort(ByDate(days))
	return days
}

// OldMatches gets the number of matches in a tournament that are finished.
//
func (t *Tournament) OldMatches(c appengine.Context) int {
	matches := GetAllMatchesFromTournament(c, t)
	old := 0
	for _, m := range matches {
		if m.Finished {
			old++
		}
	}
	return old
}
