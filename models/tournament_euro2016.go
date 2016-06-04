/*
 * Copyright (c) 2016 Santiago Arias | Remy Jourde
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

// EuroTournament2016 represents a Euro tournament.
//
type EuroTournament2016 struct {
}

// MapOfGroups is a map containing the groups of a tournament.
// key: group name, value: string array of teams.
//
func (et EuroTournament2016) MapOfGroups() map[string][]string {
	var mapEuroGroups map[string][]string
	mapEuroGroups = make(map[string][]string)
	mapEuroGroups["A"] = []string{"Albania", "France", "Romania", "Switzerland"}
	mapEuroGroups["B"] = []string{"England", "Russia", "Slovakia", "Wales"}
	mapEuroGroups["C"] = []string{"Germany", "Northern Ireland", "Poland", "Ukraine"}
	mapEuroGroups["D"] = []string{"Croatia", "Czech Republic", "Spain", "Turkey"}
	mapEuroGroups["E"] = []string{"Belgium", "Italy", "Republic of Ireland", "Sweden"}
	mapEuroGroups["F"] = []string{"Austria", "Hungary", "Iceland", "Portugal"}

	return mapEuroGroups
}

// MapOfTeamCodes is map containing the team codes.
// key: team name, value: code
// example: Brazil: BR
//
func (et EuroTournament2016) MapOfTeamCodes() map[string]string {

	var codes map[string]string
	codes = make(map[string]string)

	codes["Albania"] = "al"
	codes["France"] = "fr"
	codes["Romania"] = "ro"
	codes["Switzerland"] = "ch"
	codes["England"] = "gb"
	codes["Russia"] = "ru"
	codes["Slovakia"] = "sk"
	codes["Wales"] = "gb"
	codes["Germany"] = "de"
	codes["Northern Ireland"] = "gb"
	codes["Poland"] = "pl"
	codes["Ukraine"] = "ua"
	codes["Croatia"] = "hr"
	codes["Czech Republic"] = "cz"
	codes["Spain"] = "es"
	codes["Turkey"] = "tr"
	codes["Belgium"] = "be"
	codes["Italy"] = "it"
	codes["Republic of Ireland"] = "ie"
	codes["Sweden"] = "se"
	codes["Austria"] = "at"
	codes["Hungary"] = "hu"
	codes["Iceland"] = "is"
	codes["Portugal"] = "pt"

	return codes
}

// MapOfGroupMatches is a map containing the matches accessible by group.
//
// Example:
//
// 	Group A:[{"1", "Jun/12/2014", "Brazil", "Croatia", "Arena de São Paulo, São Paulo"}, ...]
//
func (et EuroTournament2016) MapOfGroupMatches() map[string][][]string {

	mapGroupMatches := make(map[string][][]string)

	const (
		cMatchID       = 0
		cMatchDate     = 1
		cMatchTeam1    = 2
		cMatchTeam2    = 3
		cMatchLocation = 4
	)

	mA1 := []string{"1", "Jun/10/2016", "France", "Romania", "Stade de France, Saint-Denis"}
	mA2 := []string{"2", "Jun/11/2016", "Albania", "Switzerland", "Stade Felix Bollaert, Lens"}
	mA3 := []string{"13", "Jun/15/2016", "Romania", "Switzerland", "Parc des Princes, Paris"}
	mA4 := []string{"14", "Jun/15/2016", "France", "Albania", "Stade Vélodrome, Marseille"}
	mA5 := []string{"25", "Jun/19/2016", "Switzerland", "France", "Stade Pierre Mauroy, Lille"}
	mA6 := []string{"26", "Jun/19/2016", "Romania", "Albania", "Parc OL, Lyon"}

	mB1 := []string{"3", "Jun/11/2016", "Wales", "Slovakia", "Matmut Atlantique, Bordeaux"}
	mB2 := []string{"4", "Jun/11/2016", "England", "Russia", "Stade Vélodrome, Marseille"}
	mB3 := []string{"15", "Jun/15/2016", "Russia", "Slovakia", "Stade Pierre Mauroy, Lille"}
	mB4 := []string{"16", "Jun/16/2016", "England", "Wales", "Stade Felix Bollaert, Lens"}
	mB5 := []string{"27", "Jun/20/2016", "Slovakia", "England", "Stade Geoffroy Guichard, Saint-Etienne"}
	mB6 := []string{"28", "Jun/20/2016", "Russia", "Wales", "Stadium de Toulouse, Toulouse"}

	mC1 := []string{"5", "Jun/12/2016", "Poland", "Northern Ireland", "Allianz Riviera, Nice"}
	mC2 := []string{"6", "Jun/12/2016", "Germany", "Ukraine", "Stade Pierre Mauroy, Lille"}
	mC3 := []string{"17", "Jun/16/2016", "Ukraine", "Northern Ireland", "Parc OL, Lyon"}
	mC4 := []string{"18", "Jun/16/2016", "Germany", "Poland", "Stade de France, Saint-Denis"}
	mC5 := []string{"29", "Jun/21/2016", "Northern Ireland", "Germany", "Parc des Princes, Paris"}
	mC6 := []string{"30", "Jun/21/2016", "Ukraine", "Poland", "Stade Vélodrome, Marseille"}

	mD1 := []string{"7", "Jun/12/2016", "Turkey", "Croatia", "Parc des Princes, Paris"}
	mD2 := []string{"8", "Jun/13/2016", "Spain", "Czech Republic", "Stadium de Toulouse, Toulouse"}
	mD3 := []string{"19", "Jun/17/2016", "Czech Republic", "Croatia", "Stade Geoffroy Guichard, Saint-Etienne"}
	mD4 := []string{"20", "Jun/17/2016", "Spain", "Turkey", "Allianz Riviera, Nice"}
	mD5 := []string{"31", "Jun/21/2016", "Croatia", "Spain", "Matmut Atlantique, Bordeaux"}
	mD6 := []string{"32", "Jun/21/2016", "Czech Republic", "Turkey", "Stade Felix Bollaert, Lens"}

	mE1 := []string{"9", "Jun/13/2016", "Republic of Ireland", "Sweden", "Stade de France, Saint-Denis"}
	mE2 := []string{"10", "Jun/13/2016", "Belgium", "Italy", "Parc OL, Lyon"}
	mE3 := []string{"21", "Jun/17/2016", "Italy", "Sweden", "Stadium de Toulouse, Toulouse"}
	mE4 := []string{"22", "Jun/18/2016", "Belgium", "Republic of Ireland", "Matmut Atlantique, Bordeaux"}
	mE5 := []string{"33", "Jun/22/2016", "Sweden", "Belgium", "Allianz Riviera, Nice"}
	mE6 := []string{"34", "Jun/22/2016", "Italy", "Republic of Ireland", "Stade Pierre Mauroy, Lille"}

	mF1 := []string{"11", "Jun/14/2016", "Austria", "Hungary", "Matmut Atlantique, Bordeaux"}
	mF2 := []string{"12", "Jun/14/2016", "Portugal", "Iceland", "Stade Geoffroy Guichard, Saint-Etienne"}
	mF3 := []string{"23", "Jun/18/2016", "Iceland", "Hungary", "Stade Vélodrome, Marseille"}
	mF4 := []string{"24", "Jun/18/2016", "Portugal", "Austria", "Parc des Princes, Paris"}
	mF5 := []string{"35", "Jun/22/2016", "Iceland", "Austria", "Stade de France, Saint-Denis"}
	mF6 := []string{"36", "Jun/22/2016", "Hungary", "Portugal", "Parc OL, Lyon"}

	var matchesA [][]string
	var matchesB [][]string
	var matchesC [][]string
	var matchesD [][]string
	var matchesE [][]string
	var matchesF [][]string

	matchesA = append(matchesA, mA1, mA2, mA3, mA4, mA5, mA6)
	matchesB = append(matchesB, mB1, mB2, mB3, mB4, mB5, mB6)
	matchesC = append(matchesC, mC1, mC2, mC3, mC4, mC5, mC6)
	matchesD = append(matchesD, mD1, mD2, mD3, mD4, mD5, mD6)
	matchesE = append(matchesE, mE1, mE2, mE3, mE4, mE5, mE6)
	matchesF = append(matchesF, mF1, mF2, mF3, mF4, mF5, mF6)

	mapGroupMatches["A"] = matchesA
	mapGroupMatches["B"] = matchesB
	mapGroupMatches["C"] = matchesC
	mapGroupMatches["D"] = matchesD
	mapGroupMatches["E"] = matchesE
	mapGroupMatches["F"] = matchesF

	return mapGroupMatches
}

// MapOf2ndRoundMatches returns the Map of 2nd round matches, of the euro tournament.
// key: round number, value: array of array of strings with match information ( MatchId, MatchDate, MatchTeam1, MatchTeam2, MatchLocation)
//
// Example:
//
//	round 16:[{"1", "Jun/12/2014", "Brazil", "Croatia", "Arena de São Paulo, São Paulo"}, ...]
//
func (et EuroTournament2016) MapOf2ndRoundMatches() map[string][][]string {

	// Round of 16
	m2nd1 := []string{"37", "Jun/25/2016", "2A", "2C", "Saint-Etienne"}
	m2nd2 := []string{"38", "Jun/25/2016", "1B", "3A/C/D", "Paris"}
	m2nd3 := []string{"39", "Jun/25/2016", "1D", "3B/E/F", "Lens"}
	m2nd4 := []string{"40", "Jun/26/2016", "1A", "3C/D/E", "Lyon"}
	m2nd5 := []string{"41", "Jun/26/2016", "1C", "3A/B/F", "Lille"}
	m2nd6 := []string{"42", "Jun/26/2016", "1F", "2E", "Toulouse"}
	m2nd7 := []string{"43", "Jun/27/2016", "1E", "2D", "Saint-Denis"}
	m2nd8 := []string{"44", "Jun/27/2016", "2B", "2F", "Nice"}
	// 17 Quarter-finals
	m2nd9 := []string{"45", "Jun/30/2016", "W37", "W39", "Marseille"}
	m2nd10 := []string{"46", "Jul/01/2016", "W38", "W42", "Lille"}
	m2nd11 := []string{"47", "Jul/02/2016", "W41", "W43", "Bordeaux"}
	m2nd12 := []string{"48", "Jul/03/2016", "W40", "W44", "Saint-Denis"}
	// 18 Semi-finals
	m2nd13 := []string{"49", "Jul/06/2016", "W45", "W46", "Lyon"}
	m2nd14 := []string{"50", "Jul/07/2016", "W47", "W48", "Marseille"}
	//19 Final
	m2nd15 := []string{"51", "Jul/10/2016", "W49", "50", "Saint-Denis"}

	var round16 [][]string
	var round17 [][]string
	var round18 [][]string
	var round19 [][]string

	round16 = append(round16, m2nd1, m2nd2, m2nd3, m2nd4, m2nd5, m2nd6, m2nd7, m2nd8)
	round17 = append(round17, m2nd9, m2nd10, m2nd11, m2nd12)
	round18 = append(round18, m2nd13, m2nd14)
	round19 = append(round19, m2nd15)

	mapMatches2ndRound := make(map[string][][]string)
	mapMatches2ndRound["16"] = round16
	mapMatches2ndRound["17"] = round17
	mapMatches2ndRound["18"] = round18
	mapMatches2ndRound["19"] = round19

	return mapMatches2ndRound
}

// ArrayOfPhases returns an array of the phases names of champions league tournament: (QuarterFinals, SemiFinals, Finals).
//
func (et EuroTournament2016) ArrayOfPhases() []string {
	return []string{cFirstStage, cRoundOf16, cQuarterFinals, cSemiFinals, cFinals}
}

// MapOfPhaseIntervals builds a map with key the corresponding phase in the world cup tournament
// at value a tuple that represent the match number interval in which the phase take place:
// first stage: matches 1 to 48
// Round of 16: matches 49 to 56
// Quarte-finals: matches 57 to 60
// Semi-finals: matches 61 to 62
// Third Place: match 63
// Finals: match 64
//
func (et EuroTournament2016) MapOfPhaseIntervals() map[string][]int64 {

	limits := make(map[string][]int64)
	limits[cFirstStage] = []int64{1, 36}
	limits[cRoundOf16] = []int64{37, 44}
	limits[cQuarterFinals] = []int64{45, 48}
	limits[cSemiFinals] = []int64{49, 50}
	limits[cFinals] = []int64{51, 51}
	return limits
}

// CreateEuro2016 creates euro 2016 tournament.
//
func CreateEuro2016(c appengine.Context, adminID int64) (*Tournament, error) {
	// create new tournament
	log.Infof(c, "Euro 2016: start")

	et := EuroTournament2016{}

	mapWCGroups := et.MapOfGroups()
	mapCountryCodes := et.MapOfTeamCodes()
	mapTeamID := make(map[string]int64)

	// mapGroupMatches is a map where the key is a string which represent the group
	// the key is a two dimensional string array. each element in the array represent a specific field in the match
	// map[string][][]string
	// example: {"1", "Jun/12/2014", "Brazil", "Croatia", "Arena de São Paulo, São Paulo"}
	mapGroupMatches := et.MapOfGroupMatches()

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
	matches1stStageIds := make([]int64, 6*len(matchesA))

	log.Infof(c, "Euro: maps ready")

	// build groups, and teams
	groups := make([]Tgroup, len(mapWCGroups))
	for groupName, teams := range mapWCGroups {
		log.Infof(c, "---------------------------------------")
		log.Infof(c, "Euro: working with group: %v", groupName)
		log.Infof(c, "Euro: teams: %v", teams)

		var group Tgroup
		group.Name = groupName
		groupIndex := int64(groupName[0]) - 65
		group.Teams = make([]Tteam, len(teams))
		group.Points = make([]int64, len(teams))
		group.GoalsF = make([]int64, len(teams))
		group.GoalsA = make([]int64, len(teams))
		for i, teamName := range teams {
			log.Infof(c, "Euro: team: %v", teamName)

			teamID, _, err1 := datastore.AllocateIDs(c, "Tteam", nil, 1)
			if err1 != nil {
				return nil, err1
			}
			log.Infof(c, "Euro: team: %v allocateIDs ok", teamName)

			teamkey := datastore.NewKey(c, "Tteam", "", teamID, nil)
			log.Infof(c, "Euro: team: %v NewKey ok", teamName)

			team := &Tteam{teamID, teamName, mapCountryCodes[teamName]}
			log.Infof(c, "Euro: team: %v instance of team ok", teamName)

			_, err := datastore.Put(c, teamkey, team)
			if err != nil {
				return nil, err
			}
			log.Infof(c, "Euro: team: %v put in datastore ok", teamName)
			group.Teams[i] = *team
			mapTeamID[teamName] = teamID
		}

		// build group matches:
		log.Infof(c, "Euro: building group matches")

		groupMatches := mapGroupMatches[groupName]
		group.Matches = make([]Tmatch, len(groupMatches))

		for matchIndex, matchData := range groupMatches {
			log.Infof(c, "Euro: match data: %v", matchData)

			matchID, _, err1 := datastore.AllocateIDs(c, "Tmatch", nil, 1)
			if err1 != nil {
				return nil, err1
			}

			log.Infof(c, "Euro: match: %v allocateIDs ok", matchID)

			matchkey := datastore.NewKey(c, "Tmatch", "", matchID, nil)
			log.Infof(c, "Euro: match: new key ok")

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
			log.Infof(c, "Euro: match: build match ok")

			_, err := datastore.Put(c, matchkey, match)
			if err != nil {
				return nil, err
			}
			log.Infof(c, "Euro: match: %v put in datastore ok", matchData)
			group.Matches[matchIndex] = *match

			// save in an array of int64 all the allocate IDs to store them in the tournament for easy retreival later on.
			matches1stStageIds[int64(matchInternalID)-1] = matchID
		}

		groupID, _, err1 := datastore.AllocateIDs(c, "Tgroup", nil, 1)
		if err1 != nil {
			return nil, err1
		}

		log.Infof(c, "Euro: Group: %v allocate ID ok", groupName)

		groupkey := datastore.NewKey(c, "Tgroup", "", groupID, nil)
		log.Infof(c, "Euro: Group: %v New Key ok", groupName)

		group.ID = groupID
		groups[groupIndex] = group
		_, err := datastore.Put(c, groupkey, &group)
		if err != nil {
			return nil, err
		}
		log.Infof(c, "Euro: Group: %v put in datastore ok", groupName)
	}

	// build array of group ids
	groupIds := make([]int64, 6)
	for i := range groupIds {
		groupIds[i] = groups[i].ID
	}

	log.Infof(c, "Euro: build of groups ids complete: %v", groupIds)

	// matches 2nd stage
	var matches2ndStageIds []int64
	var userIds []int64
	var teamIds []int64
	// mapMatches2ndRound  is a map where the key is a string which represent the rounds
	// the key is a two dimensional string array. each element in the array represent a specific field in the match
	// mapMatches2ndRound is a map[string][][]string
	// example: {"64", "Jul/13/2014", "W61", "W62", "Rio de Janeiro"}
	mapMatches2ndRound := et.MapOf2ndRoundMatches()

	// build matches 2nd phase
	for roundNumber, roundMatches := range mapMatches2ndRound {
		log.Infof(c, "Euro: building 2nd round matches: round number %v", roundNumber)
		for _, matchData := range roundMatches {
			log.Infof(c, "Euro: second phase match data: %v", matchData)

			matchID, _, err1 := datastore.AllocateIDs(c, "Tmatch", nil, 1)
			if err1 != nil {
				return nil, err1
			}

			log.Infof(c, "Euro: match: %v allocateIDs ok", matchID)

			matchkey := datastore.NewKey(c, "Tmatch", "", matchID, nil)
			log.Infof(c, "Euro: match: new key ok")

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
			log.Infof(c, "Euro: match 2nd round: build match ok")

			_, err := datastore.Put(c, matchkey, match)
			if err != nil {
				return nil, err
			}
			log.Infof(c, "Euro: 2nd round match: %v put in datastore ok", matchData)
			// save in an array of int64 all the allocate IDs to store them in the tournament for easy retreival later on.
			matches2ndStageIds = append(matches2ndStageIds, matchID)
		}
	}

	tstart, _ := time.Parse(shortForm, "Jun/10/2016")
	tend, _ := time.Parse(shortForm, "Jul/10/2016")
	adminIds := make([]int64, 1)
	adminIds[0] = adminID
	name := "2016 UEFA Euro"
	description := "France"
	var tournament *Tournament
	var err error
	if tournament, err = CreateTournament(c, name, description, tstart, tend, adminID); err != nil {
		log.Infof(c, "Euro: something went wrong when creating tournament.")
	} else {
		tournament.GroupIds = groupIds
		tournament.Matches1stStage = matches1stStageIds
		tournament.Matches2ndStage = matches2ndStageIds
		tournament.UserIds = userIds
		tournament.TeamIds = teamIds
		tournament.IsFirstStageComplete = false
		if err1 := tournament.Update(c); err1 != nil {
			log.Infof(c, "Euro: unable to udpate tournament.")
		}
	}

	log.Infof(c, "Euro: instance of tournament ready")

	return tournament, nil
}

// MapOfIDTeams builds a map of teams from tournament entity.
//
func (et EuroTournament2016) MapOfIDTeams(c appengine.Context, tournament *Tournament) map[int64]string {

	mapIDTeams := make(map[int64]string)

	groups := Groups(c, tournament.GroupIds)
	for _, g := range groups {
		for _, t := range g.Teams {
			mapIDTeams[t.ID] = t.Name
		}
	}
	return mapIDTeams
}
