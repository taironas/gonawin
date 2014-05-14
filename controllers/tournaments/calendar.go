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

package tournaments

import (
	"errors"
	"net/http"
	"sort"
	"time"

	"appengine"

	"github.com/santiaago/gonawin/helpers"
	"github.com/santiaago/gonawin/helpers/handlers"
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

// A PhaseJson is a variable to hold a the name of a phase and an array of days.
// We use it to group tournament matches information by phases.
type PhaseJson struct {
	Name string
	Days []DayJson
	Completed bool
}

// Json tournament calendar handler:
// Use this handler to get the calendar of a tournament.
// The calendar structure is an array of the tournament matches with the following information:
// * the location
// * the teams involved
// * the date
// by default the data returned is grouped by days.This means we will return an array of days, each of which can have an array of matches.
// You can also specify the 'groupby' parameter to be 'day' or 'phase' in which case you would have an array of phases,
// each of which would have an array of days who would have an array of matches.
func Calendar(w http.ResponseWriter, r *http.Request, u *mdl.User) error {
	c := appengine.NewContext(r)
	desc := "Tournament Calendar Handler:"

	if r.Method == "GET" {
		tournamentId, err := handlers.PermalinkID(r, c, 3)
		if err != nil {
			log.Errorf(c, "%s error extracting permalink err:%v", desc, err)
			return &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeTournamentNotFound)}
		}

		var t *mdl.Tournament
		t, err = mdl.TournamentById(c, tournamentId)
		if err != nil {
			log.Errorf(c, "%s tournament with id:%v was not found %v", desc, tournamentId, err)
			return &helpers.NotFound{Err: errors.New(helpers.ErrorCodeTournamentNotFound)}
		}

		groupby := r.FormValue("groupby")
		// if wrong data we set groupby to "day"
		if groupby != "day" && groupby != "phase" {
			groupby = "day"
		}

		if groupby == "day" {
			log.Infof(c, "%s ready to build days array", desc)
			matchesJson := buildMatchesFromTournament(c, t, u)

			days := matchesGroupByDay(matchesJson)

			data := struct {
				Days []DayJson
			}{
				days,
			}

			return templateshlp.RenderJson(w, c, data)

		} else if groupby == "phase" {
			log.Infof(c, "%s ready to build phase array", desc)
			matchesJson := buildMatchesFromTournament(c, t, u)
			phases := matchesGroupByPhase(matchesJson)

			data := struct {
				Phases []PhaseJson
			}{
				phases,
			}
			return templateshlp.RenderJson(w, c, data)
		}
	}
	return &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeNotSupported)}
}

// From an array of Matches, create an array of Phases where the matches are grouped in.
// We use the Phases intervals and the IdNumber of each match to do this operation.
func matchesGroupByPhase(matches []MatchJson) []PhaseJson {
	limits := mdl.MapOfPhaseIntervals()
	phaseNames := mdl.ArrayOfPhases()

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
		phases[i].Days = matchesGroupByDay(filteredMatches)
		lastDayOfPhase := len(phases[i].Days) - 1
		lastMatchOfPhase := len(phases[i].Days[lastDayOfPhase].Matches) - 1
		if phases[i].Days[lastDayOfPhase].Matches[lastMatchOfPhase].Finished{
			phases[i].Completed = true
		}else{
			phases[i].Completed = false
		}
	}
	return phases
}

// From an array of matches, create an array of Days where the matches are grouped in.
// We use the Date of each match to do this.
func matchesGroupByDay(matches []MatchJson) []DayJson {

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
