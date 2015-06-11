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

	"github.com/santiaago/gonawin/extract"
	"github.com/santiaago/gonawin/helpers"
	"github.com/santiaago/gonawin/helpers/log"
	templateshlp "github.com/santiaago/gonawin/helpers/templates"

	mdl "github.com/santiaago/gonawin/models"
)

// A DayJson is a variable to hold a date and match field.
// We use it to group tournament matches information by days.
type DayJson struct {
	Date    time.Time
	Matches []MatchJson
}

type DayWithPredictionJson struct {
	Date    time.Time
	Matches []MatchWithPredictionJson
}

type MatchWithPredictionJson struct {
	Match        MatchJson
	Participants []UserPredictionJson
}

type UserPredictionJson struct {
	Id       int64
	Username string
	Alias    string
	Predict  string
}

// A PhaseJson is a variable to hold a the name of a phase and an array of days.
// We use it to group tournament matches information by phases.
type PhaseJson struct {
	Name      string
	Days      []DayJson
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
		log.Infof(c, "%s ready to build days array", desc)
		matchesJson := buildMatchesFromTournament(c, t, u)

		days := matchesGroupByDay(t, matchesJson)

		data := struct {
			Days []DayJson
		}{
			days,
		}

		return templateshlp.RenderJson(w, c, data)

	} else if groupby == "phase" {
		log.Infof(c, "%s ready to build phase array", desc)
		matchesJson := buildMatchesFromTournament(c, t, u)
		phases := matchesGroupByPhase(t, matchesJson)
		data := struct {
			Phases []PhaseJson
		}{
			phases,
		}
		return templateshlp.RenderJson(w, c, data)
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
		predicts := mdl.PredictsByIds(c, p.PredictIds)
		predictsByPlayer[i] = predicts
	}

	groupby := r.FormValue("groupby")
	// if wrong data we set groupby to "day"
	if groupby != "day" && groupby != "phase" {
		groupby = "day"
	}

	if groupby == "day" {
		log.Infof(c, "%s ready to build days array", desc)
		matchesJson := buildMatchesFromTournament(c, t, u)

		days := matchesGroupByDay(t, matchesJson)
		dayswp := make([]DayWithPredictionJson, len(days))
		for i, d := range days {
			dayswp[i].Date = d.Date
			matcheswp := make([]MatchWithPredictionJson, len(d.Matches))
			for j, m := range d.Matches {
				matcheswp[j].Match = m
				uwp := make([]UserPredictionJson, len(players))
				for k, p := range players {
					// get player prediction
					uwp[k].Id = p.Id
					uwp[k].Username = p.Username
					uwp[k].Alias = p.Alias
					if hasMatch, l := predictsByPlayer[k].ContainsMatchId(m.Id); hasMatch == true {
						uwp[k].Predict = fmt.Sprintf("%v - %v", predictsByPlayer[k][l].Result1, predictsByPlayer[k][l].Result2)
					} else {
						uwp[k].Predict = "-"
					}
				}
				matcheswp[j].Participants = uwp
			}
			dayswp[i].Matches = matcheswp
		}

		data := struct {
			Days []DayWithPredictionJson
		}{
			dayswp,
		}

		return templateshlp.RenderJson(w, c, data)

	} else if groupby == "phase" {
		// @santiaago: right now not supported.
	}

	return &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeNotSupported)}
}

// From an array of Matches, create an array of Phases where the matches are grouped in.
// We use the Phases intervals and the IdNumber of each match to do this operation.
func matchesGroupByPhase(t *mdl.Tournament, matches []MatchJson) []PhaseJson {

	var tb mdl.TournamentBuilder
	if tb = mdl.GetTournamentBuilder(t); tb == nil {
		return []PhaseJson{}
	}

	limits := tb.MapOfPhaseIntervals()
	phaseNames := tb.ArrayOfPhases()

	phases := make([]PhaseJson, len(limits))
	for i, _ := range phases {
		phases[i].Name = phaseNames[i]
		low := limits[phases[i].Name][0]
		high := limits[phases[i].Name][1]

		var filteredMatches []MatchJson
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
func matchesGroupByDay(t *mdl.Tournament, matches []MatchJson) []DayJson {

	mapOfDays := make(map[string][]MatchJson)

	const shortForm = "Jan/02/2006"
	for _, m := range matches {
		currentDate := m.Date.Format(shortForm)
		_, ok := mapOfDays[currentDate]
		if ok {
			mapOfDays[currentDate] = append(mapOfDays[currentDate], m)
		} else {
			var arrayMatches []MatchJson
			arrayMatches = append(arrayMatches, m)
			mapOfDays[currentDate] = arrayMatches
		}
	}

	var days []DayJson
	days = make([]DayJson, len(mapOfDays))
	i := 0
	for key, value := range mapOfDays {
		days[i].Date, _ = time.Parse(shortForm, key)
		days[i].Matches = value
		i++
	}

	sort.Sort(ByDate(days))
	return days
}

// ByDate type implements the sort.Interface for []DayJson based on the date field.
type ByDate []DayJson

func (a ByDate) Len() int           { return len(a) }
func (a ByDate) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByDate) Less(i, j int) bool { return a[i].Date.Before(a[j].Date) }
