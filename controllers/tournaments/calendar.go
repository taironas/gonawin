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

package tournaments

import (
	"errors"
	"fmt"
	"net/http"
	"sort"
	"time"

	"appengine"

	"github.com/taironas/gonawin/extract"
	"github.com/taironas/gonawin/helpers"
	"github.com/taironas/gonawin/helpers/log"
	templateshlp "github.com/taironas/gonawin/helpers/templates"

	mdl "github.com/taironas/gonawin/models"
)

// DayJSON is a variable to hold a date and match field.
// We use it to group tournament matches information by days.
//
type DayJSON struct {
	Date    time.Time
	Matches []MatchJSON
}

// DayWithPredictionJSON is a variable to hold a date and match with prediction.
//
type DayWithPredictionJSON struct {
	Date    time.Time
	Matches []MatchWithPredictionJSON
}

// MatchWithPredictionJSON is a variable to hold a match and an array of participants
// who put a predict on the match.
//
type MatchWithPredictionJSON struct {
	Match        MatchJSON
	Participants []UserPredictionJSON
}

// UserPredictionJSON is a variable to hold user data and a prediction.
//
type UserPredictionJSON struct {
	Id       int64 `json:"Id,omitempty"`
	Username string
	Alias    string
	Predict  string
}

// PhaseJSON is a variable to hold a the name of a phase and an array of days.
// We use it to group tournament matches information by phases.
//
type PhaseJSON struct {
	Name      string
	Days      []DayJSON
	Completed bool
}

// Calendar handler gets you the calendar data of a specific tournament.
// Use this handler to get the calendar of a tournament.
// The calendar structure is an array of the tournament matches with the following information:
// * the location
// * the teams involved
// * the date
// by default the data returned is grouped by days.This means we will return an array of days, each of which can have an array of matches.
// You can also specify the 'groupby' parameter to be 'day' or 'phase' in which case you would have an array of phases,
// each of which would have an array of days who would have an array of matches.
//
func Calendar(w http.ResponseWriter, r *http.Request, u *mdl.User) error {
	if r.Method != "GET" {
		return &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeNotSupported)}
	}

	c := appengine.NewContext(r)
	desc := "Tournament Calendar Handler:"

	extract := extract.NewContext(c, desc, r)

	var err error
	var t *mdl.Tournament

	if t, err = extract.Tournament(); err != nil {
		return err
	}

	groupby := r.FormValue("groupby")
	// if wrong data we set groupby to "day"
	if groupby != "day" && groupby != "phase" {
		groupby = "day"
	}

	if groupby == "day" {
		matchesJSON := buildMatchesFromTournament(c, t, u)

		days := matchesGroupByDay(t, matchesJSON)

		data := struct {
			Days []DayJSON
		}{
			days,
		}

		return templateshlp.RenderJSON(w, c, data)

	} else if groupby == "phase" {
		matchesJSON := buildMatchesFromTournament(c, t, u)
		phases := matchesGroupByPhase(t, matchesJSON)
		data := struct {
			Phases []PhaseJSON
		}{
			phases,
		}
		return templateshlp.RenderJSON(w, c, data)
	}
	return &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeNotSupported)}
}

// CalendarWithPrediction handler give you a JSON tournament calendar with the user predictions.
// Use this handler to get the calendar of a tournament with the predictions of the players in a specific team.
// The calendar structure is an array of the tournament matches with the following information:
// * the location
// * the teams involved
// * the date
// by default the data returned is grouped by days.This means we will return an array of days, each of which can have an array of matches.
// the 'groupby' parameter does not support 'phases' yet.
func CalendarWithPrediction(w http.ResponseWriter, r *http.Request, u *mdl.User) error {
	if r.Method != "GET" {
		return &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeNotSupported)}
	}

	c := appengine.NewContext(r)
	desc := "Tournament Calendar with prediction Handler:"
	extract := extract.NewContext(c, desc, r)

	var err error
	var t *mdl.Tournament

	if t, err = extract.Tournament(); err != nil {
		return err
	}

	var team *mdl.Team
	if team, err = extract.Team(); err != nil {
		return err
	}

	var players []*mdl.User
	if players, err = team.Players(c); err != nil {
		return &helpers.InternalServerError{Err: errors.New(helpers.ErrorCodeInternal)}
	}

	predictsByPlayer := make([]mdl.Predicts, len(players))
	for i, p := range players {
		var predicts []*mdl.Predict
		if predicts, err = mdl.PredictsByIds(c, p.PredictIds); err != nil {
			log.Infof(c, "%v something failed when calling PredictsByIds for player %v : %v", desc, p.Id, err)
			continue
		}
		predictsByPlayer[i] = predicts
	}

	groupby := r.FormValue("groupby")
	// if wrong data we set groupby to "day"
	if groupby != "day" && groupby != "phase" {
		groupby = "day"
	}

	if groupby == "day" {
		vm := buildTournamentCalendarViewModel(c, t, u, predictsByPlayer, players)
		return templateshlp.RenderJSON(w, c, vm)

	} else if groupby == "phase" {
		// @taironas: right now not supported.
	}

	return &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeNotSupported)}
}

type tournamentCalendarViewModel struct {
	Days []DayWithPredictionJSON
}

func buildTournamentCalendarViewModel(c appengine.Context, t *mdl.Tournament, u *mdl.User, predictsByPlayer []mdl.Predicts, players []*mdl.User) tournamentCalendarViewModel {

	matches := buildMatchesFromTournament(c, t, u)
	matchesByDay := matchesGroupByDay(t, matches)

	daysWithPredictions := make([]DayWithPredictionJSON, len(matchesByDay))

	for i, day := range matchesByDay {
		daysWithPredictions[i].Date = day.Date
		matchesWithPredictions := matchesWithPredictions(day, players, predictsByPlayer)
		daysWithPredictions[i].Matches = matchesWithPredictions
	}
	return tournamentCalendarViewModel{daysWithPredictions}
}

func matchesWithPredictions(day DayJSON, players []*mdl.User, predictsByPlayer []mdl.Predicts) []MatchWithPredictionJSON {

	matchesWithPredictions := make([]MatchWithPredictionJSON, len(day.Matches))

	for i, m := range day.Matches {
		matchesWithPredictions[i].Match = m
		participants := matchParticipants(m, players, predictsByPlayer)
		matchesWithPredictions[i].Participants = participants
	}
	return matchesWithPredictions
}

func matchParticipants(m MatchJSON, players []*mdl.User, predictsByPlayer []mdl.Predicts) []UserPredictionJSON {

	participants := make([]UserPredictionJSON, len(players))
	for i, p := range players {
		participants[i].Id = p.Id
		participants[i].Username = p.Username
		participants[i].Alias = p.Alias
		prediction := "-"
		if ok, index := predictsByPlayer[i].ContainsMatchID(m.Id); ok {
			prediction = fmt.Sprintf("%v - %v", predictsByPlayer[i][index].Result1, predictsByPlayer[i][index].Result2)
		}
		participants[i].Predict = prediction
	}
	return participants
}

// From an array of Matches, create an array of Phases where the matches are grouped in.
// We use the Phases intervals and the IdNumber of each match to do this operation.
func matchesGroupByPhase(t *mdl.Tournament, matches []MatchJSON) []PhaseJSON {

	var tb mdl.TournamentBuilder
	if tb = mdl.GetTournamentBuilder(t); tb == nil {
		return []PhaseJSON{}
	}

	limits := tb.MapOfPhaseIntervals()
	phaseNames := tb.ArrayOfPhases()

	phases := make([]PhaseJSON, len(limits))
	for i := range phases {
		phases[i].Name = phaseNames[i]
		low := limits[phases[i].Name][0]
		high := limits[phases[i].Name][1]

		var filteredMatches []MatchJSON
		for _, v := range matches {
			if v.IdNumber >= low && v.IdNumber <= high {
				filteredMatches = append(filteredMatches, v)
			}
		}
		phases[i].Days = matchesGroupByDay(t, filteredMatches)
		lastDayOfPhase := len(phases[i].Days) - 1
		lastMatchOfPhase := len(phases[i].Days[lastDayOfPhase].Matches) - 1
		if phases[i].Days[lastDayOfPhase].Matches[lastMatchOfPhase].Finished {
			phases[i].Completed = true
		} else {
			phases[i].Completed = false
		}
	}
	return phases
}

// From an array of matches, create an array of Days where the matches are grouped in.
// We use the Date of each match to do this.
func matchesGroupByDay(t *mdl.Tournament, matches []MatchJSON) []DayJSON {

	mapOfDays := make(map[string][]MatchJSON)

	const shortForm = "Jan/02/2006"
	for _, m := range matches {
		currentDate := m.Date.Format(shortForm)
		_, ok := mapOfDays[currentDate]
		if ok {
			mapOfDays[currentDate] = append(mapOfDays[currentDate], m)
		} else {
			var arrayMatches []MatchJSON
			arrayMatches = append(arrayMatches, m)
			mapOfDays[currentDate] = arrayMatches
		}
	}

	var days []DayJSON
	days = make([]DayJSON, len(mapOfDays))
	i := 0
	for key, value := range mapOfDays {
		days[i].Date, _ = time.Parse(shortForm, key)
		days[i].Matches = value
		i++
	}

	sort.Sort(ByDate(days))
	return days
}

// ByDate type implements the sort.Interface for []DayJSON based on the date field.
type ByDate []DayJSON

func (a ByDate) Len() int           { return len(a) }
func (a ByDate) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByDate) Less(i, j int) bool { return a[i].Date.Before(a[j].Date) }
