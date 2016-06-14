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
	"fmt"
	"strconv"
	"time"

	"appengine"
	"appengine/datastore"

	"github.com/taironas/gonawin/helpers/log"
)

// CopaAmericaTournament2016 is a placeholder for the Copa America Tournament.
//
type CopaAmericaTournament2016 struct {
}

// MapOfGroups represents the groupsof a tournament, key: group name, value: string array of teams.
//
func (cat CopaAmericaTournament2016) MapOfGroups() map[string][]string {
	var mapCAGroups map[string][]string
	mapCAGroups = make(map[string][]string)
	mapCAGroups["A"] = []string{"United States", "Colombia", "Costa Rica", "Paraguay"}
	mapCAGroups["B"] = []string{"Peru", "Ecuador", "Brazil", "Haiti"}
	mapCAGroups["C"] = []string{"Mexico", "Venezuela", "Uruguay", "Jamaica"}
	mapCAGroups["D"] = []string{"Argentina", "Chile", "Panama", "Bolivia"}

	return mapCAGroups
}

// MapOfTeamCodes is the map of country codes, key: team name, value: ISO code
// example: Brazil: BR
//
func (cat CopaAmericaTournament2016) MapOfTeamCodes() map[string]string {

	var codes map[string]string
	codes = make(map[string]string)

	codes["United States"] = "us"
	codes["Colombia"] = "co"
	codes["Costa Rica"] = "cr"
	codes["Paraguay"] = "py"

	codes["Peru"] = "pe"
	codes["Ecuador"] = "ec"
	codes["Brazil"] = "br"
	codes["Haiti"] = "ht"

	codes["Mexico"] = "mx"
	codes["Venezuela"] = "ve"
	codes["Uruguay"] = "uy"
	codes["Jamaica"] = "jm"

	codes["Argentina"] = "ar"
	codes["Chile"] = "ch"
	codes["Panama"] = "pa"
	codes["Bolivia"] = "bo"

	return codes
}

// MapOfGroupMatches return a map with the maches of all the groups
// key: group name, value: array of array of strings with match information:
// MatchId, MatchDate, MatchTeam1, MatchTeam2, MatchLocation
//
// Example:
//
// 	Group A:[{"1", "Jun/12/2014", "Brazil", "Croatia", "Arena de São Paulo, São Paulo"}, ...]
//
func (cat CopaAmericaTournament2016) MapOfGroupMatches() map[string][][]string {

	mapGroupMatches := make(map[string][][]string)

	const (
		cMatchID       = 0
		cMatchDate     = 1
		cMatchTeam1    = 2
		cMatchTeam2    = 3
		cMatchLocation = 4
	)

	mA1 := []string{"1", "Jun/03/2016", "United States", "Colombia", "Levi's Stadium, Santa Clara"}
	mA2 := []string{"2", "Jun/04/2016", "Costa Rica", "Paraguay", "Camping World Stadium, Orlando"}
	mA3 := []string{"9", "Jun/07/2016", "United States", "Costa Rica", "Soldier Field, Chicago"}
	mA4 := []string{"10", "Jun/07/2016", "Colombia", "Paraguay", "Rose Bowl, Pasadena"}
	mA5 := []string{"17", "Jun/11/2016", "United States", "Paraguay", "Lincoln Financial Field, Philadelphia"}
	mA6 := []string{"18", "Jun/11/2016", "Colombia", "Costa Rica", "NRG Stadium, Houston"}

	mB1 := []string{"3", "Jun/04/2016", "Haiti", "Peru", "CenturyLink Field, Seattle"}
	mB2 := []string{"4", "Jun/04/2016", "Brazil", "Ecuador", "Rose Bowl, Pasadena"}
	mB3 := []string{"11", "Jun/08/2016", "Brazil", "Haiti", "Camping World Stadium, Orlando"}
	mB4 := []string{"12", "Jun/08/2016", "Ecuador", "Peru", "University of Phoenix Stadium, Glendale"}
	mB5 := []string{"19", "Jun/12/2016", "Ecuador", "Haiti", "MetLife Stadium, East Rutherford"}
	mB6 := []string{"20", "Jun/12/2016", "Brazil", "Peru", "Gillette Stadium, Foxborough"}

	mC1 := []string{"5", "Jun/05/2016", "Jamaica", "Venezuela", "Soldier Field, Chicago"}
	mC2 := []string{"6", "Jun/05/2016", "Mexico", "Uruguay", "University of Phoenix Stadium, Glendale"}
	mC3 := []string{"13", "Jun/09/2016", "Uruguay", "Venezuela", "Lincoln Financial Field, Philadelphia"}
	mC4 := []string{"14", "Jun/09/2016", "Mexico", "Jamaica", "Rose Bowl, Pasadena"}
	mC5 := []string{"21", "Jun/13/2016", "Mexico", "Venezuela", "NRG Stadium, Houston"}
	mC6 := []string{"22", "Jun/13/2016", "Uruguay", "Jamaica", "Levi's Stadium, Santa Clara"}

	mD1 := []string{"7", "Jun/06/2016", "Panama", "Bolivia", "Camping World Stadium, Orlando"}
	mD2 := []string{"8", "Jun/06/2016", "Argentina", "Chile", "Levi's Stadium, Santa Clara"}
	mD3 := []string{"15", "Jun/10/2016", "Chile", "Bolivia", "Gillette Stadium, Foxborough"}
	mD4 := []string{"16", "Jun/10/2016", "Argentina", "Panama", "Soldier Field, Chicago"}
	mD5 := []string{"23", "Jun/14/2016", "Chile", "Panama", "Lincoln Financial Field, Philadelphia"}
	mD6 := []string{"24", "Jun/14/2016", "Argentina", "Bolivia", "CenturyLink Field, Seattle"}

	var matchesA [][]string
	var matchesB [][]string
	var matchesC [][]string
	var matchesD [][]string

	matchesA = append(matchesA, mA1, mA2, mA3, mA4, mA5, mA6)
	matchesB = append(matchesB, mB1, mB2, mB3, mB4, mB5, mB6)
	matchesC = append(matchesC, mC1, mC2, mC3, mC4, mC5, mC6)
	matchesD = append(matchesD, mD1, mD2, mD3, mD4, mD5, mD6)

	mapGroupMatches["A"] = matchesA
	mapGroupMatches["B"] = matchesB
	mapGroupMatches["C"] = matchesC
	mapGroupMatches["D"] = matchesD

	return mapGroupMatches
}

// MapOf2ndRoundMatches returns the map of 2nd round matches, of the tournament.
// key: round number, value: array of array of strings with match information:
// MatchId, MatchDate, MatchTeam1, MatchTeam2, MatchLocation
//
// Example:
//
//	round 16:[{"1", "Jun/12/2014", "Brazil", "Croatia", "Arena de São Paulo, São Paulo"}, ...]
//
func (cat CopaAmericaTournament2016) MapOf2ndRoundMatches() map[string][][]string {

	// 17 Quarter-finals
	m2nd1 := []string{"25", "Jun/16/2016", "1A", "2B", "CenturyLink Field, Seattle"}
	m2nd2 := []string{"26", "Jun/17/2016", "1B", "2A", "MetLife Stadium, East Rutherford"}
	m2nd3 := []string{"27", "Jun/18/2016", "1D", "2C", "Gillette Stadium, Foxborough"}
	m2nd4 := []string{"28", "Jun/18/2016", "1C", "2D", "Levi's Stadium, Santa Clara"}
	// 18 Semi-finals
	m2nd5 := []string{"29", "Jun/21/2016", "W25", "W27", "NRG Stadium, Houston"}
	m2nd6 := []string{"30", "Jun/22/2016", "W26", "W28", "Soldier Field, Chicago"}
	//19 Round 19  -  Match for third place
	m2nd7 := []string{"31", "Jun/25/2016", "L29", "L30", "University of Phoenix Stadium, Glendale"}
	//"20" Final
	m2nd8 := []string{"32", "Jun/26/2016", "W29", "W30", "MetLife Stadium, East Rutherford"}

	var round17 [][]string
	var round18 [][]string
	var round19 [][]string
	var round20 [][]string

	round17 = append(round17, m2nd1, m2nd2, m2nd3, m2nd4)
	round18 = append(round18, m2nd5, m2nd6)
	round19 = append(round19, m2nd7)
	round20 = append(round20, m2nd8)

	mapMatches2ndRound := make(map[string][][]string)
	mapMatches2ndRound["17"] = round17
	mapMatches2ndRound["18"] = round18
	mapMatches2ndRound["19"] = round19
	mapMatches2ndRound["20"] = round20

	return mapMatches2ndRound
}

// ArrayOfPhases returns an array of the phases names of world cup tournament:
// FirstStage, RoundOf16, QuarterFinals, SemiFinals, ThirdPlace, Finals
//
func (cat CopaAmericaTournament2016) ArrayOfPhases() []string {
	return []string{cFirstStage, cQuarterFinals, cSemiFinals, cThirdPlace, cFinals}
}

// MapOfPhaseIntervals returns a map with key the corresponding phase in the copa america tournament
// at value a tuple that represent the match number interval in which the phase take place:
// first stage: matches 1 to 48
// Round of 16: matches 49 to 56
// Quarte-finals: matches 57 to 60
// Semi-finals: matches 61 to 62
// Third Place: match 63
// Finals: match 64
func (cat CopaAmericaTournament2016) MapOfPhaseIntervals() map[string][]int64 {

	limits := make(map[string][]int64)
	limits[cFirstStage] = []int64{1, 24}
	limits[cQuarterFinals] = []int64{25, 28}
	limits[cSemiFinals] = []int64{29, 30}
	limits[cThirdPlace] = []int64{31, 31}
	limits[cFinals] = []int64{32, 32}
	return limits
}

// CreateCopaAmerica create Copa America tournament entity 2016.
//
func CreateCopaAmerica2016(c appengine.Context, adminID int64) (*Tournament, error) {
	desc := "Copa America"
	// create new tournament
	log.Infof(c, "%s: start", desc)

	cat := CopaAmericaTournament2016{}

	mapCATGroups := cat.MapOfGroups()
	mapCountryCodes := cat.MapOfTeamCodes()
	mapTeamID := make(map[string]int64)

	// mapGroupMatches is a map where the key is a string which represent the group
	// the key is a two dimensional string array. each element in the array represent a specific field in the match
	// map[string][][]string
	// example: {"1", "Jun/12/2014", "Brazil", "Croatia", "Arena de São Paulo, São Paulo"}
	mapGroupMatches := cat.MapOfGroupMatches()

	const (
		cMatchID       = 0
		cMatchDate     = 1
		cMatchTeam1    = 2
		cMatchTeam2    = 3
		cMatchLocation = 4
	)

	// for date parsing
	const shortForm = "Jan/02/2006"

	// matches1stStageIds is an array of  int64
	// where we allocate IDs of the Tmatches entities
	// we will store them in the tournament entity for easy retreival later on.
	matchesA := mapGroupMatches["A"]
	matches1stStageIds := make([]int64, len(mapGroupMatches)*len(matchesA))

	log.Infof(c, "%s: maps ready", desc)

	// build groups, and teams
	groups := make([]Tgroup, len(mapCATGroups))
	for groupName, teams := range mapCATGroups {
		log.Infof(c, "---------------------------------------")
		log.Infof(c, "%s: working with group: %v", desc, groupName)
		log.Infof(c, "%s: teams: %v", desc, teams)

		var group Tgroup
		group.Name = groupName
		groupIndex := int64(groupName[0]) - 65
		group.Teams = make([]Tteam, len(teams))
		group.Points = make([]int64, len(teams))
		group.GoalsF = make([]int64, len(teams))
		group.GoalsA = make([]int64, len(teams))

		for i, teamName := range teams {
			log.Infof(c, "%s: team: %v", desc, teamName)

			teamID, _, err1 := datastore.AllocateIDs(c, "Tteam", nil, 1)
			if err1 != nil {
				return nil, err1
			}
			log.Infof(c, "%s: team: %v allocateIDs ok", desc, teamName)

			teamkey := datastore.NewKey(c, "Tteam", "", teamID, nil)
			log.Infof(c, "%s: team: %v NewKey ok", desc, teamName)

			team := &Tteam{teamID, teamName, mapCountryCodes[teamName]}
			log.Infof(c, "%s: team: %v instance of team ok", desc, teamName)

			_, err := datastore.Put(c, teamkey, team)
			if err != nil {
				return nil, err
			}
			log.Infof(c, "%s: team: %v put in datastore ok", desc, teamName)
			group.Teams[i] = *team
			mapTeamID[teamName] = teamID
		}

		// build group matches:
		log.Infof(c, "%s: building group matches", desc)

		groupMatches := mapGroupMatches[groupName]
		group.Matches = make([]Tmatch, len(groupMatches))

		for matchIndex, matchData := range groupMatches {
			log.Infof(c, "%s: match data: %v", desc, matchData)

			matchID, _, err1 := datastore.AllocateIDs(c, "Tmatch", nil, 1)
			if err1 != nil {
				return nil, err1
			}

			log.Infof(c, "%s: match: %v allocateIDs ok", desc, matchID)

			matchkey := datastore.NewKey(c, "Tmatch", "", matchID, nil)
			log.Infof(c, "%s: match: new key ok", desc)

			matchTime, _ := time.Parse(shortForm, matchData[cMatchDate])
			matchInternalID, _ := strconv.Atoi(matchData[cMatchID])
			emptyrule := ""
			emptyresult := int64(0)
			match := &Tmatch{
				matchID,
				int64(matchInternalID),
				matchTime,
				mapTeamID[matchData[cMatchTeam1]],
				mapTeamID[matchData[cMatchTeam2]],
				matchData[cMatchLocation],
				emptyrule,
				emptyresult,
				emptyresult,
				false,
				true,
				true,
			}

			log.Infof(c, "%s: match: build match ok", desc)

			_, err := datastore.Put(c, matchkey, match)
			if err != nil {
				return nil, err
			}
			log.Infof(c, "%s: match: %v put in datastore ok", desc, matchData)
			group.Matches[matchIndex] = *match

			// save in an array of int64 all the allocate IDs to store them in the tournament for easy retreival later on.
			matches1stStageIds[int64(matchInternalID)-1] = matchID
		}

		groupID, _, err1 := datastore.AllocateIDs(c, "Tgroup", nil, 1)
		if err1 != nil {
			return nil, err1
		}

		log.Infof(c, "%s: Group: %v allocate Id ok", desc, groupName)

		groupkey := datastore.NewKey(c, "Tgroup", "", groupID, nil)
		log.Infof(c, "%s: Group: %v New Key ok", desc, groupName)

		group.Id = groupID
		groups[groupIndex] = group
		_, err := datastore.Put(c, groupkey, &group)
		if err != nil {
			return nil, err
		}
		log.Infof(c, "%s: Group: %v put in datastore ok", desc, groupName)
	}

	// build array of group ids
	groupIds := make([]int64, len(groups))
	for i := range groupIds {
		groupIds[i] = groups[i].Id
	}

	log.Infof(c, "%s: build of groups ids complete: %v", desc, groupIds)

	// matches 2nd stage
	var matches2ndStageIds []int64
	var userIds []int64
	var teamIds []int64

	// mapMatches2ndRound  is a map where the key is a string which represent the rounds
	// the key is a two dimensional string array. each element in the array represent a specific field in the match
	// mapMatches2ndRound is a map[string][][]string
	// example: {"64", "Jul/13/2014", "W61", "W62", "Rio de Janeiro"}
	mapMatches2ndRound := cat.MapOf2ndRoundMatches()

	// build matches 2nd phase
	for roundNumber, roundMatches := range mapMatches2ndRound {
		log.Infof(c, "%s: building 2nd round matches: round number %v", desc, roundNumber)
		for _, matchData := range roundMatches {
			log.Infof(c, "%s: second phase match data: %v", desc, matchData)

			matchID, _, err1 := datastore.AllocateIDs(c, "Tmatch", nil, 1)
			if err1 != nil {
				return nil, err1
			}

			log.Infof(c, "%s: match: %v allocateIDs ok", desc, matchID)

			matchkey := datastore.NewKey(c, "Tmatch", "", matchID, nil)
			log.Infof(c, "%s: match: new key ok", desc)

			matchTime, _ := time.Parse(shortForm, matchData[cMatchDate])
			matchInternalID, _ := strconv.Atoi(matchData[cMatchID])

			rule := fmt.Sprintf("%s %s", matchData[cMatchTeam1], matchData[cMatchTeam2])
			emptyresult := int64(0)
			match := &Tmatch{
				matchID,
				int64(matchInternalID),
				matchTime,
				0, // second round matches start with ids at 0
				0, // second round matches start with ids at 0
				matchData[cMatchLocation],
				rule,
				emptyresult,
				emptyresult,
				false,
				false,
				true,
			}
			log.Infof(c, "%s: match 2nd round: build match ok", desc)

			_, err := datastore.Put(c, matchkey, match)
			if err != nil {
				return nil, err
			}
			log.Infof(c, "%s: 2nd round match: %v put in datastore ok", desc, matchData)
			// save in an array of int64 all the allocate IDs to store them in the tournament for easy retreival later on.
			matches2ndStageIds = append(matches2ndStageIds, matchID)
		}
	}

	tstart, _ := time.Parse(shortForm, "Jun/03/2016")
	tend, _ := time.Parse(shortForm, "Jun/29/2016")
	adminIds := make([]int64, 1)
	adminIds[0] = adminID
	name := "2016 Copa America"
	description := "Centenario"
	var tournament *Tournament
	var err error
	if tournament, err = CreateTournament(c, name, description, tstart, tend, adminID); err != nil {
		log.Infof(c, "%s: something went wrong when creating tournament.", desc)
	} else {
		tournament.GroupIds = groupIds
		tournament.Matches1stStage = matches1stStageIds
		tournament.Matches2ndStage = matches2ndStageIds
		tournament.UserIds = userIds
		tournament.TeamIds = teamIds
		tournament.IsFirstStageComplete = false
		if err1 := tournament.Update(c); err1 != nil {
			log.Infof(c, "%s: unable to udpate tournament.", desc)
		}
	}

	log.Infof(c, "%s: instance of tournament ready", desc)

	return tournament, nil
}

// MapOfIDTeams returns a map of team IDs as keys and team names as values.
//
func (cat CopaAmericaTournament2016) MapOfIDTeams(c appengine.Context, tournament *Tournament) map[int64]string {

	mapIDTeams := make(map[int64]string)

	groups := Groups(c, tournament.GroupIds)
	for _, g := range groups {
		for _, t := range g.Teams {
			mapIDTeams[t.Id] = t.Name
		}
	}
	return mapIDTeams
}
